// package kernel -- библиотека гибкого универсального облегчённого ядра для любого микросервиса
package kernel

import (
	"github.com/svi/kern/kernel_ctx"
	"github.com/svi/kern/kernel_serv_http"
	"github.com/svi/kern/kernel_store"
	. "github.com/svi/kern/kernel_types"
	"github.com/svi/kern/safe_bool"
)

// NewKernelCtx -- возвращает новый контекст ядра
func NewKernelCtx() IKernelCtx {
	ctx := kernel_ctx.GetKernelCtx()
	return ctx
}

// NewKernelStore -- возвращает новое хранилище ядра
func NewKernelStore(ctx IKernelCtx) IKernelStore {
	store := kernel_store.GetKernelStore(ctx)
	return store
}

// NewKernelServerHttp -- возвращает новый веб-сервер ядра
func NewKernelServerHttp(ctx IKernelCtx) IKernelServerHttp {
	kernServHttp := kernel_serv_http.GetKernelServHttp(ctx)
	return kernServHttp
}

// NewSafeBool -- возвращает новый потокобезопасный булевый признак
func NewSafeBool() ISafeBool {
	sb := safe_bool.NewSafeBool()
	return sb
}
