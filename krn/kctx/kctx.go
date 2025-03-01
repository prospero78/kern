// package kctx -- контекст ядра
package kctx

import (
	"context"
	"sync"

	"github.com/prospero78/kern/kc/local_ctx"
	"github.com/prospero78/kern/krn/kctx/kernel_keeper"
	"github.com/prospero78/kern/krn/kctx/kwg"
	. "github.com/prospero78/kern/krn/ktypes"
)

// kCtx -- контекст ядра
type kCtx struct {
	ILocalCtx
	log        ILogBuf
	ctxBg      context.Context // Неотменяемый контекст ядра
	ctx        context.Context // Отменяемый контекст ядра
	fnCancel   func()          // Функция отмены контекста ядра
	kernKeeper IKernelKeeper   // Встроенный сторож отмены контекста системным сигналом
	kernWg     IKernelWg       // Встроенный ожидатель потока
}

var (
	kernCtx *kCtx // Глобальный объект контекста приложения
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
	sf := &kCtx{
		ctxBg:    ctxBg,
		ctx:      ctx,
		fnCancel: fnCancel,
	}
	sf.ILocalCtx = local_ctx.NewLocalCtx(sf.ctx)
	sf.log = sf.Log()
	sf.kernWg = kwg.GetKernelWg(sf.ctx)
	sf.kernKeeper = kernel_keeper.GetKernelKeeper(sf.ctx, sf.fnCancel, sf.kernWg)
	kernCtx = sf
	return kernCtx
}

// Keeper -- возвращает сторож системных сигналов
func (sf *kCtx) Keeper() IKernelKeeper {
	return sf.kernKeeper
}

// Wg -- возвращает ожидатель потоков
func (sf *kCtx) Wg() IKernelWg {
	return sf.kernWg
}

// Done -- блокирующий вызов ожидания отмены контекста ядра
func (sf *kCtx) Done() {
	<-sf.ctx.Done()
	sf.log.Debug("kCtx.Done()")
}

// CtxBg -- возвращает неотменяемый контекст ядра (лучше не использовать)
func (sf *kCtx) CtxBg() context.Context {
	return sf.ctxBg
}

// BaseCtx -- возвращает контекст ядра
func (sf *kCtx) BaseCtx() context.Context {
	return sf.ctx
}

// Cancel -- отменяет контекст ядра
func (sf *kCtx) Cancel() {
	sf.fnCancel()
	sf.log.Debug("kCtx.Cancel()")
}
