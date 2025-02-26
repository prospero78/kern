// package kern -- библиотека гибкого универсального облегчённого ядра для любого микросервиса
package kern

import (
	. "github.com/prospero78/kern/kernel_alias"
	"github.com/prospero78/kern/kernel_bus/kernel_bus_http"
	"github.com/prospero78/kern/kernel_bus/kernel_bus_local"
	"github.com/prospero78/kern/kernel_ctx"
	"github.com/prospero78/kern/kernel_module"
	"github.com/prospero78/kern/kernel_monolit"
	"github.com/prospero78/kern/kernel_serv_http"
	"github.com/prospero78/kern/kernel_store"
	. "github.com/prospero78/kern/kernel_types"
	"github.com/prospero78/kern/safe_bool"
)

// NewKernelCtx -- возвращает контекст ядра
func NewKernelCtx() IKernelCtx {
	ctx := kernel_ctx.GetKernelCtx()
	return ctx
}

// NewKernelStore -- возвращает хранилище ядра
func NewKernelStore() IKernelStore {
	store := kernel_store.GetKernelStore()
	return store
}

// NewKernelServerHttp -- возвращает веб-сервер ядра
func NewKernelServerHttp() IKernelServerHttp {
	kernServHttp := kernel_serv_http.GetKernelServHttp()
	return kernServHttp
}

// NewSafeBool -- возвращает новый потокобезопасный булевый признак
func NewSafeBool() ISafeBool {
	sb := safe_bool.NewSafeBool()
	return sb
}

// NewKernelBusLocal -- возвращает локальную шину данных
func NewKernelBusLocal() IKernelBus {
	bus := kernel_bus_local.GetKernelBusLocal()
	return bus
}

// NewKernelBusHttp -- возвращает HTTP шину данных
func NewKernelBusHttp() IKernelBus {
	bus := kernel_bus_http.GetKernelBusHttp()
	return bus
}

// NewMonolitLocal -- возвращает монолит с локальной шиной
func NewMonolitLocal() IKernelMonolit {
	ctx := kernel_ctx.GetKernelCtx()
	_ = kernel_bus_local.GetKernelBusLocal()
	ctx.Set("isLocal", true)
	monolit := kernel_monolit.NewMonolit()
	return monolit
}

// NewMonolitHttp -- возвращает монолит с локальной шиной поверх HTTP
func NewMonolitHttp() IKernelMonolit {
	ctx := kernel_ctx.GetKernelCtx()
	_ = kernel_bus_http.GetKernelBusHttp()
	ctx.Set("isLocal", false)
	monolit := kernel_monolit.NewMonolit()
	return monolit
}

// NewKernelModule -- возвращает новый модуль на ядре
func NewKernelModule(name AModuleName) IKernelModule {
	mod := kernel_module.NewKernelModule(name)
	return mod
}
