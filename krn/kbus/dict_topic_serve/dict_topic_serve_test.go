package dict_topic_serve

import (
	"testing"

	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/mock/mock_hand_serve"
)

type tester struct {
	t    *testing.T
	dict *dictServe
	hand *mock_hand_serve.MockHandlerServe
}

func TestDictSub(t *testing.T) {
	sf := &tester{
		t:    t,
		hand: mock_hand_serve.NewMockHandlerServe("topic_dict_serve", "name_dict_serve"),
	}
	sf.new()
	sf.addBad1()
	sf.addGood1()
	sf.addBad2()
	sf.sendBad1()
	sf.sendBad2()
	sf.sendGood1()
	sf.delBad1()
	sf.delGood2()
	sf.callBad3()
}

// Работа ядра завершена
func (sf *tester) callBad3() {
	sf.t.Log("callBad3")
	ctx := kctx.GetKernelCtx()
	ctx.Cancel()
	ctx.Wg().Wait()
	sf.dict.Register(sf.hand)
	binMsg, err := sf.dict.SendRequest(sf.hand.Topic(), []byte("test"))
	if err == nil {
		sf.t.Fatalf("callBad3(): err==nil")
	}
	if binMsg != nil {
		sf.t.Fatalf("callBad3(): binMsg!=nil")
	}
}

func (sf *tester) delGood2() {
	sf.t.Log("delGood2()")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("delGood2(): panic=%v", _panic)
		}
	}()
	sf.dict.Unregister(sf.hand)
	sf.dict.Unregister(sf.hand)
}

// Удаляет, чего нет
func (sf *tester) delBad1() {
	sf.t.Log("delBad1()")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("delBad1(): panic==nil")
		}
	}()
	sf.dict.Unregister(nil)
}

func (sf *tester) sendGood1() {
	sf.t.Log("sendGood1")
	TimeoutDefault = 5000
	binMsg, err := sf.dict.SendRequest(sf.hand.Topic(), []byte("test_good"))
	if err != nil {
		sf.t.Fatalf("sendGood1(): err=%v", err)
	}
	if binMsg == nil {
		sf.t.Fatalf("sendGood1(): binMsg==nil")
	}
}

// Обработчик вернул ошибку
func (sf *tester) sendBad2() {
	sf.t.Log("sendBad2")
	sf.hand.IsBad_.Set()
	binMsg, err := sf.dict.SendRequest(sf.hand.Topic(), []byte("test"))
	if err == nil {
		sf.t.Fatalf("sendBad2(): err==nil")
	}
	if binMsg != nil {
		sf.t.Fatalf("sendBad2(): binMsg!=nil")
	}
	sf.hand.IsBad_.Reset()
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
	sf.dict.Register(sf.hand)
}

// Правильное добавление обработчика подписки
func (sf *tester) addGood1() {
	sf.t.Log("addGood1()")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("addGood1(): panic=%v", _panic)
		}
	}()
	sf.dict.Register(sf.hand)
}

// Вместо обработчика пустышка
func (sf *tester) addBad1() {
	sf.t.Log("addBad1()")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("addBad1(): panic==nil")
		}
	}()
	sf.dict.Register(nil)
}

// Вызов несуществующего топика
func (sf *tester) sendBad1() {
	sf.t.Log("sendBad1")
	binRes, err := sf.dict.SendRequest("test_bad_topic", []byte("test"))
	if err == nil {
		sf.t.Fatalf("sendBad1(): err==nil")
	}
	if binRes != nil {
		sf.t.Fatalf("sendBad1(): binRes!=nil")
	}
}

// Создание словаря подписчиков
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("newGood1(): panic=%v", _panic)
		}
	}()
	sf.dict = NewDictServe().(*dictServe)
	if sf.dict == nil {
		sf.t.Fatalf("newGood1(): DictServe==nil")
	}
}
