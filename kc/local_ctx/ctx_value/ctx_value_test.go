package ctx_value

import (
	"testing"

	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
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
