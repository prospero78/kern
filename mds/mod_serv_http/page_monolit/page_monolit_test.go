package page_monolit

import (
	"net/http"
	"os"
	"testing"

	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/krn/kmonolit"
	"github.com/prospero78/kern/krn/kserv_http"
	. "github.com/prospero78/kern/krn/ktypes"
	"github.com/prospero78/kern/mock/mock_env"
)

type tester struct {
	t    *testing.T
	ctx  IKernelCtx
	serv IKernelServerHttp
	page *PageMonolit
}

func TestPageMonolit(t *testing.T) {
	sf := &tester{
		t:   t,
		ctx: kctx.GetKernelCtx(),
	}
	sf.new()
	sf.getMonolit()
	sf.postMonolitState()
	sf.postMonolitCtx()
	sf.postMonolitLog()
	sf.done()
}

// Возвращает состояние лога монолита
func (sf *tester) postMonolitLog() {
	sf.t.Log("postMonolitCtx")
	mon := sf.ctx.Get("monolit").Val().(IKernelMonolit)
	log := mon.Log()
	log.Debug("test msg")
	fiberApp := sf.serv.Fiber()
	req, err := http.NewRequest("POST", "/monolit_log", nil)
	if err != nil {
		sf.t.Fatalf("postMonolitLog(): in net request, err=%v", err)
	}
	resp, err := fiberApp.Test(req)
	if err != nil {
		sf.t.Fatalf("postMonolitLog(): in make POST, err=%v", err)
	}
	if resp.StatusCode != 200 {
		sf.t.Fatalf("postMonolitLog(): status(%v)!=200", resp.StatusCode)
	}
}

// Возвращает состояние контекста монолита
func (sf *tester) postMonolitCtx() {
	sf.t.Log("postMonolitCtx")
	fiberApp := sf.serv.Fiber()
	req, err := http.NewRequest("POST", "/monolit_ctx", nil)
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

// Возвращает состояние монолита
func (sf *tester) postMonolitState() {
	sf.t.Log("postMonolitState")
	fiberApp := sf.serv.Fiber()
	req, err := http.NewRequest("POST", "/monolit_state", nil)
	if err != nil {
		sf.t.Fatalf("postMonolitState(): in net request, err=%v", err)
	}
	resp, err := fiberApp.Test(req)
	if err != nil {
		sf.t.Fatalf("postMonolitState(): in make POST, err=%v", err)
	}
	if resp.StatusCode != 200 {
		sf.t.Fatalf("postMonolitState(): status(%v)!=200", resp.StatusCode)
	}
}

// Возвращает главную страницу монолита
func (sf *tester) getMonolit() {
	sf.t.Log("getMonolit")
	fiberApp := sf.serv.Fiber()
	req, err := http.NewRequest("GET", "/monolit", nil)
	if err != nil {
		sf.t.Fatalf("getMonolit(): in net request, err=%v", err)
	}
	resp, err := fiberApp.Test(req)
	if err != nil {
		sf.t.Fatalf("getMonolit(): in make GET, err=%v", err)
	}
	if resp.StatusCode != 200 {
		sf.t.Fatalf("getMonolit(): status(%v)!=200", resp.StatusCode)
	}
}

// Освобождает ресурсы
func (sf *tester) done() {
	sf.t.Log("done")
	sf.ctx.Cancel()
	sf.ctx.Wg().Wait()
}

// Создаёт новую страницу монолита
func (sf *tester) new() {
	sf.t.Log("new")
	_ = mock_env.MakeEnv()
	_ = os.Unsetenv("LOCAL_HTTP_URL")
	os.Setenv("LOCAL_HTTP_URL", "http://localhost:18310/")
	sf.ctx.Set("isLocal", true, "testing")
	_ = kmonolit.GetMonolit("test_monolit")
	sf.serv = kserv_http.GetKernelServHttp()

	sf.page = GetPageMonolit()
	if sf.page == nil {
		sf.t.Fatalf("new(): page==nil")
	}
	_ = GetPageMonolit()
	go sf.serv.Run()
}
