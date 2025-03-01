// package mod_keeper -- модуль сторожа ядра
package mod_keeper

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

// ModuleKeeper -- модуль сторожа
type ModuleKeeper struct {
	IKernelModule
	kCtx      IKernelCtx
	kServHttp IKernelServerHttp
	log       ILogBuf
}

// NewModuleKeeper -- возвращает новый модуль сторожа ядра
func NewModuleKeeper() *ModuleKeeper {
	sf := &ModuleKeeper{
		kCtx:          kctx.GetKernelCtx(),
		IKernelModule: kmodule.NewKernelModule("kKeeper"),
		kServHttp:     kserv_http.GetKernelServHttp(),
	}
	sf.log = sf.kCtx.Keeper().Log()
	_ = page_monolit.GetPageMonolit()
	_ = page_modules.GetPageModules()
	_ = page_module.GetPageModule()

	_ = http_api.NewHttpApi()
	return sf
}

// Run -- запускает модуль в работу
func (sf *ModuleKeeper) Run() {
	sf.log.Info("ModuleKernelCtx.Run(): module=%v, is run", sf.Name())
	go sf.kServHttp.Run()
}

// IsWork -- признак работы модуля
func (sf *ModuleKeeper) IsWork() bool {
	return sf.kCtx.Wg().IsWork()
}
