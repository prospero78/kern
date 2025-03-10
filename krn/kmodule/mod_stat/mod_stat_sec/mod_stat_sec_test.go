package mod_stat_sec

import (
	"testing"
)

type tester struct {
	t    *testing.T
	stat *ModStatSec
}

func TestModStatSec(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.add()
	sf.sum()
}

// Возвращает сумму значений за минуту
func (sf *tester) sum() {
	sf.t.Log("sum")
	if sum := sf.stat.Sum(); sum == 0 {
		sf.t.Fatalf("sum(): sum==0")
	}
}

// Добавляет событие в стату
func (sf *tester) add() {
	sf.t.Log("add")
	sf.stat.Add(12)
	sf.stat.momentAt = 0
	sf.stat.Add(11)
	sf.stat.momentAt = 0
	sf.stat.Add(3)
	sf.stat.Add(3)
	if svg := sf.stat.Svg(); svg == "" {
		sf.t.Fatal("add(): svg is empty")
	}
}

// Создаёт новую секундную статистику модуля
func (sf *tester) new() {
	sf.t.Log("new")
	sf.stat = NewModStatSec()
	if sf.stat == nil {
		sf.t.Fatalf("new(): stat==nil")
	}
	if svg := sf.stat.Svg(); svg == "" {
		sf.t.Fatal("new(): svg is empty")
	}
}
