package client_bus_http

import (
	"os"
	"strings"
	"testing"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_base"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	"gitp78su.ipnodns.ru/svi/kern/krn/kserv_http"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/mock/mock_env"
	"gitp78su.ipnodns.ru/svi/kern/mock/mock_hand_serve"
	"gitp78su.ipnodns.ru/svi/kern/mock/mock_hand_sub_http"
)

type tester struct {
	t        *testing.T
	ctx      IKernelCtx
	cl       *ClientBusHttp
	handSub  *mock_hand_sub_http.MockHandSubHttp
	handServ *mock_hand_serve.MockHandlerServe
	bus      *kbus_base.KBusBase
}

func TestClientBusHttp(t *testing.T) {
	handSub := mock_hand_sub_http.NewMockHandSubHttp("test_topic_sub", "http://localhost:18314/").(*mock_hand_sub_http.MockHandSubHttp)
	handServ := mock_hand_serve.NewMockHandlerServe("test_topic_serv", "local_hook")
	sf := &tester{
		t:        t,
		ctx:      kctx.GetKernelCtx(),
		handSub:  handSub,
		handServ: handServ,
		bus:      kbus_base.GetKernelBusBase(),
	}
	sf.new()
	sf.unsub()
	sf.sub()
	sf.pub()
	sf.unsubGood1()
	sf.reg()
	sf.send()
}

// Отправка запросов
func (sf *tester) send() {
	sf.t.Log("send")
	sf.sendBad1()
	sf.sendBad2()
	sf.sendGood1()
}

func (sf *tester) sendGood1() {
	sf.t.Log("sendGood1")
	binResp, err := sf.cl.SendRequest("test_topic_serv", []byte("test msg 456"))
	if err != nil {
		sf.t.Fatalf("sendGood1(): err=%v", err)
	}
	if binResp == nil {
		sf.t.Fatalf("sendGood1(): binResp==nil")
	}
	strResp := string(binResp)
	if strResp != "test msg 456" {
		sf.t.Fatalf("sendGood1(): strResp(%v)!='test msg 456'", strResp)
	}
}

// Левый адрес
func (sf *tester) sendBad2() {
	sf.t.Log("sendBad2")
	urlRemote := sf.cl.urlRemote
	sf.cl.urlRemote = "tra-ta-ta"
	defer func() {
		sf.cl.urlRemote = urlRemote
	}()
	binResp, err := sf.cl.SendRequest("test_topic_serv", []byte("test msg"))
	if err == nil {
		sf.t.Fatalf("sendBad2(): err==nil")
	}
	if binResp != nil {
		sf.t.Fatalf("sendBad2(): binResp!=nil")
	}
}

// Нет такого топика
func (sf *tester) sendBad1() {
	sf.t.Log("sendBad1")
	binResp, err := sf.cl.SendRequest("test_bad_topic", []byte("test msg"))
	if err == nil {
		sf.t.Fatalf("sendBad1(): err==nil")
	}
	if binResp != nil {
		sf.t.Fatalf("sendBad1(): binResp!=nil")
	}
}

// Регистрация серверного обработчика
func (sf *tester) reg() {
	sf.t.Log("reg")
	sf.regBad1()
	sf.regGood1()
}

func (sf *tester) regGood1() {
	sf.t.Log("regGood1")
	sf.cl.RegisterServe(sf.handServ)
}

// Нет серверного обработчика
func (sf *tester) regBad1() {
	sf.t.Log("regBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("regBad1(): panic==nil")
		}
	}()
	sf.cl.RegisterServe(nil)
}

// Правильная отписка
func (sf *tester) unsubGood1() {
	sf.t.Log("unsubGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("unsubGood1(): panic=%v", _panic)
		}
	}()
	sf.cl.Unsubscribe(sf.handSub)
}

// Публикация сообщения
func (sf *tester) pub() {
	sf.t.Log("pub")
	sf.pubBad1()
	sf.pubGood1()
	sf.pubBad2()
}

// С шиной что-то случилось
func (sf *tester) pubBad2() {
	sf.t.Log("pubBad2")
	sf.bus.IsWork_.Reset()
	defer sf.bus.IsWork_.Set()
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("pubBad2(): panic=%v", _panic)
		}
	}()
	sf.handSub.BinMsg_ = []byte{}
	err := sf.cl.Publish("test_topic_sub", []byte("test_msg_456"))
	if err != nil {
		if strings.Contains(err.Error(), "topic='test_topic_sub',bus already closed") {
			return
		}
		sf.t.Fatalf("pubBad2(): err=%v", err)
	}
	for {
		SleepMs()
		if sf.handSub.Msg() != "" {
			break
		}
	}
	if msg := sf.handSub.Msg(); msg != "test_msg_456" {
		sf.t.Fatalf("bad msg(%v)", msg)
	}
}

