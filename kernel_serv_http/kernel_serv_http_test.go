package kernel_serv_http

import (
	"testing"
	"time"

	"github.com/svi/kern/kernel_ctx"
	. "github.com/svi/kern/kernel_types"

	// . "github.com/svi/kern/helpers"
	"github.com/svi/kern/mock/mock_env"
)

type tester struct {
	t   *testing.T
	ctx IKernelCtx
	wg  IKernelWg
	me  *mock_env.MockEnv
}

func TestKernelServHttp(t *testing.T) {
	ctx := kernel_ctx.GetKernelCtx()
	sf := &tester{
		t:   t,
		ctx: ctx,
		wg:  ctx.Wg(),
	}
	sf.new()
	sf.close()
}

// Закрытие HTTP-сервера
func (sf *tester) close() {
	sf.t.Log("close")
	sf.ctx.Cancel()
	sf.wg.Wait()
	kernServHttp.close()
}

// Создание сервера HTTP
func (sf *tester) new() {
	sf.t.Log("new()")
	sf.newBad1()
	sf.newBad2()
	sf.newGood1()
	sf.newBad3()
}

// Повторный запуск сервера на том же порту
func (sf *tester) newBad3() {
	sf.t.Log("newBad3()")
	ctx := kernel_ctx.GetKernelCtx()
	serv := GetKernelServHttp(ctx)
	go serv.Run()
	time.Sleep(time.Millisecond * 20)
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1()")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("newGood1(): panic=%v", _panic)
		}
	}()
	sf.me = mock_env.MakeEnv()
	serv := GetKernelServHttp(sf.ctx)
	if serv != kernServHttp {
		sf.t.Fatalf("newGood1(): bad IKernelServHttp")
	}
	if webFiber := serv.Fiber(); webFiber != kernServHttp.fiberApp {
		sf.t.Fatalf("newGood1(): webFiber==serv.appFiber")
	}
	go serv.Run()
}

// Не указана SERVER_HTTP_PORT
func (sf *tester) newBad2() {
	sf.t.Log("newBad2()")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad2(): panic==nil")
		}
	}()
	_ = GetKernelServHttp(sf.ctx)
}

// Нет контекста ядра
func (sf *tester) newBad1() {
	sf.t.Log("newBad1()")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	var ctx IKernelCtx
	_ = GetKernelServHttp(ctx)
}
