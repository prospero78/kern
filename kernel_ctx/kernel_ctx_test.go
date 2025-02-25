package kernel_ctx

import (
	"testing"
)

type tester struct {
	t *testing.T
}

func TestKernelCtx(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
}

// Создание контекста ядра
func (sf *tester) new() {
	sf.t.Log("new")
	ctx := GetKernelCtx()
	if ctx == nil {
		sf.t.Fatalf("new(): KernelCtx==nil")
	}
	if ctx := ctx.CtxBg(); ctx != kernCtx.ctxBg {
		sf.t.Fatalf("new(): ctx!=ctxBg")
	}
	if ctx := ctx.Ctx(); ctx != kernCtx.ctx {
		sf.t.Fatalf("new(): ctx!=kernel.ctx")
	}
	ctx.Add("counter", 5)
	counter := ctx.Get("counter").(int)
	if counter != 5 {
		sf.t.Fatalf("new(): counter(%v)!=5", counter)
	}
	ctx.Cancel()
	ctx.Done()
	ctx = GetKernelCtx()
	if ctx == nil {
		sf.t.Fatalf("new(): KernelCtx==nil")
	}
	if wg := ctx.Wg(); wg == nil {
		sf.t.Fatalf("new(): IKernelWg==nil")
	}
}
