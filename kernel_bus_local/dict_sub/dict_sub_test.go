package dict_sub

import (
	"testing"

	. "github.com/prospero78/kern/kernel_alias"
)

type handler struct {
	t     *testing.T
	chMsg chan []byte // Для обратного вызова
}

// Функция обратного вызова подписки
func (sf *handler) FnBack(binMsg []byte) {
	sf.t.Log("FnBack")
	sf.chMsg <- binMsg
}

// Возвращает топик для обработчика подписки
func (sf *handler) Topic() ATopic {
	return "test_topic"
}

type tester struct {
	t    *testing.T
	dict *DictSub
	hand *handler
}

func TestDictSub(t *testing.T) {
	sf := &tester{
		t: t,
		hand: &handler{
			t:     t,
			chMsg: make(chan []byte, 2),
		},
	}
	sf.new()
	sf.addBad1()
	sf.addGood1()
	sf.addBad2()
	sf.callBad1()
	sf.callGood1()
	sf.callBad2()
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

func (sf *tester) callGood1() {
	sf.t.Log("callGood1")
	sf.dict.Call(sf.hand.Topic(), []byte("test_good"))
	<-sf.hand.chMsg
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

// Нет топика
func (sf *tester) callBad2() {
	sf.t.Log("callBad2")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("callBad2(): panic==nil")
		}
	}()
	sf.dict.Call("", []byte("test_msg"))
}

// Вызов несуществующего топика
func (sf *tester) callBad1() {
	sf.t.Log("callBad1")
	sf.dict.Call("test_bad_topic", []byte("test"))
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
	sf.dict = NewDictSub()
	if sf.dict == nil {
		sf.t.Fatalf("newGood1(): DictSub==nil")
	}
}
