package safe_bool

import (
	"testing"

	. "github.com/svi/kern/kernel_types"
)

type tester struct {
	t  *testing.T
	sb ISafeBool
}

func TestSAfeBool(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.create()
	sf.set()
	sf.reset()
}

// Сбрасывает хранимое значение
func (sf *tester) reset() {
	sf.t.Log("reset")
	sf.sb.Reset()
	if sf.sb.Get() {
		sf.t.Fatalf("reset(): SafeBool==true")
	}
	sf.sb.Reset()
	if sf.sb.Get() {
		sf.t.Fatalf("reset(): SafeBool==true")
	}
}

// Установка хранимого значения
func (sf *tester) set() {
	sf.t.Log("set")
	sf.sb.Set()
	if !sf.sb.Get() {
		sf.t.Fatalf("set(): SafeBool==true")
	}
	sf.sb.Set()
	if !sf.sb.Get() {
		sf.t.Fatalf("set(): SafeBool==true")
	}
}

// Создаёт потокобезопасный булевый признак
func (sf *tester) create() {
	sf.t.Log("create")
	sf.sb = NewSafeBool()
	if sf.sb == nil {
		sf.t.Fatalf("create(): SafeBool==nil")
	}
	if sf.sb.Get() {
		sf.t.Fatalf("create(): SafeBool==true")
	}
}
