// package kernel_store -- локальное хранилище ядра
package kernel_store

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"

	. "github.com/svi/kern/helpers"
	. "github.com/svi/kern/kernel_types"
	"github.com/svi/kern/safe_bool"
)

const (
	storeStreamName = "kernel_store" // Имя потока для ожидателя потоков
)

// kernelStore -- локальное хранилище ядра
type kernelStore struct {
	ctx        IKernelCtx
	wg         IKernelWg
	storePath  string
	db         *badger.DB
	isWork     ISafeBool
	block      sync.RWMutex
	blockClose sync.Mutex
}

var (
	kernStore *kernelStore // Глобальный объект
	block     sync.Mutex
)

// GetKernelStore -- возвращает новое локальное хранилище ядра
func GetKernelStore(ctx IKernelCtx) IKernelStore {
	log.Println("GetKernelStore()")
	block.Lock()
	defer block.Unlock()
	if kernStore != nil {
		return kernStore
	}
	Hassert(ctx != nil, "GetKernelStore(): IKernelCtx==nil")
	sf := &kernelStore{
		ctx:    ctx,
		wg:     ctx.Wg(),
		isWork: safe_bool.NewSafeBool(),
	}
	sf.open()
	kernStore = sf
	ctx.Add("kernStore", kernStore)
	return kernStore
}

// Set -- устанавливает значение по ключу
func (sf *kernelStore) Set(key string, val []byte) error {
	sf.block.Lock()
	defer sf.block.Unlock()
	// if !sf.isWork.Get() {
	// 	return fmt.Errorf("kernelStore.Set(): DB already close")
	// }
	fnSet := func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), val)
		return err
	}
	err := sf.db.Update(fnSet)
	if err != nil {
		return fmt.Errorf("kernelStore.Set(): key=%v, err=\n\t%w", key, err)
	}
	return nil
}

// Get -- возвращает значение по ключу
func (sf *kernelStore) Get(key string) ([]byte, error) {
	sf.block.RLock()
	defer sf.block.RUnlock()
	var binVal []byte
	fnGet := func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		binVal, err = item.ValueCopy(binVal)
		return err
	}
	err := sf.db.View(fnGet)
	if err != nil {
		return nil, fmt.Errorf("kernelStore.Delete(): key=%v, err=\n\t%w", key, err)
	}
	return binVal, nil
}

// Delete -- удалить ключ из хранилища
func (sf *kernelStore) Delete(key string) error {
	sf.block.Lock()
	defer sf.block.Unlock()

	fnDelete := func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	}
	err := sf.db.Update(fnDelete)
	if err != nil {
		return fmt.Errorf("kernelStore.Delete(): key=%v, err=\n\t%w", key, err)
	}
	return nil
}

// Открывает базу при создании
func (sf *kernelStore) open() {
	sf.block.Lock()
	defer sf.block.Unlock()
	strPath := os.Getenv("LOCAL_STORE_PATH")
	Hassert(strPath != "", "kernelStore.open(): env LOCAL_STORE_PATH not set")
	pwd, err := os.Getwd()
	Hassert(err == nil, "kernelStore.open(): in get PWD, err=\n\t%v", err)
	sf.storePath = pwd + strPath + "/db_local"
	err = os.MkdirAll(sf.storePath, 0750)
	Hassert(err == nil, "kernelStore.open(): in make dir %v, err=\n\t%v", sf.storePath, err)
	sf.db, err = badger.Open(badger.DefaultOptions(sf.storePath))
	Hassert(err == nil, "kernelStore.open(): in open DB %v, err=\n\t%v", sf.storePath, err)
	err = sf.wg.Add(storeStreamName)
	Hassert(err == nil, "kernelStore.open(): in add name stream to IKernelWg, err=\n\t%v", err)
	sf.isWork.Set()
	go sf.close()
	go sf.clean()
}

// Выполняет периодическую сборку мусора в файле
func (sf *kernelStore) clean() {
	chRun := make(chan int, 2)
	defer close(chRun)
	fnClean := func() {
		sf.block.Lock()
		defer sf.block.Unlock()
		_ = sf.db.RunValueLogGC(0.7)
	}
	chRun <- 1
	for {
		select {
		case <-sf.ctx.Ctx().Done(): // надо прекратить работу
			return
		case <-chRun: // Пора поработать
			fnClean()
		}
		time.Sleep(time.Second * 1)
	}
}

// Ожидает последнего потока под отдельной блокировкой
func (sf *kernelStore) wait(chWait chan int) {
	for {
		time.Sleep(time.Millisecond * 5)
		if sf.wg.Len() <= 1 {
			break
		}
	}
	close(chWait)
}

// Ожидает закрытия контекста ядра, закрывает хранилище
func (sf *kernelStore) close() {
	<-sf.ctx.Ctx().Done()
	sf.blockClose.Lock()
	defer sf.blockClose.Unlock()
	if !sf.isWork.Get() {
		return
	}
	chWait := make(chan int, 2)
	go sf.wait(chWait)
	<-chWait
	sf.isWork.Reset()
	err := sf.db.Close()
	Assert(err == nil, "kernelStore.close(): in close DB, err=\n\t%v", err)
	sf.wg.Done(storeStreamName)
	log.Println("kernelStore.close(): done")
}
