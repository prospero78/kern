// package kern -- библиотека гибкого универсального облегчённого ядра для любого микросервиса
package kern

import (
	"github.com/prospero78/kern/kernel_bus_local"
	"github.com/prospero78/kern/kernel_ctx"
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
func NewKernelStore(ctx IKernelCtx) IKernelStore {
	store := kernel_store.GetKernelStore(ctx)
	return store
}

// NewKernelServerHttp -- возвращает веб-сервер ядра
func NewKernelServerHttp(ctx IKernelCtx) IKernelServerHttp {
	kernServHttp := kernel_serv_http.GetKernelServHttp(ctx)
	return kernServHttp
}

// NewSafeBool -- возвращает новый потокобезопасный булевый признак
func NewSafeBool() ISafeBool {
	sb := safe_bool.NewSafeBool()
	return sb
}

// NewKernelBusLocal -- возвращает локальную шину данных
func NewKernelBusLocal(ctx IKernelCtx) IKernelBus {
	bus := kernel_bus_local.GetKernelBusLocal(ctx)
	return bus
}
