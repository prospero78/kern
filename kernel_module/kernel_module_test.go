package kernel_module

import (
	"testing"

	. "github.com/prospero78/kern/kernel_types"
)

type tester struct {
	t   *testing.T
	mod IKernelModule
}

func TestKernelModule(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.run()
	sf.isWork()
}

// Проверить признак работы
func (sf *tester) isWork() {
	sf.t.Log("isWork")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("isWork(): panic==nil")
		}
	}()
	_ = sf.mod.IsWork()
}

// Запускает модуль в работу
func (sf *tester) run() {
	sf.t.Log("run")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("run(): panic==nil")
		}
	}()
	sf.mod.Run()
}

// Создание нового модуля ядра
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	sf.mod = NewKernelModule("test_module")
	if name := sf.mod.Name(); name != "test_module" {
		sf.t.Fatalf("newGood1(): name(%v)!='test_module'", name)
	}
	if ctx := sf.mod.Ctx(); ctx == nil {
		sf.t.Fatalf("newGood1(): ctx==nil")
	}
}

// Нет имени модуля
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	_ = NewKernelModule("")
}
