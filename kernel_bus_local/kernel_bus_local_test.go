package kernel_bus_local

import (
	"fmt"
	"testing"
	"time"

	. "github.com/prospero78/kern/kernel_alias"
	"github.com/prospero78/kern/kernel_ctx"
	. "github.com/prospero78/kern/kernel_types"
	"github.com/prospero78/kern/safe_bool"
)

type tester struct {
	t        *testing.T
	ctx      IKernelCtx
	bus      IKernelBus
	handSub  *handSub
	handServ *handServ
}

func TestKernelBusLocal(t *testing.T) {
	sf := &tester{
		t:   t,
		ctx: kernel_ctx.GetKernelCtx(),
		handSub: &handSub{
			t:     t,
			chMsg: make(chan []byte, 1),
		},
		handServ: &handServ{
			t:      t,
			chMsg:  make(chan []byte, 2),
			isLong: safe_bool.NewSafeBool(),
		},
	}
	sf.new()
	sf.subBad1()
	sf.subGood1()
	sf.pubBad1()
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
	binMsg, err := sf.bus.Request(sf.ctx.Ctx(), "test_topic", []byte("test_msg"))
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
	sf.bus.Serve(sf.handServ)
}

// Нет обработчика для обслуживания запросов
func (sf *tester) servBad1() {
	sf.t.Log("servBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("servBad1(): panic==nil")
		}
	}()
	sf.bus.Serve(nil)
}

// Нет такого топика
func (sf *tester) reqBad1() {
	sf.t.Log("reqBad1")
	binMsg, err := sf.bus.Request(sf.ctx.Ctx(), "test_topic1", []byte("test_msg"))
	if err == nil {
		sf.t.Fatalf("reqBad1(): err==nil")
	}
	if binMsg != nil {
		sf.t.Fatalf("reqBad1(): binMsg!=nil")
	}
}

// Нет такого топика
func (sf *tester) pubBad1() {
	sf.t.Log("pubBad1")
	err := sf.bus.Publish(sf.ctx.Ctx(), "test_topic1", []byte("test_msg"))
	if err != nil {
		sf.t.Fatalf("pubBad1(): err=%v", nil)
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
	err := sf.bus.Subscribe(nil)
	if err == nil {
		sf.t.Fatalf("subBad1(): err==nil")
	}
}

// Создание локальной шины
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()

}

// Закрытие шины
func (sf *tester) close() {
	sf.t.Log("close")
	sf.ctx.Cancel()
	sf.ctx.Wg().Wait()
	sf.bus.(*kernelBusLocal).close()
	if sf.bus.IsWork() {
		sf.t.Fatalf("close(): bus work")
	}
	err := sf.bus.Subscribe(sf.handSub)
	if err == nil {
		sf.t.Fatalf("close(): err==nil")
	}
	err = sf.bus.Publish(sf.ctx.Ctx(), "test_topic1", []byte("test_msg"))
	if err == nil {
		sf.t.Fatalf("close(): err==nil")
	}
	_, err = sf.bus.Request(sf.ctx.Ctx(), "test_topic", []byte("test_msg"))
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
	sf.bus = GetKernelBusLocal(sf.ctx)
	sf.bus = GetKernelBusLocal(sf.ctx)
	if sf.bus == nil {
		sf.t.Fatalf("newGood1(): IKernelBus==nil")
	}
	if !sf.bus.IsWork() {
		sf.t.Fatalf("newGood1(): bus not work")
	}
}

// Нет контекста ядра
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	_ = GetKernelBusLocal(nil)
}

type handSub struct {
	t     *testing.T
	chMsg chan []byte // Для обратного вызова
}

// Функция обратного вызова подписки
func (sf *handSub) FnBack(binMsg []byte) {
	sf.t.Log("FnBack")
	sf.chMsg <- binMsg
}

// Возвращает топик для обработчика подписки
func (sf *handSub) Topic() ATopic {
	return "test_topic"
}

type handServ struct {
	t      *testing.T
	chMsg  chan []byte // Для обратного вызова
	isBad  bool        // Признак сбоя при вызове
	isLong ISafeBool   // Долгое выполнение вызова
}

// Функция обратного вызова подписки
func (sf *handServ) FnBack(binMsg []byte) ([]byte, error) {
	sf.t.Log("FnBack")
	if sf.isBad {
		return nil, fmt.Errorf("FnBack(): isBad==true")
	}
	if sf.isLong.Get() {
		time.Sleep(time.Millisecond * 20)
	}
	sf.chMsg <- binMsg
	return []byte("response"), nil
}

// Возвращает топик для обработчика подписки
func (sf *handServ) Topic() ATopic {
	return "test_topic"
}
