package safe_int

import (
	"testing"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

type tester struct {
	t  *testing.T
	si ISafeInt
}

func TestSafeInt(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.set()
	sf.reset()
}

// Сбрасывает хранимое значение
func (sf *tester) reset() {
	sf.t.Log("reset")
	sf.si.Reset()
	if sf.si.Get() != 0 {
		sf.t.Fatalf("reset(): ISafeInt!=0")
	}
	sf.si.Reset()
	if sf.si.Get() != 0 {
		sf.t.Fatalf("reset(): ISafeInt!=0")
	}
}

// Установка хранимого значения
func (sf *tester) set() {
	sf.t.Log("set")
	sf.si.Set(77)
	if sf.si.Get() != 77 {
		sf.t.Fatalf("set(): ISafeInt!=77")
	}
	sf.si.Set(-56)
	if sf.si.Get() != -56 {
		sf.t.Fatalf("set(): ISafeInt!=-56")
	}
}

// Создаёт потокобезопасный булевый признак
func (sf *tester) new() {
	sf.t.Log("new")
	sf.si = NewSafeInt()
	if sf.si == nil {
		sf.t.Fatalf("new(): ISafeInt==nil")
	}
	if sf.si.Get() != 0 {
		sf.t.Fatalf("new(): ISafeInt!=0")
	}
}
