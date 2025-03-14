package lst_sort

import (
	"testing"

	"gitp78su.ipnodns.ru/svi/kern/kc/local_ctx/ctx_value"
)

type tester struct {
	t   *testing.T
	lst *LstSort
}

func TestLstSort(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.new()
	sf.add()
	sf.del()
	sf.del2()
	sf.del3()
}

// Нет такого элемента
func (sf *tester) del3() {
	sf.t.Log("del3")
	val5 := ctx_value.NewCtxValue("val5", 5, "test 5")
	sf.lst.Del(val5)
	sf.lst.Del(nil)
}

// Удаляет элемент в середине
func (sf *tester) del2() {
	sf.t.Log("del2")
	val3 := sf.lst.Get(3)
	sf.lst.Del(val3)
	if _len := sf.lst.Size(); _len != 4 {
		sf.t.Fatalf("del2(): len(%v)!=4", _len)
	}
	val4 := sf.lst.Get(3)
	if val4.Key() != "val4" {
		sf.t.Fatalf("del2(): key(%v)!='val4'", val4.Key())
	}
	sf.lst.Add(val3)
	if _len := sf.lst.Size(); _len != 5 {
		sf.t.Fatalf("del2(): len(%v)!=5", _len)
	}
}

// Удаляет элемент из списка в конце
func (sf *tester) del() {
	sf.t.Log("del")

	val4 := sf.lst.Get(4)

	sf.lst.Del(val4)
	if _len := sf.lst.Size(); _len != 4 {
		sf.t.Fatalf("del(): len(%v)!=4", _len)
	}
	val3 := sf.lst.Get(3)
	if val3.Key() != "val3" {
		sf.t.Fatalf("del(): key(%v)!='val3'", val3.Key())
	}
	sf.lst.Add(val4)
	if _len := sf.lst.Size(); _len != 5 {
		sf.t.Fatalf("del(): len(%v)!=5", _len)
	}
	_ = sf.lst.List()
}

// Добавление элементов в список
func (sf *tester) add() {
	sf.t.Log("add")
	val0 := ctx_value.NewCtxValue("val0", 0, "test 0")
	val1 := ctx_value.NewCtxValue("val1", 1, "test 1")
	val2 := ctx_value.NewCtxValue("val2", 2, "test 2")
	val3 := ctx_value.NewCtxValue("val3", 3, "test 3")
	val4 := ctx_value.NewCtxValue("val4", 4, "test 4")
	sf.lst.Add(val1)
	sf.lst.Add(val4)
	sf.lst.Add(val0)
	sf.lst.Add(val2)
	sf.lst.Add(val3)
}

// Создание сортированного списка
func (sf *tester) new() {
	sf.t.Log("new")
	sf.lst = NewLstSort()
	if sf.lst == nil {
		sf.t.Fatalf("new(): lst==nil")
	}
}
