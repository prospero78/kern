package kernel_types

import (
	. "github.com/prospero78/kern/kernel_alias"
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
