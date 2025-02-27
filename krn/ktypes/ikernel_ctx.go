// package kernel_types -- интерфейсы проекта
package ktypes

import "context"

// IKernelCtx -- интерфейс к контексту ядра
type IKernelCtx interface {
	ILocalCtx
	// CtxBg -- возвращает неотменяемый контекст ядра
	CtxBg() context.Context
	// BaseCtx -- возвращает отменяемы контекст ядра
	BaseCtx() context.Context
	// Wg -- возвращает ожидатель потоков
	Wg() IKernelWg
}
