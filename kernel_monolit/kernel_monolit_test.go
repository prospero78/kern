package kernel_monolit

import (
	"testing"

	"github.com/prospero78/kern/kernel_ctx"
	. "github.com/prospero78/kern/kernel_types"
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
	sf.add()
}

// Добавление модуля в монолит
func (sf *tester) add() {
	sf.t.Log("add")
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
	ctx := kernel_ctx.GetKernelCtx()
	ctx.Set("isLocal", true)
	sf.mon = NewMonolit()
}

// Нет признака локальности
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	_ = NewMonolit()
}
