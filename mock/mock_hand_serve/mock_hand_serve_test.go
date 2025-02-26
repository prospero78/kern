package mock_hand_serve

import (
	"testing"
)

type tester struct {
	t    *testing.T
	hand *MockHandlerServe
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
	sf.backBad1()
	sf.backGood1()
}

func (sf *tester) backGood1() {
	sf.t.Log("backGood1")
	binMsg, err := sf.hand.FnBack([]byte("test_msg"))
	if err != nil {
		sf.t.Fatalf("backGood1(): err=%v", err)
	}
	if string(binMsg) != "test_msg" {
		sf.t.Fatalf("backGood1(): binMsg(%v)!='test_msg'", string(binMsg))
	}
}

// Что-то сломалось
func (sf *tester) backBad1() {
	sf.t.Log("backBad1")
	sf.hand.IsBad_.Set()
	binMsg, err := sf.hand.FnBack([]byte("test_msg"))
	if err == nil {
		sf.t.Fatalf("backBad1(): err==nil")
	}
	if binMsg != nil {
		sf.t.Fatalf("backBad1(): binMsg!=nil")
	}
	sf.hand.IsBad_.Reset()
}

// Создание мок-обработчика запросов
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	sf.hand = NewMockHandlerServe("test_topic", "test_name")
	if sf.hand == nil {
		sf.t.Fatalf("newGood1(): handler==nil")
	}
	if name := sf.hand.Name(); name != "test_name" {
		sf.t.Fatalf("newGood1(): name(%v)!='test_name'", name)
	}
	if topic := sf.hand.Topic(); topic != "test_topic" {
		sf.t.Fatalf("newGood1(): topic(%v)!='test_topic'", topic)
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
	_ = NewMockHandlerServe("", "test_name")
}
