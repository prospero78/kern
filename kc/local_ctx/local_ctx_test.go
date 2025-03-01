package local_ctx

import (
	"context"
	"testing"
)

type tester struct {
	t   *testing.T
	ctx *LocalCtx
}

func TestLocalCtx(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.set()
	sf.get()
	sf.del()
	sf.done()
}

// Ожидает отмены контекста
func (sf *tester) done() {
	sf.t.Log("done")
	go sf.ctx.Cancel()
	sf.ctx.Done()
}

// Удаляет несуществующий ключ из локального контекста
func (sf *tester) del() {
	sf.t.Log("del")
	sf.ctx.Del("123")
	sf.ctx.Del("count")
	if _len := len(sf.ctx.dictVal); _len != 0 {
		sf.t.Fatalf("del(): len dict(%v)!=0", _len)
	}
}

// Возвращает хранимое значение
func (sf *tester) get() {
	sf.t.Log("get")
	val := sf.ctx.Get("count")
	count := val.Val().(int)
	if count != 15 {
		sf.t.Fatalf("get(): count(%v)!=15", count)
	}
}

// Устанавливает значение
func (sf *tester) set() {
	sf.t.Log("set")
	sf.ctx.Set("count", 5, "test_val")
	sf.ctx.Set("count", 15, "test_val1")
}

// Создание нового локального контекста
func (sf *tester) new() {
	sf.t.Log("new")
	sf.newBad1()
	sf.newGood1()
}

func (sf *tester) newGood1() {
	sf.t.Log("newGood1")
	ctx := context.Background()
	sf.ctx = NewLocalCtx(ctx).(*LocalCtx)
	_ = sf.ctx.Log()
	if lst := sf.ctx.SortedList(); lst == nil {
		sf.t.Fatalf("newGood1(): lst==nil")
	}
}

// Нет контекста ядра
func (sf *tester) newBad1() {
	sf.t.Log("newBad1")
	defer func() {
		if _panic := recover(); _panic == nil {
			sf.t.Fatalf("newBad1(): panic==nil")
		}
	}()
	var ctx context.Context
	sf.ctx = NewLocalCtx(ctx).(*LocalCtx)
}
