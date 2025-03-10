package page_modules

import (
	"net/http"
	"os"
	"testing"
	"time"

	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	"gitp78su.ipnodns.ru/svi/kern/krn/kmodule"
	"gitp78su.ipnodns.ru/svi/kern/krn/kmonolit"
	"gitp78su.ipnodns.ru/svi/kern/krn/kserv_http"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/mock/mock_env"
)

type tester struct {
	t    *testing.T
	ctx  IKernelCtx
	serv IKernelServerHttp
	page *PageModules
}

func TestPageMonolit(t *testing.T) {
	sf := &tester{
		t:   t,
		ctx: kctx.GetKernelCtx(),
	}
	sf.new()
	sf.postModulesState()
	sf.postModulesState1()
	sf.done()
}

// Возвращает состояние модуля 1 (теперь добавлен)
func (sf *tester) postModulesState1() {
	sf.t.Log("postModulesState1")
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
	req, err := http.NewRequest("POST", "/modules", nil)
	if err != nil {
		sf.t.Fatalf("postModulesState1(): in net request, err=%v", err)
	}
	resp, err := fiberApp.Test(req)
	if err != nil {
		sf.t.Fatalf("postModulesState1(): in make POST, err=%v", err)
	}
	if resp.StatusCode != 200 {
		sf.t.Fatalf("postModulesState1(): status(%v)!=200", resp.StatusCode)
	}
}

// Возвращает состояние модуля
func (sf *tester) postModulesState() {
	sf.t.Log("postModulesState")
	fiberApp := sf.serv.Fiber()
	req, err := http.NewRequest("POST", "/modules", nil)
	if err != nil {
		sf.t.Fatalf("postModulesState(): in net request, err=%v", err)
	}
	resp, err := fiberApp.Test(req)
	if err != nil {
		sf.t.Fatalf("postModulesState(): in make POST, err=%v", err)
	}
	if resp.StatusCode != 200 {
		sf.t.Fatalf("postModulesState(): status(%v)!=200", resp.StatusCode)
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
	os.Setenv("LOCAL_HTTP_URL", "http://localhost:18315/")
	sf.ctx.Set("isLocal", true, "testing")
	_ = kmonolit.GetMonolit("test_monolit")
	sf.serv = kserv_http.GetKernelServHttp()

	sf.page = GetPageModules()
	if sf.page == nil {
		sf.t.Fatalf("new(): page==nil")
	}
	_ = GetPageModules()
	go sf.serv.Run()
}
