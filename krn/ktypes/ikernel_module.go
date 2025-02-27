package ktypes

import (
	. "github.com/prospero78/kern/krn/kalias"
)

// IKernelModule -- интерфейс к модулю на основе ядра
type IKernelModule interface {
	// Run -- запускает модуль в работу
	Run()
	// IsWork -- возвращает состояние модуля
	IsWork() bool
	// Name -- возвращает уникальное имя модуля
	Name() AModuleName
	// Ctx -- возвращает контекст ядра
	Ctx() IKernelCtx
}
