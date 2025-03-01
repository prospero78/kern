// package mod_kctx -- модуль контекста ядра
package mod_kctx

import (
	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/krn/kmodule"
	"github.com/prospero78/kern/krn/kserv_http"
	. "github.com/prospero78/kern/krn/ktypes"
	"github.com/prospero78/kern/mds/mod_serv_http/http_api"
	"github.com/prospero78/kern/mds/mod_serv_http/page_module"
	"github.com/prospero78/kern/mds/mod_serv_http/page_modules"
	"github.com/prospero78/kern/mds/mod_serv_http/page_monolit"
)

// ModuleKernelCtx -- модуль контекста ядра
type ModuleKernelCtx struct {
	IKernelModule
	ctx       IKernelCtx
	kServHttp IKernelServerHttp
	log       ILogBuf
}

// NewModuleKernelCtx -- возвращает новый модуль контекста ядра
func NewModuleKernelCtx() *ModuleKernelCtx {
	sf := &ModuleKernelCtx{
		ctx:           kctx.GetKernelCtx(),
		IKernelModule: kmodule.NewKernelModule("kCtx"),
		kServHttp:     kserv_http.GetKernelServHttp(),
	}
	sf.log = sf.ctx.Log()
	_ = page_monolit.GetPageMonolit()
	_ = page_modules.GetPageModules()
	_ = page_module.GetPageModule()

	_ = http_api.NewHttpApi()
	return sf
}

// Run -- запускает модуль в работу
func (sf *ModuleKernelCtx) Run() {
	sf.log.Info("ModuleKernelCtx.Run(): module=%v, is run", sf.Name())
	go sf.kServHttp.Run()
}

// IsWork -- признак работы модуля
func (sf *ModuleKernelCtx) IsWork() bool {
	return sf.ctx.Wg().IsWork()
}
