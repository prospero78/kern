package mock_hand_sub_local

import (
	"strings"
	"testing"
)

type tester struct {
	t    *testing.T
	hand *MockHandlerSub
}

func TestMockHandleServe(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.back()
}

// Проверка обратного вызова
func (sf *tester) back() {
	sf.t.Log("back")
	sf.backGood1()
}

func (sf *tester) backGood1() {
	sf.t.Log("backGood1")
	sf.hand.FnBack([]byte("test_msg"))
	if sf.hand.Msg() != "test_msg" {
		sf.t.Fatalf("backGood1(): binMsg(%v)!='test_msg'", string(sf.hand.Msg_))
	}
}

// Создание мок-обработчика запросов
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	sf.hand = NewMockHandlerSub("test_topic", "test_name")
	if sf.hand == nil {
		sf.t.Fatalf("newGood1(): handler==nil")
	}
	if name := sf.hand.Name(); !strings.Contains(string(name), "test_name_") {
		sf.t.Fatalf("newGood1(): name(%v)!='test_name_'", name)
	}
	if topic := sf.hand.Topic(); topic != "test_topic" {
		sf.t.Fatalf("newGood1(): topic(%v)!='test_topic'", topic)
	}
	if msg := sf.hand.Msg(); msg != "" {
		sf.t.Fatalf("newGood1(): msg not empty")
	}
}

// Нет топика для создания
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("")
		}
	}()
	_ = NewMockHandlerSub("", "test_name")
}
