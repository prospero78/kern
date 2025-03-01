package kctx

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
	if ctx := ctx.BaseCtx(); ctx != kernCtx.ctx {
		sf.t.Fatalf("new(): ctx!=kernel.ctx")
	}
	ctx.Set("counter", 5, "test_counter")
	if ctx.Get("counter") == nil {
		sf.t.Fatalf("new(): counter==nil")
	}
	counter := ctx.Get("counter").Val().(int)
	if counter != 5 {
		sf.t.Fatalf("new(): counter(%v)!=5", counter)
	}
	ctx.Del("counter")
	ctx.Cancel()
	ctx.Done()
	ctx = GetKernelCtx()
	if ctx == nil {
		sf.t.Fatalf("new(): KernelCtx==nil")
	}
	if wg := ctx.Wg(); wg == nil {
		sf.t.Fatalf("new(): IKernelWg==nil")
	}
	if keep := ctx.Keeper(); keep == nil {
		sf.t.Fatalf("new(): IKernelKeeper==nil")
	}

}
