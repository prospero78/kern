// package kernel_types -- интерфейсы проекта
package kernel_types

import "context"

// IKernelCtx -- интерфейс к контексту ядра
type IKernelCtx interface {
	// CtxBg -- возвращает неотменяемый контекст ядра
	CtxBg() context.Context
	// Ctx -- возвращает отменяемы контекст ядра
	Ctx() context.Context
	// Cancel -- отменяет контекст ядра
	Cancel()
	// Done -- ожидает отмены контекста ядра
	Done()
	// Add -- добавляет значение в контекст
	Add(key string, val any)
	// Get -- извлекает значение из контекста
	Get(key string) any
	// Del -- удаляет значение из контекста
	Del(key string)
	// Wg -- возвращает ожидатель потоков
	Wg() IKernelWg
}
