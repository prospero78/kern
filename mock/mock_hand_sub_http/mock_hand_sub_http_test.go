package mock_hand_sub_http

import (
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	"gitp78su.ipnodns.ru/svi/kern/krn/kserv_http"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/mock/mock_env"
)

type tester struct {
	t     *testing.T
	ctx   IKernelCtx
	isBad bool // Признак испорченности обработчика
	hand  *MockHandSubHttp
}

func TestHandlerHttpSub(t *testing.T) {
	sf := &tester{
		t:   t,
		ctx: kctx.GetKernelCtx(),
	}
	sf.ctx.Set("monolitName", "test_monolit", "comment")
	sf.new()
	sf.back()
	sf.ctx.Cancel()
	sf.ctx.Wg().Wait()
}

// Проверка работы обратного вызова
func (sf *tester) back() {
	sf.t.Log("back")
	sf.backBad1()
	sf.backBad2()
	sf.backGood1()
}
func (sf *tester) backGood1() {
	sf.t.Log("backGood1")
	sf.hand.FnBack([]byte("test_msg"))
}

// Что-то случилось на той стороне
func (sf *tester) backBad2() {
	sf.t.Log("backBad2")
	_ = mock_env.MakeEnv()
	kernServ := kserv_http.GetKernelServHttp()
	fiberApp := kernServ.Fiber()
	sf.hand.WebHook_ = "http://localhost:18200/test/local"
	fiberApp.Post("/test/local", sf.testLocal)
	go kernServ.Run()
	sf.isBad = true
	sf.hand.FnBack([]byte("test_msg"))
	sf.isBad = false
}

// Эндпоинт на HTTP-сервере
func (sf *tester) testLocal(ctx *fiber.Ctx) error {
	if sf.isBad {
		return ctx.SendStatus(400)
	}
	return ctx.SendString("ok")
}

// Нет обработчика для запроса
func (sf *tester) backBad1() {
	sf.t.Log("backBad1")
	sf.hand.FnBack([]byte("hello_test"))
}

// Создание HTTP-обработчика подписки
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("newGood1(): panic=%v", _panic)
		}
	}()
	sf.hand = NewMockHandSubHttp("test_topic", "/test/local").(*MockHandSubHttp)
	if sf.hand == nil {
		sf.t.Fatalf("newGood1(): handler==nil")
	}
	if name := sf.hand.Name(); !strings.Contains(string(name), "/test/local_") {
		sf.t.Fatalf("newGood1(): name(%v)!='/test/local_'", name)
	}
	sf.hand.SetName("test_name")
	if name := sf.hand.Name(); name != "test_name" {
		sf.t.Fatalf("newGood1(): name(%v)!='test_name'", name)
	}
	if topic := sf.hand.Topic(); topic != "test_topic" {
		sf.t.Fatalf("newGood1(): bad topic(%v) 'test_topic'", topic)
	}
	if msg := sf.hand.Msg(); msg != "" {
		sf.t.Fatalf("newGood1(): msg(%v) not empty", msg)
	}
}

// Нет нужных полей
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	_ = NewMockHandSubHttp("", "/test/local")
}
