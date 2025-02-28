package kmonolit

import (
	"testing"

	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/krn/kmodule"
	. "github.com/prospero78/kern/krn/ktypes"
)

type tester struct {
	t   *testing.T
	mon IKernelMonolit
}

func TestKernMono(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.run()
	sf.add()
	sf.done()
}

func (sf *tester) done() {
	sf.t.Log("done")
	ctx := kctx.GetKernelCtx()
	ctx.Cancel()
	ctx.Wg().Wait()
	sf.mon.(*kMonolit).close()
	sf.mon.Run()
}

// Добавление модуля в монолит
func (sf *tester) add() {
	sf.t.Log("add")
	sf.addGood1()
}

type mod struct {
	IKernelModule
}

func newMod(name AModuleName) IKernelModule {
	sf := &mod{
		IKernelModule: kmodule.NewKernelModule(name),
	}
	return sf
}

func (sf *mod) Run() {}

func (sf *tester) addGood1() {
	sf.t.Log("addGood1")
	mod := newMod("test_module")
	sf.mon.Add(mod)
}

func (sf *tester) run() {
	sf.t.Log("run")
	mod := newMod("test_mod1")
	sf.mon.Add(mod)
	sf.mon.Run()
	isWork := sf.mon.IsWork()
	if !isWork {
		sf.t.Fatalf("newGood1(): isWork==false")
	}
}

// Создаёт новый монолит
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("newGood1(): panic=%v", _panic)
		}
	}()
	ctx := kctx.GetKernelCtx()
	ctx.Set("isLocal", true, "type bus")
	sf.mon = GetMonolit("test_monolit")
	isLocal := sf.mon.IsLocal()
	if !isLocal {
		sf.t.Fatalf("newGood1(): isLocal==false")
	}
	if name := sf.mon.Name(); name != "test_monolit" {
		sf.t.Fatalf("newGood1(): name(%v)!='test_monolit'", name)
	}
	if log := sf.mon.Log(); log == nil {
		sf.t.Fatalf("newGood1(): log==nil")
	}
	sf.mon = GetMonolit("")
}

// Нет признака локальности
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	_ = GetMonolit("test_32")
}
