package safe_string

import (
	"testing"

	. "github.com/prospero78/kern/krn/ktypes"
)

type tester struct {
	t  *testing.T
	ss ISafeString
}

func TestSafeString(t *testing.T) {
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
	sf.ss.Reset()
	if sf.ss.Get() != "" {
		sf.t.Fatalf("reset(): ISafeString!=''")
	}
	sf.ss.Reset()
	if sf.ss.Get() != "" {
		sf.t.Fatalf("reset(): ISafeString!=''")
	}
	_ = sf.ss.Byte()
}

// Установка хранимого значения
func (sf *tester) set() {
	sf.t.Log("set")
	sf.ss.Set("77")
	if sf.ss.Get() != "77" {
		sf.t.Fatalf("set(): ISafeString!='77'")
	}
	sf.ss.Set("-56")
	if sf.ss.Get() != "-56" {
		sf.t.Fatalf("set(): ISafeString!='-56'")
	}
}

// Создаёт потокобезопасный булевый признак
func (sf *tester) new() {
	sf.t.Log("new")
	sf.ss = NewSafeString()
	if sf.ss == nil {
		sf.t.Fatalf("new(): ISafeString==nil")
	}
	if sf.ss.Get() != "" {
		sf.t.Fatalf("new(): ISafeString!=''")
	}
}
