package kernel_store

import (
	"os"
	"testing"
	"time"

	"github.com/prospero78/kern/kernel_ctx"
	. "github.com/prospero78/kern/kernel_types"
	"github.com/prospero78/kern/mock/mock_env"
)

type tester struct {
	t   *testing.T
	me  *mock_env.MockEnv
	ctx IKernelCtx
	wg  IKernelWg
}

func TestKernelStore(t *testing.T) {
	ctx := kernel_ctx.GetKernelCtx()
	sf := &tester{
		t:   t,
		me:  mock_env.MakeEnv(),
		ctx: ctx,
		wg:  ctx.Wg(),
	}
	_ = os.Unsetenv("LOCAL_STORE_PATH")
	_ = os.Setenv("LOCAL_STORE_PATH", "/store/store_kernel")
	fnClear := func() {
		pwd := sf.me.Pwd() + "/store/store_kernel"
		_ = os.RemoveAll(pwd)
	}
	fnClear()
	sf.new()
	sf.set()
	sf.get()
	sf.del()
	sf.close()
	sf.workBad1()
	time.Sleep(time.Second * 2)
}

// Установка ключа после закрытия хранилища
func (sf *tester) workBad1() {
	sf.t.Log("workBad1")
	err := kernStore.Set("test_key", []byte("test_val"))
	if err == nil {
		sf.t.Fatalf("workBad1(): set, err==nil")
	}
	_, err = kernStore.Get("test_key")
	if err == nil {
		sf.t.Fatalf("workBad1(): get, err==nil")
	}
	err = kernStore.Delete("test_key")
	if err == nil {
		sf.t.Fatalf("workBad1(): del, err==nil")
	}
}

// Удаляет значение
func (sf *tester) del() {
	sf.t.Log("del")
	err := kernStore.Delete("test_key")
	if err != nil {
		sf.t.Fatalf("del(): err=%v", err)
	}
	_, err = kernStore.Get("test_key")
	if err == nil {
		sf.t.Fatalf("del(): err==nil")
	}
}

// Получение значения по ключу
func (sf *tester) get() {
	sf.t.Log("get")
	binVal, err := kernStore.Get("test_key")
	if err != nil {
		sf.t.Fatalf("get(): err=%v", err)
	}
	strVal := string(binVal)
	if strVal != "test_val" {
		sf.t.Fatalf("get(): strVal(%v)!='test_val'", strVal)
	}
}

// Добавление ключа в хранилище
func (sf *tester) set() {
	sf.t.Log("set")
	err := kernStore.Set("test_key", []byte("test_val"))
	if err != nil {
		sf.t.Fatalf("set(): err=%v", err)
	}
}

// Закрытие хранилища
func (sf *tester) close() {
	sf.t.Log("close")
	err := sf.wg.Add("123")
	if err != nil {
		sf.t.Fatalf("close(): err=%v", err)
	}
	sf.ctx.Cancel()
	time.Sleep(time.Millisecond * 20)
	go sf.wg.Done("123")
	sf.wg.Wait()
	kernStore.close()
}

// Создаёт новое хранилище ядра
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Fatalf("newGood1(): panic=%v", _panic)
		}
	}()
	store := GetKernelStore()
	if store == nil {
		sf.t.Fatalf("newGood1(): KernelStore==nil")
	}
	store = GetKernelStore()
	if store == nil {
		sf.t.Fatalf("newGood1(): KernelStore==nil")
	}
}
