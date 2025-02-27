package kbus_base

import (
	"testing"

	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
	"github.com/prospero78/kern/mock/mock_hand_serve"
	"github.com/prospero78/kern/mock/mock_hand_sub"
)

type tester struct {
	t        *testing.T
	bus      IKernelBus
	handSub  *mock_hand_sub.MockHandlerSub
	handServ *mock_hand_serve.MockHandlerServe
}

func TestKernelBusLocal(t *testing.T) {
	sf := &tester{
		t:        t,
		handSub:  mock_hand_sub.NewMockHandlerSub("topic_hand_sub", "mock_hand_sub"),
		handServ: mock_hand_serve.NewMockHandlerServe("topic_hand_serv", "mock_hand_serv"),
	}
	sf.new()
	sf.subBad1()
	sf.subGood1()
	sf.pubGood10()
	sf.reqBad1()
	sf.servBad1()
	sf.servGood1()
	sf.reqGood1()
	sf.close()
	sf.unsubBad1()
	sf.unsubGood1()
}

func (sf *tester) unsubGood1() {
	sf.t.Log("unsubBad1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("unsubGood1(): panic=%v", _panic)
		}
	}()
	sf.bus.Unsubscribe(sf.handSub)
	sf.bus.Unsubscribe(sf.handSub)
}

// Отписка от топика, нет обработчика
func (sf *tester) unsubBad1() {
	sf.t.Log("unsubBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("unsubBad1(): panic==nil")
		}
	}()
	sf.bus.Unsubscribe(nil)
}

func (sf *tester) reqGood1() {
	sf.t.Log("reqGood1")
	binMsg, err := sf.bus.SendRequest(sf.handServ.Topic_, []byte("test_msg"))
	if err != nil {
		sf.t.Fatalf("reqGood1(): err=%v", err)
	}
	if binMsg == nil {
		sf.t.Fatalf("reqGood1(): binMsg==nil")
	}
}

func (sf *tester) servGood1() {
	sf.t.Log("servGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("servGood1(): panic=%v", _panic)
		}
	}()
	sf.bus.RegisterServe(sf.handServ)
}

// Нет обработчика для обслуживания запросов
func (sf *tester) servBad1() {
	sf.t.Log("servBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("servBad1(): panic==nil")
		}
	}()
	sf.bus.RegisterServe(nil)
}

// Нет такого топика
func (sf *tester) reqBad1() {
	sf.t.Log("reqBad1")
	binMsg, err := sf.bus.SendRequest("test_topic1", []byte("test_msg"))
	if err == nil {
		sf.t.Fatalf("reqBad1(): err==nil")
	}
	if binMsg != nil {
		sf.t.Fatalf("reqBad1(): binMsg!=nil")
	}
}

// Нет читателей топика
func (sf *tester) pubGood10() {
	sf.t.Log("pubGood10")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("pubGood10(): panic=%v", _panic)
		}
	}()
	err := sf.bus.Publish("test_topic1", []byte("test_msg"))
	if err != nil {
		sf.t.Fatalf("pubGood10(): err=%v", nil)
	}
}

func (sf *tester) subGood1() {
	sf.t.Log("subGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("subGood1(): panic=%v", _panic)
		}
	}()
	err := sf.bus.Subscribe(sf.handSub)
	if err != nil {
		sf.t.Fatalf("subGood1(): err=%v", err)
	}
}

// Нет обработчик подписки
func (sf *tester) subBad1() {
	sf.t.Log("subBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("subBad1(): panic==nil")
		}
	}()
	_ = sf.bus.Subscribe(nil)
}

// Создание локальной шины
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newGood1()

}

// Закрытие шины
func (sf *tester) close() {
	sf.t.Log("close")
	ctx := kctx.GetKernelCtx()
	ctx.Cancel()
	ctx.Wg().Wait()
	sf.bus.(*KernelBusBase).close()
	if sf.bus.IsWork() {
		sf.t.Fatalf("close(): bus work")
	}
	err := sf.bus.Subscribe(sf.handSub)
	if err == nil {
		sf.t.Fatalf("close(): err==nil")
	}
	err = sf.bus.Publish("test_topic1", []byte("test_msg"))
	if err == nil {
		sf.t.Fatalf("close(): err==nil")
	}
	_, err = sf.bus.SendRequest("test_topic", []byte("test_msg"))
	if err == nil {
		sf.t.Fatalf("close(): err==nil")
	}
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("newGood1(): panic=%v", _panic)
		}
	}()
	sf.bus = GetKernelBusBase()
	sf.bus = GetKernelBusBase()
	if sf.bus == nil {
		sf.t.Fatalf("newGood1(): IKernelBus==nil")
	}
	if !sf.bus.IsWork() {
		sf.t.Fatalf("newGood1(): bus not work")
	}
}
