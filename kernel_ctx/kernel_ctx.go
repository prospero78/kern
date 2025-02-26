// package kernel_ctx -- контекст ядра
package kernel_ctx

import (
	"context"
	"log"
	"sync"

	"github.com/prospero78/kern/kernel_ctx/kernel_keeper"
	"github.com/prospero78/kern/kernel_ctx/kernel_wg"
	. "github.com/prospero78/kern/kernel_types"
)

// kernelCtx -- контекст ядра
type kernelCtx struct {
	ctxBg      context.Context        // Неотменяемый контекст ядра
	ctx        context.Context        // Отменяемый контекст ядра
	fnCancel   func()                 // Функция отмены контекста ядра
	dictVal    map[string]interface{} // Словарь различных значений
	kernKeeper IKernelKeeper          // Встроенный сторож отмены контекста системным сигналом
	kernWg     IKernelWg              // Встроенный ожидатель потока
	block      sync.RWMutex
}

var (
	kernCtx *kernelCtx // Глобальный объект контекста приложения
	block   sync.Mutex
)

// GetKernelCtx -- возвращает контекст ядра
func GetKernelCtx() IKernelCtx {
	block.Lock()
	defer block.Unlock()
	if kernCtx != nil {
		return kernCtx
	}
	ctxBg := context.Background()
	ctx, fnCancel := context.WithCancel(ctxBg)
	sf := &kernelCtx{
		ctxBg:    ctxBg,
		ctx:      ctx,
		fnCancel: fnCancel,
		dictVal:  map[string]interface{}{},
	}
	sf.kernWg = kernel_wg.GetKernelWg(sf.ctx)
	sf.kernKeeper = kernel_keeper.GetKernelKeeper(sf.ctx, sf.fnCancel, sf.kernWg)
	kernCtx = sf
	return kernCtx
}

// Wg -- возвращает ожидатель потоков
func (sf *kernelCtx) Wg() IKernelWg {
	return sf.kernWg
}

// Get -- возвращает хранимое значение
func (sf *kernelCtx) Get(key string) interface{} {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.dictVal[key]
}

// Del -- удаляет значение из контекста
func (sf *kernelCtx) Del(key string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	delete(sf.dictVal, key)
}

// Add -- добавляет значение в контекст
func (sf *kernelCtx) Add(key string, val interface{}) {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.dictVal[key] = val
}

// Done -- блокирующий вызов ожидания отмены контекста ядра
func (sf *kernelCtx) Done() {
	<-sf.ctx.Done()
	log.Println("kernelCtx.Done()")
}

// CtxBg -- возвращает неотменяемый контекст ядра (лучше не использовать)
func (sf *kernelCtx) CtxBg() context.Context {
	return sf.ctxBg
}

// Ctx -- возвращает контекст ядра
func (sf *kernelCtx) Ctx() context.Context {
	return sf.ctx
}

// Cancel -- отменяет контекст ядра
func (sf *kernelCtx) Cancel() {
	log.Println("kernelCtx.Cancel()")
	sf.fnCancel()
}
