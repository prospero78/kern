package kernel_wg

import (
	"context"
	"testing"
	"time"

	. "github.com/prospero78/kern/kernel_types"
)

type tester struct {
	t        *testing.T
	ctx      context.Context
	fnCancel func()
	wg       IKernelWg
}

func TestKernelWG(t *testing.T) {
	ctxBg := context.Background()
	ctx, fnCancel := context.WithCancel(ctxBg)
	defer fnCancel()
	sf := &tester{
		t:        t,
		ctx:      ctx,
		fnCancel: fnCancel,
	}
	sf.new()
	sf.add()
	sf.done()
	sf.wait()
	sf.addBad3()
}

// Попытка добавления после закрытия ожидателя
func (sf *tester) addBad3() {
	sf.t.Log("addBad3")
	err := sf.wg.Add("test_stream")
	if err == nil {
		sf.t.Fatalf("addBad3(): err==nil")
	}
	_len := sf.wg.Len()
	if _len != 0 {
		sf.t.Fatalf("addBad3(): len(%v)!=0", _len)
	}
}

// Убирает имя потока из ожидателя
func (sf *tester) done() {
	sf.t.Log("done")
	sf.wg.Done("test_stream")
	lst := sf.wg.List()
	if len(lst) != 0 {
		sf.t.Fatalf("addBad1(): len(lst)!=0, lst=%#v", lst)
	}
}

// Добавление потока ожидания
func (sf *tester) add() {
	sf.t.Log("add")
	sf.addGood1()
	sf.addBad1()
	sf.addBad2()
}

// Уже есть такое имя потока
func (sf *tester) addBad2() {
	sf.t.Log("addBad2")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("addBad1(): panic==nil")
		}
		lst := sf.wg.List()
		if len(lst) != 1 {
			sf.t.Fatalf("addBad1(): len(lst)!=1, lst=%#v", lst)
		}
	}()
	_ = sf.wg.Add("test_stream")
}

// Пустое имя потока
func (sf *tester) addBad1() {
	sf.t.Log("addBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("addBad1(): panic==nil")
		}
	}()
	_ = sf.wg.Add("")
}

func (sf *tester) addGood1() {
	sf.t.Log("addGood1")
	err := sf.wg.Add("test_stream")
	if err != nil {
		sf.t.Fatalf("addGood1(): err=%v", err)
	}
}

// Ожидание завершения ожидателя потоков
func (sf *tester) wait() {
	sf.t.Log("wait")
	go sf.wg.Wait()
	time.Sleep(time.Millisecond * 10)
	sf.fnCancel()
	time.Sleep(time.Millisecond * 10)
	if sf.wg.IsWork() {
		sf.t.Fatalf("wait(): isWork==true")
	}
}

// Создаёт ожидатель потоков ядра
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
	sf.wg = GetKernelWg(sf.ctx)
	if sf.wg == nil {
		sf.t.Fatalf("newGood1(): KernelWg==nil")
	}
	if !sf.wg.IsWork() {
		sf.t.Fatalf("wait(): isWork==false")
	}
	wg := GetKernelWg(sf.ctx)
	if sf.wg != wg {
		sf.t.Fatalf("newGood1(): bad IKernelWg")
	}
}

// Нет контекста ядра
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic=nil")
		}
	}()
	// defer sf.panicUse("newBad1(): ")
	var ctx context.Context
	_ = GetKernelWg(ctx)
}
