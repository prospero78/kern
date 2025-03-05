package mod_stat

import (
	"testing"
	"time"
)

type tester struct {
	t    *testing.T
	stat *ModStat
}

func TestModStat(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.event()
}

// Проверка генерации меток времени
func (sf *tester) event() {
	sf.t.Log("event")
	time.Sleep(time.Millisecond * 100)
}

// Создание статистики модуля
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	sf.stat = NewModStat("test")
	// Обязательно установить время сразу, для покрытия тестами
	sf.stat.timeMinute.Set(1)
	sf.stat.Add(23)
	if svg := sf.stat.SvgSec(); svg == "" {
		sf.t.Fatalf("newGood1(): svg is empty")
	}
	if svg := sf.stat.SvgMin(); svg == "" {
		sf.t.Fatalf("newGood1(): svg is empty")
	}
	if svg := sf.stat.SvgDay(); svg == "" {
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
