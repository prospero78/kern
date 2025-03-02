package safe_int

import (
	"testing"

	. "github.com/prospero78/kern/krn/ktypes"
)

type tester struct {
	t  *testing.T
	sb ISafeInt
}

func TestSAfeBool(t *testing.T) {
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
	sf.sb.Reset()
	if sf.sb.Get() != 0 {
		sf.t.Fatalf("reset(): ISafeInt!=0")
	}
	sf.sb.Reset()
	if sf.sb.Get() != 0 {
		sf.t.Fatalf("reset(): ISafeInt!=0")
	}
}

// Установка хранимого значения
func (sf *tester) set() {
	sf.t.Log("set")
	sf.sb.Set(77)
	if sf.sb.Get() != 77 {
		sf.t.Fatalf("set(): ISafeInt!=77")
	}
	sf.sb.Set(-56)
	if sf.sb.Get() != -56 {
		sf.t.Fatalf("set(): ISafeInt!=-56")
	}
}

// Создаёт потокобезопасный булевый признак
func (sf *tester) new() {
	sf.t.Log("new")
	sf.sb = NewSafeInt()
	if sf.sb == nil {
		sf.t.Fatalf("new(): ISafeInt==nil")
	}
	if sf.sb.Get() != 0 {
		sf.t.Fatalf("new(): ISafeInt!=0")
	}
}