func (sf *tester) pubGood1() {
	sf.t.Log("pubGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("pubGood1(): panic=%v", _panic)
		}
	}()
	sf.handSub.BinMsg_ = []byte{}
	err := sf.cl.Publish("test_topic_sub", []byte("test_msg_456"))
	if err != nil {
		if strings.Contains(err.Error(), "topic='test_topic_sub',bus already closed") {
			return
		}
		sf.t.Fatalf("pubGood1(): err=%v", err)
	}
	for {
		SleepMs()
		if sf.handSub.Msg() != "" {
			break
		}
	}
	if msg := sf.handSub.Msg(); msg != "test_msg_456" {
		sf.t.Fatalf("bad msg(%v)", msg)
	}
}

// Левый адрес
func (sf *tester) pubBad1() {
	sf.t.Log("pubBad1")
	urlRemote := sf.cl.urlRemote
	sf.cl.urlRemote = "tra-ta-ta"
	defer func() {
		sf.cl.urlRemote = urlRemote
	}()
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("pubBad1(): panic=%v", _panic)
		}
	}()
	err := sf.cl.Publish("test_topic", []byte("test_msg"))
	if err == nil {
		sf.t.Fatalf("pubBad1(): err==nil")
	}
}

// Подписка на топик
func (sf *tester) sub() {
	sf.t.Log("sub")
	sf.subBad1()
	sf.subBad2()
	sf.subBad3()
	sf.subGood1()
}

// С шиной что-то случилось
func (sf *tester) subBad3() {
	sf.t.Log("subBad3")
	bus := kbus_base.Bus_
	bus.IsWork_.Reset()
	defer bus.IsWork_.Set()
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("subBad3(): panic=%v", _panic)
		}
	}()
	err := sf.cl.Subscribe(sf.handSub)
	if err != nil {
		if strings.Contains(err.Error(), "bus already closed") {
			return
		}
		sf.t.Fatalf("subBad3(): err=%v", err)
	}
}

func (sf *tester) subGood1() {
	sf.t.Log("subGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("subGood1(): panic=%v", _panic)
		}
	}()
	err := sf.cl.Subscribe(sf.handSub)
	if err != nil {
		if strings.Contains(err.Error(), "bus already closed") {
			return
		}
		sf.t.Fatalf("subGood1(): err=%v", err)
	}
}

// Левый адрес
func (sf *tester) subBad2() {
	sf.t.Log("subBad2")
	urlRemote := sf.cl.urlRemote
	sf.cl.urlRemote = "tra-ta-ta"
	defer func() {
		sf.cl.urlRemote = urlRemote
	}()
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("subBad2(): panic=%v", _panic)
		}
	}()
	err := sf.cl.Subscribe(sf.handSub)
	if err == nil {
		sf.t.Fatalf("subBad2(): err==nil")
	}
}

// Нет обработчика
func (sf *tester) subBad1() {
	sf.t.Log("subBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("subBad1(): panic==nil")
		}
	}()
	_ = sf.cl.Subscribe(nil)
}

// Отписка от топика
func (sf *tester) unsub() {
	sf.t.Log("unsub")
	sf.unsubBad1()
	sf.unsubBad2()
	sf.unsubBad3()
}

// Левый адрес
func (sf *tester) unsubBad3() {
	sf.t.Log("unsubBad3")
	urlRemote := sf.cl.urlRemote
	sf.cl.urlRemote = "tra-ta-ta"
	defer func() {
		sf.cl.urlRemote = urlRemote
	}()
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("unsubBad3(): panic=%v", _panic)
		}
	}()
	sf.cl.Unsubscribe(sf.handSub)
}

// Нет такой подписки
func (sf *tester) unsubBad2() {
	sf.t.Log("unsubBad2")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("unsubBad2(): panic=%v", _panic)
		}
	}()
	sf.cl.Unsubscribe(sf.handSub)
}

// Нет обработчика
func (sf *tester) unsubBad1() {
	sf.t.Log("unsubBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("unsubBad1(): panic==nil")
		}
	}()
	sf.cl.Unsubscribe(nil)
}

// Создание нового клиента
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newBad2()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("newGood1(): panic=%v", _panic)
		}
	}()
	_ = mock_env.MakeEnv()
	_ = os.Unsetenv("LOCAL_HTTP_URL")
	_ = os.Setenv("LOCAL_HTTP_URL", "http://localhost:18314/")
	sf.ctx.Set("monolitName", "test_monolit", "comment")
	sf.cl = NewClientBusHttp("http://localhost:18314/").(*ClientBusHttp)
	kServHttp := kserv_http.GetKernelServHttp()
	go kServHttp.Run()
	for {
		SleepMs()
		if kServHttp.IsWork() {
			break
		}
	}
	if log := sf.cl.Log(); log == nil {
		sf.t.Fatalf("newGood1(): log==nil")
	}
	if isWork := sf.cl.IsWork(); !isWork {
		sf.t.Fatalf("newGood1(): isWork==false")
	}
}

func (sf *tester) newBad2() {
	sf.t.Log("newBad2")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad2(): panic==nil")
		}
	}()
	_ = mock_env.MakeEnv()
	_ = NewClientBusHttp("")
}

// Нет ничего для HTTP-шины
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	_ = NewClientBusHttp("url")
}
