package kernel_keeper

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/prospero78/kern/kernel_ctx/kernel_wg"
	. "github.com/prospero78/kern/kernel_types"
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
	wg := kernel_wg.GetKernelWg(ctx)
	defer fnCancel()
	sf := &tester{
		t:        t,
		ctx:      ctx,
		fnCancel: fnCancel,
		wg:       wg,
	}
	sf.get()
	sf.get2()
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
	close(keep.chSys_)
	time.Sleep(time.Millisecond * 10)
}
