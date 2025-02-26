// package kernel_module -- модуль на основе ядра
package kernel_module

import (
	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	"github.com/prospero78/kern/kernel_ctx"
	. "github.com/prospero78/kern/kernel_types"
)

// kernelModule -- модуль на основе ядра
type kernelModule struct {
	ctx  IKernelCtx
	name AModuleName
}

// NewKernelModule -- возвращает новый модуль на основе ядра
func NewKernelModule(name AModuleName) IKernelModule {
	Hassert(name != "", "NewKernelModule(): name is empty")
	sf := &kernelModule{
		ctx:  kernel_ctx.GetKernelCtx(),
		name: name,
	}
	return sf
}

// Ctx -- возвращает контекст ядра
func (sf *kernelModule) Ctx() IKernelCtx {
	return sf.ctx
}

// Run -- запускает модуль в работу
func (sf *kernelModule) Run() {
	Hassert(false, "kernelModule.Run(): module='%v', parent not realised this method", sf.name)
}

// Name -- возвращает уникальное имя модуля
func (sf *kernelModule) Name() AModuleName {
	return sf.name
}

// IsWork -- возвращает признак состояния работы
func (sf *kernelModule) IsWork() bool {
	Hassert(false, "kernelModule.IsWork(): module='%v', parent not realised this method", sf.name)
	return false
}
