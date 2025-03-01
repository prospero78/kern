package page_module

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/krn/kmodule"
	"github.com/prospero78/kern/krn/kmonolit"
	"github.com/prospero78/kern/krn/kserv_http"
	. "github.com/prospero78/kern/krn/ktypes"
	"github.com/prospero78/kern/mock/mock_env"
)

type tester struct {
	t    *testing.T
	ctx  IKernelCtx
	serv IKernelServerHttp
	page *PageModule
}

func TestPageMonolit(t *testing.T) {
	sf := &tester{
		t:   t,
		ctx: kctx.GetKernelCtx(),
	}
	sf.new()
	sf.postModuleState()
	sf.postMonolitCtx()
	sf.postModuleLog()
	sf.postModuleState1()
	sf.postMonolitCtx()
	sf.postModuleLog()
	sf.done()
}

// Возвращает состояние модуля 1 (теперь добавлен)
func (sf *tester) postModuleState1() {
	sf.t.Log("postModuleState1")
	mon := kmonolit.GetMonolit("test_monolit")
	ctxMon := mon.Ctx()
	module := kmodule.NewKernelModule("test_module")
	module.Log().Debug("test msg")
	module.Log().Debug("test msg")
	ctxMod := module.Ctx()
	ctxMod.Set("demo_key", "demo value", "for demo comment")
	time.Sleep(time.Millisecond * 20)
	ctxMon.Set("module_1", module, "test_module")
	fiberApp := sf.serv.Fiber()
	req, err := http.NewRequest("POST", "/module/1", nil)
	if err != nil {
		sf.t.Fatalf("postModuleState1(): in net request, err=%v", err)
	}
	resp, err := fiberApp.Test(req)
	if err != nil {
		sf.t.Fatalf("postModuleState1(): in make POST, err=%v", err)
	}
	if resp.StatusCode != 200 {
		sf.t.Fatalf("postModuleState1(): status(%v)!=200", resp.StatusCode)
	}
}

// Возвращает состояние лога модуля
func (sf *tester) postModuleLog() {
	sf.t.Log("postMonolitCtx")
	mon := sf.ctx.Get("monolit").Val().(IKernelMonolit)
	log := mon.Log()
	log.Debug("test msg")
	fiberApp := sf.serv.Fiber()
	req, err := http.NewRequest("POST", "/module_log/1", nil)
	if err != nil {
		sf.t.Fatalf("postModuleLog(): in net request, err=%v", err)
	}
	resp, err := fiberApp.Test(req)
	if err != nil {
		sf.t.Fatalf("postModuleLog(): in make POST, err=%v", err)
	}
	if resp.StatusCode != 200 {
		sf.t.Fatalf("postModuleLog(): status(%v)!=200", resp.StatusCode)
	}
}

// Возвращает состояние контекста модуля
func (sf *tester) postMonolitCtx() {
	sf.t.Log("postMonolitCtx")
	fiberApp := sf.serv.Fiber()
	req, err := http.NewRequest("POST", "/module_ctx/1", nil)
	if err != nil {
		sf.t.Fatalf("postMonolitCtx(): in net request, err=%v", err)
	}
	resp, err := fiberApp.Test(req)
	if err != nil {
		sf.t.Fatalf("postMonolitCtx(): in make POST, err=%v", err)
	}
	if resp.StatusCode != 200 {
		sf.t.Fatalf("postMonolitCtx(): status(%v)!=200", resp.StatusCode)
	}
}

// Возвращает состояние модуля
func (sf *tester) postModuleState() {
	sf.t.Log("postModuleState")
	fiberApp := sf.serv.Fiber()
	req, err := http.NewRequest("POST", "/module/1", nil)
	if err != nil {
		sf.t.Fatalf("postModuleState(): in net request, err=%v", err)
	}
	resp, err := fiberApp.Test(req)
	if err != nil {
		sf.t.Fatalf("postModuleState(): in make POST, err=%v", err)
	}
	if resp.StatusCode != 200 {
		sf.t.Fatalf("postModuleState(): status(%v)!=200", resp.StatusCode)
	}
}

// Освобождает ресурсы
func (sf *tester) done() {
	sf.t.Log("done")
	sf.ctx.Cancel()
	sf.ctx.Wg().Wait()
}

// Создаёт новую страницу модуля
func (sf *tester) new() {
	sf.t.Log("new")
	_ = mock_env.MakeEnv()
	_ = os.Unsetenv("LOCAL_HTTP_URL")
	os.Setenv("LOCAL_HTTP_URL", "http://localhost:18318/")
	sf.ctx.Set("isLocal", true, "testing")
	_ = kmonolit.GetMonolit("test_monolit")
	sf.serv = kserv_http.GetKernelServHttp()

	sf.page = GetPageModule()
	if sf.page == nil {
		sf.t.Fatalf("new(): page==nil")
	}
	_ = GetPageModule()
	go sf.serv.Run()
}
