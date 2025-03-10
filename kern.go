// package kern -- библиотека гибкого универсального облегчённого ядра для любого микросервиса
package kern

import (
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	"gitp78su.ipnodns.ru/svi/kern/kc/safe_bool"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_http"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_http/client_bus_http"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_local"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_local/client_bus_local"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	"gitp78su.ipnodns.ru/svi/kern/krn/kmodule"
	"gitp78su.ipnodns.ru/svi/kern/krn/kmonolit"
	"gitp78su.ipnodns.ru/svi/kern/krn/kserv_http"
	"gitp78su.ipnodns.ru/svi/kern/krn/kstore_kv"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_kctx"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_keeper"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http"
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
	for {
		SleepMs()
		if ctx.Get("isLocal") != nil {
			break
		}
	}
	monolit := kmonolit.GetMonolit(name)
	_ = kbus_local.GetKernelBusLocal()
	return monolit
}

// NewMonolitHttp -- возвращает монолит с локальной шиной поверх HTTP
func NewMonolitHttp(name string) IKernelMonolit {
	ctx := kctx.GetKernelCtx()
	_ = kbus_http.GetKernelBusHttp()
	ctx.Set("isLocal", false, "bus type")
	for {
		SleepMs()
		if ctx.Get("isLocal") != nil {
			break
		}
	}
	monolit := kmonolit.GetMonolit(name)
	return monolit
}

// NewKernelModule -- возвращает новый модуль на ядре
func NewKernelModule(name AModuleName) IKernelModule {
	mod := kmodule.NewKernelModule(name)
	return mod
}

// NewClientBusLocal -- возвращает клиент для работы с локальной шиной
func NewClientBusLocal() IBusClient {
	client := client_bus_local.NewClientBusLocal()
	return client
}

// NewClientBusHttp -- возвращает клиент для работы с шиной поверх HTTP
func NewClientBusHttp(url string) IBusClient {
	client := client_bus_http.NewClientBusHttp(url)
	return client
}

// NewModuleServHttp -- возвращает новый модуль для IKernelServHttp
func NewModuleServHttp() IKernelModule {
	modServHttp := mod_serv_http.NewModuleServHttp()
	return modServHttp
}

// NewModuleKernelCtx -- возвращает новый модуль для IKernelCtx
func NewModuleKernelCtx() IKernelModule {
	modKernelCtx := mod_kctx.NewModuleKernelCtx()
	return modKernelCtx
}

// NewModuleKernelKeeper -- возвращает новый модуль для IKernelKeeper
func NewModuleKernelKeeper() IKernelModule {
	modKernelKeeper := mod_keeper.NewModuleKeeper()
	return modKernelKeeper
}
