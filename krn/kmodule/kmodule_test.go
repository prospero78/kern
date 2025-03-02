package kmodule

import (
	"fmt"
	"testing"
	"time"

	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
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
	sf.done()
	sf.recErr()
}

// Регистрация ошибки
func (sf *tester) recErr() {
	sf.t.Log("recErr")
	mod := sf.mod.(*kModule)
	err := fmt.Errorf("tra-la-la")
	mod.recErr(err)
}

// Работа после остановки локальной шины
func (sf *tester) done() {
	sf.t.Log("done")
	kCtx := kctx.GetKernelCtx()

	time.Sleep(time.Millisecond * 250)
	kCtx.Cancel()
	kCtx.Done()
	time.Sleep(time.Millisecond * 200)
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
	mod := sf.mod.(*kModule)
	mod.timePhase.Set(5) // Настройка переменной модуля
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
	if _log := sf.mod.Log(); _log == nil {
		sf.t.Fatalf("newGood1(): log==nil")
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
