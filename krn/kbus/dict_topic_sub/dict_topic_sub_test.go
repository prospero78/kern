package dict_topic_sub

import (
	"testing"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/mock/mock_hand_sub_local"
)

type tester struct {
	t    *testing.T
	dict *dictTopicSub
	hand *mock_hand_sub_local.MockHandlerSub
}

func TestDictSub(t *testing.T) {
	sf := &tester{
		t:    t,
		hand: mock_hand_sub_local.NewMockHandlerSub("topic_dict_sub", "name_dict_sub"),
	}
	sf.new()
	sf.addBad1()
	sf.addGood1()
	sf.addGood2()
	sf.callGood10()
	sf.callGood1()
	sf.callBad2()
	sf.delBad1()
	sf.delGood2()
	sf.unsub1()
}

// Прямой вызов отписки от топика которого точно нет
func (sf *tester) unsub1() {
	sf.t.Log("unsub1")
	sf.dict.Unsubscribe(sf.hand)
	sf.dict.Read("test_test", []byte("test test"))
	hand := mock_hand_sub_local.NewMockHandlerSub("topic_dict_sub1", "name_dict_sub")
	sf.dict.Unsubscribe(hand)
	count := 0
	for count < 100 {
		SleepMs()
		count++
	}
	sf.dict.Read("topic_dict_sub", []byte("test test"))
	count = 0
	for count < 200 {
		SleepMs()
		count++
	}
}

func (sf *tester) delGood2() {
	sf.t.Log("delGood2()")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("delGood2(): panic=%v", _panic)
		}
	}()
	sf.dict.Unsubscribe(sf.hand)
	sf.dict.Unsubscribe(sf.hand)
}

// Удаляет, чего нет
func (sf *tester) delBad1() {
	sf.t.Log("delBad1()")
	sf.dict.Unsubscribe(nil)
}

func (sf *tester) callGood1() {
	sf.t.Log("callGood1")
	sf.dict.Read(sf.hand.Topic(), []byte("test_good"))
}

// повторное добавление обработчика
func (sf *tester) addGood2() {
	sf.t.Log("addGood2")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("addGood2(): panic=%v", _panic)
		}
	}()
	sf.dict.Subscribe(sf.hand)
}

// Правильное добавление обработчика подписки
func (sf *tester) addGood1() {
	sf.t.Log("addGood1()")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("addGood1(): panic=%v", _panic)
		}
	}()
	sf.dict.Subscribe(sf.hand)
}

// Вместо обработчика пустышка
func (sf *tester) addBad1() {
	sf.t.Log("addBad1()")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("addBad1(): panic==nil")
		}
	}()
	sf.dict.Subscribe(nil)
}

// Нет топика
func (sf *tester) callBad2() {
	sf.t.Log("callBad2")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("callBad2(): panic==nil")
		}
	}()
	sf.dict.Read("", []byte("test_msg"))
}

// Нет данных в топике
func (sf *tester) callGood10() {
	sf.t.Log("callGood10")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("callGood10(): panic=%v", _panic)
		}
	}()
	sf.dict.Read("test_bad_topic", []byte("test"))
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
	sf.dict = NewDictTopicSub().(*dictTopicSub)
	if sf.dict == nil {
		sf.t.Fatalf("newGood1(): DictSub==nil")
	}
}
