package kernel_keeper

import (
	"context"
	"os"
	"testing"
	"time"

	"gitp78su.ipnodns.ru/svi/kern/krn/kctx/kwg"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

type tester struct {
	t        *testing.T
	ctx      context.Context
	fnCancel func()
	wg       IKernelWg
}

func TestKernelKeeper(t *testing.T) {
	ctxBg := context.Background()
	ctx, fnCancel := context.WithCancel(ctxBg)
	wg := kwg.GetKernelWg(ctx)
	defer fnCancel()
	sf := &tester{
		t:        t,
		ctx:      ctx,
		fnCancel: fnCancel,
		wg:       wg,
	}
	sf.get()
	sf.get2()
	sf.done()
}

// Отмена контекста приложения
func (sf *tester) done() {
	sf.t.Log("done")
	sf.fnCancel()
	time.Sleep(time.Millisecond * 10)
	chSys := make(chan os.Signal, 1)
	kernKeep.run(chSys)
}

type sysSig struct {
}

func (sf *sysSig) String() string {
	return "test_sig"
}

func (sf *sysSig) Signal() {
}
func (sf *tester) get2() {
	sf.t.Log("get2")
	chSys := make(chan os.Signal, 2)
	sig := &sysSig{}
	chSys <- sig
	go kernKeep.run(chSys)
	sf.fnCancel()
	sf.wg.Wait()
}

// Получает сторож ядра
func (sf *tester) get() {
	sf.t.Log("get")
	keep := GetKernelKeeper(sf.ctx, sf.fnCancel, sf.wg)
	if keep == nil {
		sf.t.Fatalf("get(): IKernelKeeper==nil")
	}
	_ = GetKernelKeeper(sf.ctx, sf.fnCancel, sf.wg)
	if log := keep.Log(); log == nil {
		sf.t.Fatalf("get(): log==nil")
	}
	close(keep.chSys_)
	time.Sleep(time.Millisecond * 10)
}
