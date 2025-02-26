// package kern -- библиотека гибкого универсального облегчённого ядра для любого микросервиса
package kern

import (
	"github.com/prospero78/kern/kc/safe_bool"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kbus/kbus_http"
	"github.com/prospero78/kern/krn/kbus/kbus_local"
	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/krn/kmodule"
	"github.com/prospero78/kern/krn/kmonolit"
	"github.com/prospero78/kern/krn/kserv_http"
	"github.com/prospero78/kern/krn/kstore_kv"
	. "github.com/prospero78/kern/krn/ktypes"
)

// NewKernelCtx -- возвращает контекст ядра
func NewKernelCtx() IKernelCtx {
	ctx := kctx.GetKernelCtx()
	return ctx
}

// NewKernelStoreKv -- возвращает быстрое key-value хранилище ядра
func NewKernelStoreKv() IKernelStoreKv {
	store := kstore_kv.GetKernelStore()
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
	ctx := kctx.GetKernelCtx()
	ctx.Set("monolitName", "unknown monolit", "GetKernelBusLocal()")
	bus := kbus_local.GetKernelBusLocal()
	return bus
}

// NewKernelBusHttp -- возвращает HTTP шину данных
func NewKernelBusHttp() IKernelBus {
	bus := kbus_http.GetKernelBusHttp()
	return bus
}

// NewMonolitLocal -- возвращает монолит с локальной шиной
func NewMonolitLocal(name string) IKernelMonolit {
	ctx := kctx.GetKernelCtx()
	ctx.Set("isLocal", true, "bus type")
	monolit := kmonolit.NewMonolit(name)
	_ = kbus_local.GetKernelBusLocal()
	return monolit
}

// NewMonolitHttp -- возвращает монолит с локальной шиной поверх HTTP
func NewMonolitHttp(name string) IKernelMonolit {
	ctx := kctx.GetKernelCtx()
	_ = kbus_http.GetKernelBusHttp()
	ctx.Set("isLocal", false, "bus type")
	monolit := kmonolit.NewMonolit(name)
	return monolit
}

// NewKernelModule -- возвращает новый модуль на ядре
func NewKernelModule(name AModuleName) IKernelModule {
	mod := kmodule.NewKernelModule(name)
	return mod
}
