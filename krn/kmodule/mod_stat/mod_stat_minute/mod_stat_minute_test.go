package mod_stat_minute

import (
	"testing"
)

type tester struct {
	t    *testing.T
	stat *ModStatMinutes
}

func TestModStatMinute(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.add()
}

// Добавляет событие в стату
func (sf *tester) add() {
	sf.t.Log("add")
	sf.stat.Add(12)
	sf.stat.momentAt = 0
	sf.stat.Add(11)
	sf.stat.momentAt = 0
	sf.stat.Add(3)
	if svg := sf.stat.Svg(); svg == "" {
		sf.t.Fatal("add(): svg is empty")
	}
}

// Создаёт новую секундную статистику модуля
func (sf *tester) new() {
	sf.t.Log("new")
	sf.stat = NewModStatMinute()
	if sf.stat == nil {
		sf.t.Fatalf("new(): stat==nil")
	}
}
