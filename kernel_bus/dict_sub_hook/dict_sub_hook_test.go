package dict_sub_hook

import (
	"testing"

	"github.com/prospero78/kern/mock/mock_hand_sub"
)

type tester struct {
	t        *testing.T
	dict     *dictSubHook
	handSub  *mock_hand_sub.MockHandlerSub
	handSub2 *mock_hand_sub.MockHandlerSub
}

func TestDictSubWebHook(t *testing.T) {
	sf := &tester{
		t:        t,
		handSub:  mock_hand_sub.NewMockHandlerSub("hand_topic1", "hand_name1"),
		handSub2: mock_hand_sub.NewMockHandlerSub("hand_topic2", "hand_name2"),
	}
	sf.new()
	sf.add()
	sf.del()
	sf.read()
	sf.unsub()
}

// Отписка обработчика от топика
func (sf *tester) unsub() {
	sf.t.Log("unsub")
	sf.dict.Unsubscribe(sf.handSub)
}

// Чтение входящего сообщения по подписке
func (sf *tester) read() {
	sf.t.Log("read")
	sf.dict.Read([]byte("test_msg"))
}

// Удаляет хук из словаря
func (sf *tester) del() {
	sf.t.Log("del")
	sf.dict.Del("hand_name1")
	sf.dict.Del("hand_name1")
	if _len := len(sf.dict.dict); _len != 1 {
		sf.t.Fatalf("new(): len(%v)!=1", _len)
	}
}

// Добавляет хуки ыв словарь
func (sf *tester) add() {
	sf.t.Log("add")
	sf.addBad1()
	sf.addGood1()
}

func (sf *tester) addGood1() {
	sf.t.Log("addGood1")
	sf.dict.Subscribe(sf.handSub)
	sf.dict.Subscribe(sf.handSub)
	sf.dict.Subscribe(sf.handSub2)
	if _len := len(sf.dict.dict); _len != 2 {
		sf.t.Fatalf("new(): len(%v)!=2", _len)
	}
}

// нет веб-хука для добавления
func (sf *tester) addBad1() {
	sf.t.Log("addBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("addBad1(): panic==nil")
		}
	}()
	sf.dict.Subscribe(nil)
}

// Создание словаря
func (sf *tester) new() {
	sf.t.Log("new")
	sf.dict = NewDictSubHook().(*dictSubHook)
	if sf.dict == nil {
		sf.t.Fatalf("new(): dict==nil")
	}
	if _len := len(sf.dict.dict); _len != 0 {
		sf.t.Fatalf("new(): len(%v)!=0", _len)
	}
}
