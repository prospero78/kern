// package kern -- библиотека гибкого универсального облегчённого ядра для любого микросервиса
package kern

import (
	"github.com/prospero78/kern/kc/safe_bool"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kbus/kbus_http"
	"github.com/prospero78/kern/krn/kbus/kbus_local"
	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/krn/kernel_store"
	"github.com/prospero78/kern/krn/kmodule"
	"github.com/prospero78/kern/krn/kmonolit"
	"github.com/prospero78/kern/krn/kserv_http"
	. "github.com/prospero78/kern/krn/ktypes"
)

// NewKernelCtx -- возвращает контекст ядра
func NewKernelCtx() IKernelCtx {
	ctx := kctx.GetKernelCtx()
	return ctx
}

// NewKernelStore -- возвращает хранилище ядра
func NewKernelStore() IKernelStore {
	store := kernel_store.GetKernelStore()
	return store
}

// NewKernelServerHttp -- возвращает веб-сервер ядра
func NewKernelServerHttp() IKernelServerHttp {
	kernServHttp := kserv_http.GetKernelServHttp()
	return kernServHttp
}

// NewSafeBool -- возвращает новый потокобезопасный булевый признак
func NewSafeBool() ISafeBool {
	sb := safe_bool.NewSafeBool()
	return sb
}

// NewKernelBusLocal -- возвращает локальную шину данных
func NewKernelBusLocal() IKernelBus {
	bus := kbus_local.GetKernelBusLocal()
	return bus
}

// NewKernelBusHttp -- возвращает HTTP шину данных
func NewKernelBusHttp() IKernelBus {
	bus := kbus_http.GetKernelBusHttp()
	return bus
}

// NewMonolitLocal -- возвращает монолит с локальной шиной
func NewMonolitLocal() IKernelMonolit {
	ctx := kctx.GetKernelCtx()
	_ = kbus_local.GetKernelBusLocal()
	ctx.Set("isLocal", true)
	monolit := kmonolit.NewMonolit()
	return monolit
}

// NewMonolitHttp -- возвращает монолит с локальной шиной поверх HTTP
func NewMonolitHttp() IKernelMonolit {
	ctx := kctx.GetKernelCtx()
	_ = kbus_http.GetKernelBusHttp()
	ctx.Set("isLocal", false)
	monolit := kmonolit.NewMonolit()
	return monolit
}

// NewKernelModule -- возвращает новый модуль на ядре
func NewKernelModule(name AModuleName) IKernelModule {
	mod := kmodule.NewKernelModule(name)
	return mod
}
