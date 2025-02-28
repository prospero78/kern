// package mod_serv_http -- модуль HTTP-сервера
package mod_serv_http

import (
	"github.com/prospero78/kern/krn/kmodule"
	"github.com/prospero78/kern/krn/kserv_http"
	"github.com/prospero78/kern/krn/kserv_http/http_api"
	"github.com/prospero78/kern/krn/kserv_http/page_module"
	"github.com/prospero78/kern/krn/kserv_http/page_modules"
	"github.com/prospero78/kern/krn/kserv_http/page_monolit"
	. "github.com/prospero78/kern/krn/ktypes"
)

// ModuleServHttp -- модуль HTTP-сервера
type ModuleServHttp struct {
	IKernelModule
	kServHttp IKernelServerHttp
	log       ILogBuf
}

// NewModuleServHttp -- возвращает новый модуль HTTP-сервера
func NewModuleServHttp() *ModuleServHttp {
	sf := &ModuleServHttp{
		IKernelModule: kmodule.NewKernelModule("kServHttp"),
		kServHttp:     kserv_http.GetKernelServHttp(),
	}
	sf.log = sf.Ctx().Log()
	_ = page_monolit.NewPageMonolit()
	_ = page_modules.NewPageModules()
	_ = page_module.NewPageModule()

	_ = http_api.NewHttpApi()
	return sf
}

// Run -- запускает модуль в работу
func (sf *ModuleServHttp) Run() {
	sf.log.Info("ModuleServHttp.Run(): module=%v, is run", sf.Name())
	go sf.kServHttp.Run()
}

// IsWork -- признак работы модуля
func (sf *ModuleServHttp) IsWork() bool {
	return sf.kServHttp.IsWork()
}
