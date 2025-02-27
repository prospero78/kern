package log_msg

import (
	"strings"
	"testing"
)

type tester struct {
	t *testing.T
}

func TestLogMsg(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
}

// Создаёт сообщение
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	_ = NewLogMsg(-3, "test msg -3")
	_ = NewLogMsg(-2, "test msg -2")
	_ = NewLogMsg(-1, "test msg -1")
	msg := NewLogMsg(0, "test msg 0")
	if msg := msg.Msg(); msg != "test msg 0" {
		sf.t.Fatalf("newGood1(): msg(%v)!='test msg 0'", msg)
	}
	if lvl := msg.Level(); lvl != "ERRO" {
		sf.t.Fatalf("newGood1(): lvl(%v)!='ERRO'", lvl)
	}
	if create := msg.CreateAt(); create == "" {
		sf.t.Fatalf("newGood1(): create is empty")
	}
	if str := msg.String(); !strings.Contains(str,"ERRO   2") {
		sf.t.Fatalf("newGood1(): str(%v)!=`ERRO   2`", str)
	}
}

// Неправильный уровень сообщения
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	_ = NewLogMsg(-10, "test msg")
}
