// package kctx -- контекст ядра
package kctx

import (
	"context"
	"log"
	"sync"

	"github.com/prospero78/kern/kc/local_ctx"
	"github.com/prospero78/kern/krn/kctx/kernel_keeper"
	"github.com/prospero78/kern/krn/kctx/kwg"
	. "github.com/prospero78/kern/krn/ktypes"
)

// kernelCtx -- контекст ядра
type kernelCtx struct {
	ILocalCtx
	ctxBg      context.Context // Неотменяемый контекст ядра
	ctx        context.Context // Отменяемый контекст ядра
	fnCancel   func()          // Функция отмены контекста ядра
	kernKeeper IKernelKeeper   // Встроенный сторож отмены контекста системным сигналом
	kernWg     IKernelWg       // Встроенный ожидатель потока
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
	}
	sf.ILocalCtx = local_ctx.NewLocalCtx(sf.ctx)
	sf.kernWg = kwg.GetKernelWg(sf.ctx)
	sf.kernKeeper = kernel_keeper.GetKernelKeeper(sf.ctx, sf.fnCancel, sf.kernWg)
	kernCtx = sf
	return kernCtx
}

// Wg -- возвращает ожидатель потоков
func (sf *kernelCtx) Wg() IKernelWg {
	return sf.kernWg
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

// BaseCtx -- возвращает контекст ядра
func (sf *kernelCtx) BaseCtx() context.Context {
	return sf.ctx
}

// Cancel -- отменяет контекст ядра
func (sf *kernelCtx) Cancel() {
	log.Println("kernelCtx.Cancel()")
	sf.fnCancel()
}
