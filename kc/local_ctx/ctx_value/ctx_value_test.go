package ctx_value

import (
	"testing"

	. "github.com/prospero78/kern/krn/kalias"
	. "github.com/prospero78/kern/krn/ktypes"
)

type tester struct {
	t      *testing.T
	val    ICtxValue
	create ATime
}

func TestCtxValue(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.upd()
}

// Обновление значения
func (sf *tester) upd() {
	sf.t.Log("upd")
	sf.val.Update(35, "new_value")
	if val := sf.val.Val().(int); val != 35 {
		sf.t.Fatalf("upd(): val(%v)!=35", val)
	}
	create := sf.val.CreateAt()
	if create != sf.create {
		sf.t.Fatalf("upd(): create bad")
	}
	if comment := sf.val.Comment(); comment != "new_value" {
		sf.t.Fatalf("upd(): comment(%v)!='new_value'", comment)
	}
	if update := sf.val.UpdateAt(); update == "" {
		sf.t.Fatalf("upd(): update empty")
	}
}

// Создаёт значение локального контекста
func (sf *tester) new() {
	sf.t.Log("new")
	sf.val = NewCtxValue("test_key", 5, "test_value")
	if sf.val == nil {
		sf.t.Fatalf("new(): val==nil")
	}
	if val := sf.val.Val().(int); val != 5 {
		sf.t.Fatalf("new(): val(%v)!=5", val)
	}
	create := sf.val.CreateAt()
	if create == "" {
		sf.t.Fatalf("new(): create is empty")
	}
	sf.create = create
	if comment := sf.val.Comment(); comment != "test_value" {
		sf.t.Fatalf("new(): comment(%v)!='test_value'", comment)
	}
	if key := sf.val.Key(); key != "test_key" {
		sf.t.Fatalf("new(): key(%v)!='test_key'", key)
	}
	if update := sf.val.UpdateAt(); update != "" {
		sf.t.Fatalf("new(): update not empty")
	}
}
