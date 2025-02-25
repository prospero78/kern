// package kernel -- библиотека гибкого универсального облегчённого ядра для любого микросервиса
package kernel

import (
	"github.com/prospero78/kern/kernel_ctx"
	"github.com/prospero78/kern/kernel_serv_http"
	"github.com/prospero78/kern/kernel_store"
	. "github.com/prospero78/kern/kernel_types"
	"github.com/prospero78/kern/safe_bool"
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
