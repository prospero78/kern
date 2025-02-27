// package kmonolit -- модульный монолит на основе ядра
package kmonolit

import (
	"log"
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/safe_bool"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// kMonolit -- объект модульного монолита
type kMonolit struct {
	ctx     IKernelCtx
	name    string
	isLocal bool
	isWork  ISafeBool
	isEnd   bool
	dict    map[AModuleName]IKernelModule // Словарь модулей монолита
	block   sync.Mutex
}

// NewMonolit -- возвращает новый монолит
func NewMonolit(name string) IKernelMonolit {
	Hassert(name != "", "NewMonolit(): name is empty")
	ctx := kctx.GetKernelCtx()
	sf := &kMonolit{
		ctx:     ctx,
		name:    name,
		dict:    map[AModuleName]IKernelModule{},
		isWork:  safe_bool.NewSafeBool(),
		isLocal: ctx.Get("isLocal").Val().(bool),
	}
	ctx.Set("monolitName", name, "name of monolit")
	return sf
}

// Name -- возвращает имя монолита
func (sf *kMonolit) Name() string {
	return sf.name
}

// Add -- добавляет модуль в монолит
func (sf *kMonolit) Add(module IKernelModule) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(module != nil, "kMonolit.Add(): module==nil")
	_, isOk := sf.dict[module.Name()]
	Hassert(!isOk, "kMonolit.Add(): module(%v) already exists", module.Name())
	sf.dict[module.Name()] = module
	if sf.isWork.Get() {
		go module.Run()
	}
}

// Run -- запускает монолит в работу
func (sf *kMonolit) Run() {
	sf.block.Lock()
	defer sf.block.Unlock()
	if sf.isEnd {
		return
	}
	sf.isWork.Set()
	for _, module := range sf.dict {
		go module.Run()
	}
	go sf.close()
}

// IsLocal -- возвращает признак локальной шины
func (sf *kMonolit) IsLocal() bool {
	return sf.isLocal
}

// IsWork -- возвращает признак работы монолита
func (sf *kMonolit) IsWork() bool {
	return sf.isWork.Get()
}

// Ожидание завершения работы монолита
func (sf *kMonolit) close() {
	sf.ctx.Done()
	sf.ctx.Wg().Wait()
	sf.isWork.Reset()
	sf.isEnd = true
	log.Printf("kMonolit.close(): done")
}
