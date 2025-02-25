package mock_env

import (
	"testing"
)

type tester struct {
	t  *testing.T
	me *MockEnv
}

func TestMockEnv(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.make()
	sf.reset()
}

// Сброс переменной окружения
func (sf *tester) reset() {
	sf.t.Log("reset")
	sf.resetBad1()
	pwd := sf.me.Pwd()
	if pwd == "" {
		sf.t.Fatalf("reset(): pwd i empty")
	}
}

// Нет переменой окружения
func (sf *tester) resetBad1() {
	sf.t.Log("resetBad1")
	defer func() {
		if panic_ := recover(); panic_ == nil {
			sf.t.Fatalf("resetBad1(): panic=%v", panic_)
		}
	}()
	sf.me.Reset("")
}

// Создание мок-окружения
func (sf *tester) make() {
	sf.t.Log("make")
	defer func() {
		if panic_ := recover(); panic_ != nil {
			sf.t.Fatalf("make(): panic=%v", panic_)
		}
	}()
	env := MakeEnv()
	if env == nil {
		sf.t.Fatalf("make(): env==nil")
	}
	sf.me = env
}
