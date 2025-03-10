package mod_serv_http

import (
	"os"
	"testing"
	"time"

	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/mock/mock_env"
)

type tester struct {
	t   *testing.T
	ctx IKernelCtx
	mod IKernelModule
}

func TestModServHttp(t *testing.T) {
	sf := &tester{
		t:   t,
		ctx: kctx.GetKernelCtx(),
	}
	sf.ctx.Set("monolitName", "test_monolit", "test")
	sf.new()
	sf.done()
}

// Завершение работы
func (sf *tester) done() {
	sf.t.Log("done")
	sf.ctx.Cancel()
	sf.ctx.Wg().Wait()
	if isWork := sf.mod.IsWork(); isWork {
		sf.t.Fatalf("newGood1(): isWork==true")
	}
}

// Создание нового модуля HTTP-сервера
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	_ = mock_env.MakeEnv()
	_ = os.Unsetenv("LOCAL_HTTP_URL")
	os.Setenv("LOCAL_HTTP_URL", "http://localhost:18301/")
	sf.mod = NewModuleServHttp()
	if sf.mod == nil {
		sf.t.Fatalf("newGood1(): mod==nil")
	}
	go sf.mod.Run()
	for {
		time.Sleep(time.Millisecond * 1)
		if sf.mod.IsWork() {
			return
		}
	}
}

// нет ничего для создания модуля
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	_ = NewModuleServHttp()
}
