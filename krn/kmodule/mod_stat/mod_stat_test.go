package mod_stat

import (
	"testing"
)

type tester struct {
	t *testing.T
}

func TestModStat(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
}

// Создание статистики модуля
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	stat := NewModStat("test")
	stat.Add(23)
	if svg := stat.SvgSec(); svg == "" {
		sf.t.Fatalf("newGood1(): svg is empty")
	}
}

// Нет имени статистики
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	_ = NewModStat("")
}
