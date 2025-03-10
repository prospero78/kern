// package mod_kctx -- модуль контекста ядра
package mod_kctx

import (
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	"gitp78su.ipnodns.ru/svi/kern/krn/kmodule"
	"gitp78su.ipnodns.ru/svi/kern/krn/kserv_http"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http/http_api"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http/page_module"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http/page_modules"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http/page_monolit"
)

// ModuleKernelCtx -- модуль контекста ядра
type ModuleKernelCtx struct {
	IKernelModule
	kCtx      IKernelCtx
	kServHttp IKernelServerHttp
	log       ILogBuf
}

// NewModuleKernelCtx -- возвращает новый модуль контекста ядра
func NewModuleKernelCtx() *ModuleKernelCtx {
	sf := &ModuleKernelCtx{
		kCtx:          kctx.GetKernelCtx(),
		IKernelModule: kmodule.NewKernelModule("kCtx"),
		kServHttp:     kserv_http.GetKernelServHttp(),
	}
	sf.log = sf.Ctx().Log()
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
	return sf.kCtx.Wg().IsWork()
}
