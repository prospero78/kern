package dict_serve

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/prospero78/kern/kernel_alias"
	"github.com/prospero78/kern/kernel_ctx"
	. "github.com/prospero78/kern/kernel_types"
	"github.com/prospero78/kern/safe_bool"
)

type tester struct {
	t    *testing.T
	ctx  IKernelCtx
	dict *DictServe
	hand *handler
}

func TestDictSub(t *testing.T) {
	sf := &tester{
		t:   t,
		ctx: kernel_ctx.GetKernelCtx(),
		hand: &handler{
			t:      t,
			chMsg:  make(chan []byte, 2),
			isLong: safe_bool.NewSafeBool(),
		},
	}
	sf.new()
	sf.addBad1()
	sf.addGood1()
	sf.addBad2()
	sf.callBad1()
	sf.callBad2()
	sf.callBad3()
	sf.callGood1()
	sf.callBad4()
	sf.delBad1()
	sf.delGood2()
}

func (sf *tester) delGood2() {
	sf.t.Log("delGood2()")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("delGood2(): panic=%v", _panic)
		}
	}()
	sf.dict.Del(sf.hand)
	sf.dict.Del(sf.hand)
}

// Удаляет, чего нет
func (sf *tester) delBad1() {
	sf.t.Log("delBad1()")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("delBad1(): panic==nil")
		}
	}()
	sf.dict.Del(nil)
}

func (sf *tester) callBad4() {
	sf.t.Log("callBad4")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("callBad4)(: panic==nil)")
		}
	}()
	var ctx context.Context
	_, _ = sf.dict.Call(ctx, sf.hand.Topic(), []byte("test_good"))
}
func (sf *tester) callGood1() {
	sf.t.Log("callGood1")
	TimeoutDefault = 5000
	binMsg, err := sf.dict.Call(sf.ctx.Ctx(), sf.hand.Topic(), []byte("test_good"))
	if err != nil {
		sf.t.Fatalf("callGood1(): err=%v", err)
	}
	if binMsg == nil {
		sf.t.Fatalf("callGood1(): binMsg==nil")
	}
	<-sf.hand.chMsg
}

// Слишком долгое ожидание
func (sf *tester) callBad3() {
	sf.t.Log("callBad3")
	sf.hand.isLong.Set()
	TimeoutDefault = 5
	binMsg, err := sf.dict.Call(sf.ctx.Ctx(), sf.hand.Topic(), []byte("test"))
	if err == nil {
		sf.t.Fatalf("callBad3(): err==nil")
	}
	if binMsg != nil {
		sf.t.Fatalf("callBad3(): binMsg!=nil")
	}
	sf.hand.isLong.Reset()
}

// Обработчик вернул ошибку
func (sf *tester) callBad2() {
	sf.t.Log("callBad2")
	sf.hand.isBad = true
	binMsg, err := sf.dict.Call(sf.ctx.Ctx(), sf.hand.Topic(), []byte("test"))
	if err == nil {
		sf.t.Fatalf("callBad2(): err==nil")
	}
	if binMsg != nil {
		sf.t.Fatalf("callBad2(): binMsg!=nil")
	}
	sf.hand.isBad = false
}

// повторное добавление обработчика
func (sf *tester) addBad2() {
	sf.t.Log("addBad2")
	sf.t.Log("addGood1()")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("addGood1(): panic==nil")
		}
	}()
	sf.dict.Add(sf.hand)
}

// Правильное добавление обработчика подписки
func (sf *tester) addGood1() {
	sf.t.Log("addGood1()")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("addGood1(): panic=%v", _panic)
		}
	}()
	sf.dict.Add(sf.hand)
}

// Вместо обработчика пустышка
func (sf *tester) addBad1() {
	sf.t.Log("addBad1()")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("addBad1(): panic==nil")
		}
	}()
	sf.dict.Add(nil)
}

// Вызов несуществующего топика
func (sf *tester) callBad1() {
	sf.t.Log("callBad1")
	binRes, err := sf.dict.Call(sf.ctx.Ctx(), "test_bad_topic", []byte("test"))
	if err == nil {
		sf.t.Fatalf("callBad1(): err==nil")
	}
	if binRes != nil {
		sf.t.Fatalf("callBad1(): binRes!=nil")
	}
}

// Создание словаря подписчиков
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
	sf.dict = NewDictServe(sf.ctx)
	if sf.dict == nil {
		sf.t.Fatalf("newGood1(): DictServe==nil")
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
	var ctx IKernelCtx
	_ = NewDictServe(ctx)
}

type handler struct {
	t      *testing.T
	chMsg  chan []byte // Для обратного вызова
	isBad  bool        // Признак сбоя при вызове
	isLong ISafeBool   // Долгое выполнение вызова
}

// Функция обратного вызова подписки
func (sf *handler) FnBack(binMsg []byte) ([]byte, error) {
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
func (sf *handler) Topic() ATopic {
	return "test_topic"
}
