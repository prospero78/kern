package handler_http_sub

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/prospero78/kern/kernel_ctx"
	"github.com/prospero78/kern/kernel_serv_http"
	"github.com/prospero78/kern/mock/mock_env"
)

type tester struct {
	t     *testing.T
	isBad bool // Признак испорченности обработчика
	hand  *handlerHttpSub
}

func TestHandlerHttpSub(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.back()
	ctx := kernel_ctx.GetKernelCtx()
	ctx.Cancel()
	ctx.Wg().Wait()
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
	kernServ := kernel_serv_http.GetKernelServHttp()
	fiberApp := kernServ.Fiber()
	sf.hand.name = "http://localhost:18200/test/local"
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
	sf.hand = NewHandlerHttpSub("test_topic", "/test/local").(*handlerHttpSub)
	if sf.hand == nil {
		sf.t.Fatalf("newGood1(): handler==nil")
	}
	if name := sf.hand.Name(); name != "/test/local" {
		sf.t.Fatalf("newGood1(): name(%v)!='/test/local'", name)
	}
	if topic := sf.hand.Topic(); topic != "test_topic" {
		sf.t.Fatalf("newGood1(): topic(%v)!='test_topic'", topic)
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
	_ = NewHandlerHttpSub("", "/test/local")
}
