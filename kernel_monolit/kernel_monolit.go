// package kernel_monolit -- модульный монолит на основе ядра
package kernel_monolit

import (
	"log"
	"sync"

	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	"github.com/prospero78/kern/kernel_ctx"
	. "github.com/prospero78/kern/kernel_types"
	"github.com/prospero78/kern/safe_bool"
)

// kernMonolit -- объект модульного монолита
type kernMonolit struct {
	ctx     IKernelCtx
	isLocal bool
	isWork  ISafeBool
	dict    map[AModuleName]IKernelModule // Словарь модулей монолита
	block   sync.Mutex
}

// NewMonolit -- возвращает новый монолит
func NewMonolit() IKernelMonolit {
	ctx := kernel_ctx.GetKernelCtx()
	sf := &kernMonolit{
		ctx:     ctx,
		dict:    map[AModuleName]IKernelModule{},
		isWork:  safe_bool.NewSafeBool(),
		isLocal: ctx.Get("isLocal").(bool),
	}
	return sf
}

// Add -- добавляет модуль в монолит
func (sf *kernMonolit) Add(module IKernelModule) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(module != nil, "kernMonolit.Add(): module==nil")
	_, isOk := sf.dict[module.Name()]
	Hassert(!isOk, "kernMonolit.Add(): module(%v) already exists", module.Name())
	sf.dict[module.Name()] = module
	if sf.isWork.Get() {
		go module.Run()
	}
}

// Run -- запускает монолит в работу
func (sf *kernMonolit) Run() {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.isWork.Set()
	for _, module := range sf.dict {
		go module.Run()
	}
	go sf.close()
}

// IsLocal -- возвращает признак локальной шины
func (sf *kernMonolit) IsLocal() bool {
	return sf.isLocal
}

// IsWork -- возвращает признак работы монолита
func (sf *kernMonolit) IsWork() bool {
	return sf.isWork.Get()
}

// Ожидание завершения работы монолита
func (sf *kernMonolit) close() {
	sf.ctx.Wg().Wait()
	sf.isWork.Reset()
	log.Printf("kernMonolit.close(): done")
}
