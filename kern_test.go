package kern

import (
	"os"
	"testing"

	"github.com/prospero78/kern/mock/mock_env"
)

type tester struct {
	t  *testing.T
	me *mock_env.MockEnv
}

func TestBuilders(t *testing.T) {
	sf := &tester{
		t:  t,
		me: mock_env.MakeEnv(),
	}
	_ = os.Unsetenv("LOCAL_STORE_PATH")
	_ = os.Setenv("LOCAL_STORE_PATH", "/store/store_builder")
	fnClear := func() {
		pwd := sf.me.Pwd() + "/store/store_builder"
		_ = os.RemoveAll(pwd)
	}
	fnClear()
	fnClear()
	sf.new()
}

// создание компонентов
func (sf *tester) new() {
	sf.t.Log("new")
	ctx := NewKernelCtx()
	if ctx == nil {
		sf.t.Fatalf("new(): IKernelCtx==nil")
	}
	store := NewKernelStore()
	err := store.Delete("test_builders")
	if err != nil {
		sf.t.Fatalf("new(): get empty key, store, err=%v", err)
	}

	safeBool := NewSafeBool()
	if safeBool == nil {
		sf.t.Fatalf("new(): ISafeBool==nil")
	}

	kernBus := NewKernelBusLocal()
	if kernBus == nil {
		sf.t.Fatalf("new(): (local) IKernelBus==nil")
	}

	kernServHttp := NewKernelServerHttp()
	go kernServHttp.Run()
	ctx.Cancel()
	ctx.Wg().Wait()
	ctx.Cancel()
	ctx.Wg().Wait()
}
