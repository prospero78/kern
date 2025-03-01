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
	_ = os.Unsetenv("LOCAL_HTTP_URL")
	os.Setenv("LOCAL_HTTP_URL", "http://localhost:18311/")
	fnClear := func() {
		pwd := sf.me.Pwd() + "/store/store_builder"
		_ = os.RemoveAll(pwd)
	}
	fnClear()
	fnClear()
	sf.new()
	sf.newModBad()
}

func (sf *tester) newModBad() {
	sf.t.Log("newModBad")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newModBad(): panic==nil")
		}
	}()
	_ = NewKernelModule("")
}

// создание компонентов
func (sf *tester) new() {
	sf.t.Log("new")
	ctx := NewKernelCtx()
	if ctx == nil {
		sf.t.Fatalf("new(): IKernelCtx==nil")
	}
	store := NewKernelStoreKv()
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

	kernBusHttp := NewKernelBusHttp()
	if kernBusHttp == nil {
		sf.t.Fatalf("new(): (http) IKernelBus==nil")
	}

	monLocal := NewMonolitLocal("mon_local")
	if monLocal == nil {
		sf.t.Fatalf("new(): (local) IKernelMonolit==nil")
	}

	monHttp := NewMonolitHttp("mon_http")
	if monHttp == nil {
		sf.t.Fatalf("new(): (http) IKernelMonolit==nil")
	}

	mod := NewKernelModule("test_mod")
	if mod == nil {
		sf.t.Fatalf("new(): IKernelModule==nil")
	}

	clientLocal := NewClientBusLocal()
	if clientLocal == nil {
		sf.t.Fatalf("new(): (local) IBusClient==nil")
	}
	clientHttp := NewClientBusHttp("test_url")
	if clientHttp == nil {
		sf.t.Fatalf("new(): (http) IBusClient==nil")
	}

	modServHttp := NewModuleServHttp()
	if modServHttp == nil {
		sf.t.Fatalf("new(): modServHttp==nil")
	}

	modKernelCtx := NewModuleKernelCtx()
	if modKernelCtx == nil {
		sf.t.Fatalf("new(): modKernelCtx==nil")
	}

	kernServHttp := NewKernelServerHttp()
	go kernServHttp.Run()
	ctx.Cancel()
	ctx.Wg().Wait()
	ctx.Cancel()
	ctx.Wg().Wait()
}
