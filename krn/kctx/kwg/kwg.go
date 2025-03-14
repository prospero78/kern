// package kwg -- именованный ожидатель потоков ядра
//
// Не позволяет завершиться ядру, если есть хоть один работающий поток
package kwg

import (
	"context"
	"fmt"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	"gitp78su.ipnodns.ru/svi/kern/kc/log_buf"
	"gitp78su.ipnodns.ru/svi/kern/kc/safe_bool"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// kernelWg -- именованный ожидатель потоков ядра
type kernelWg struct {
	sync.RWMutex
	ctx        context.Context
	dictStream map[AStreamName]bool // Словарь имён потоков с признаком работы
	isWork     ISafeBool
	log        ILogBuf
}

var (
	kernWg *kernelWg // Глобальный объект
	block  sync.Mutex
)

// GetKernelWg -- возвращает новый именованный ожидатель потоков ядра
func GetKernelWg(ctx context.Context) IKernelWg {
	block.Lock()
	defer block.Unlock()
	if kernWg != nil {
		kernWg.log.Debug("GetKernelWg()")
		return kernWg
	}
	Hassert(ctx != nil, "GetKernelWg(): ctx==nil")
	sf := &kernelWg{
		ctx:        ctx,
		dictStream: map[AStreamName]bool{},
		isWork:     safe_bool.NewSafeBool(),
		log:        log_buf.NewLogBuf(),
	}
	sf.log.Debug("GetKernelWg(): run")
	go sf.close()
	sf.isWork.Set()
	kernWg = sf
	return kernWg
}

// Log -- возвращает лог ожидателя потоков
func (sf *kernelWg) Log() ILogBuf {
	return sf.log
}

// Len -- возвращает размер списка ожидания потоков
func (sf *kernelWg) Len() int {
	sf.RLock()
	defer sf.RUnlock()
	return len(sf.dictStream)
}

// IsWork -- возвращает признак работы ядра
func (sf *kernelWg) IsWork() bool {
	return sf.isWork.Get()
}

// List -- возвращает список имён потоков на ожидании
func (sf *kernelWg) List() []AStreamName {
	sf.RLock()
	defer sf.RUnlock()
	lst := []AStreamName{}
	for name := range sf.dictStream {
		lst = append(lst, name)
	}
	return lst
}

// Done -- удаляет поток из ожидания
func (sf *kernelWg) Done(name AStreamName) {
	sf.Lock()
	defer sf.Unlock()
	delete(sf.dictStream, name)
	sf.log.Debug("kernelWg.Done(): stream(%v) done", name)
}

// Wait -- блокирующий вызов; возвращает управление, только когда все потоки завершили работу
func (sf *kernelWg) Wait() {
	for {
		SleepMs()
		if !sf.isWork.Get() {
			break
		}
	}
	sf.log.Debug("kernelWg.Wait(): done")
}

// Add -- добавляет поток в ожидание
func (sf *kernelWg) Add(name AStreamName) error {
	sf.Lock()
	defer sf.Unlock()
	sf.log.Debug("kernelWg.Add(): stream='%v'", name)
	if !sf.isWork.Get() {
		return fmt.Errorf("kernelWg.Add(): stream=%v, work end", name)
	}
	Hassert(name != "", "kernelWg.Add(): name stream is empty")
	_, isOk := sf.dictStream[name]
	Hassert(!isOk, "kernelWg.Add(): stream '%v' already exists", name)
	sf.dictStream[name] = true
	return nil
}

// Ожидает окончания работы ожидателя групп
func (sf *kernelWg) close() {
	<-sf.ctx.Done()
	fnDone := func() bool {
		sf.Lock()
		defer sf.Unlock()
		return len(sf.dictStream) == 0
	}
	for {
		SleepMs()
		if fnDone() {
			break
		}
	}
	sf.Lock()
	defer sf.Unlock()
	if !sf.isWork.Get() {
		return
	}
	sf.isWork.Reset()
	sf.log.Debug("kernelWg.close(): end")
}
