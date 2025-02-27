// package kmodule -- модуль на основе ядра
package kmodule

import (
	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// kernelModule -- модуль на основе ядра
type kernelModule struct {
	kCtx IKernelCtx
	ctx  ILocalCtx
	name AModuleName
}

// NewKernelModule -- возвращает новый модуль на основе ядра
func NewKernelModule(name AModuleName) IKernelModule {
	Hassert(name != "", "NewKernelModule(): name is empty")
	kCtx := kctx.GetKernelCtx()
	sf := &kernelModule{
		kCtx: kCtx,
		ctx:  local_ctx.NewLocalCtx(kCtx.BaseCtx()),
		name: name,
	}
	return sf
}

// Ctx -- возвращает контекст модуля
func (sf *kernelModule) Ctx() ILocalCtx {
	return sf.kCtx
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
