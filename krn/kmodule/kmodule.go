// package kmodule -- модуль на основе ядра
package kmodule

import (
	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// kModule -- модуль на основе ядра
type kModule struct {
	kCtx IKernelCtx
	ctx  ILocalCtx
	name AModuleName
}

// NewKernelModule -- возвращает новый модуль на основе ядра
func NewKernelModule(name AModuleName) IKernelModule {
	Hassert(name != "", "NewKernelModule(): name is empty")
	kCtx := kctx.GetKernelCtx()
	sf := &kModule{
		kCtx: kCtx,
		ctx:  local_ctx.NewLocalCtx(kCtx.BaseCtx()),
		name: name,
	}
	return sf
}

// Log -- возвращает буферный лог
func (sf *kModule) Log() ILogBuf {
	return sf.ctx.Log()
}

// Ctx -- возвращает контекст модуля
func (sf *kModule) Ctx() ILocalCtx {
	return sf.kCtx
}

// Run -- запускает модуль в работу
func (sf *kModule) Run() {
	Hassert(false, "kModule.Run(): module='%v', parent not realised this method", sf.name)
}

// Name -- возвращает уникальное имя модуля
func (sf *kModule) Name() AModuleName {
	return sf.name
}

// IsWork -- возвращает признак состояния работы
func (sf *kModule) IsWork() bool {
	Hassert(false, "kModule.IsWork(): module='%v', parent not realised this method", sf.name)
	return false
}
