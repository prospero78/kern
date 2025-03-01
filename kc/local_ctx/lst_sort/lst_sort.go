// package lst_sort -- сортированный список значений контекста
package lst_sort

import (
	"sort"

	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/ktypes"
)

// LstSort -- сортированный список значений контекста
type LstSort struct {
	chAddIn  chan ICtxValue
	chAddOut chan int

	chDelIn  chan ICtxValue
	chDelOut chan int

	chListIn  chan int
	chListOut chan []ICtxValue

	chSizeIn  chan int
	chSizeOut chan int

	chGetIn  chan int
	chGetOut chan ICtxValue
	lstVal   []ICtxValue // Сортированный список значений
}

// NewLstSort -- возвращает новый сортированный список значений контекста
func NewLstSort() *LstSort {
	sf := &LstSort{
		chAddIn:  make(chan ICtxValue, 2),
		chAddOut: make(chan int, 2),

		chDelIn:  make(chan ICtxValue, 2),
		chDelOut: make(chan int, 2),

		chListIn:  make(chan int, 2),
		chListOut: make(chan []ICtxValue, 2),

		chSizeIn:  make(chan int, 2),
		chSizeOut: make(chan int, 2),

		chGetIn:  make(chan int, 2),
		chGetOut: make(chan ICtxValue, 2),
		lstVal:   []ICtxValue{},
	}
	go sf.run()
	return sf
}

// Add -- добавляет значение в список
func (sf *LstSort) Add(val ICtxValue) {
	Hassert(val != nil, "LstSort.Add(): ICtxValue==nil")
	sf.chAddIn <- val
	<-sf.chAddOut
}

// Del -- удаляет элемент из списка
func (sf *LstSort) Del(val ICtxValue) {
	if val == nil {
		return
	}
	sf.chDelIn <- val
	<-sf.chDelOut
}

// List -- возвращает сортированный список
func (sf *LstSort) List() []ICtxValue {
	sf.chListIn <- 1
	return <-sf.chListOut
}

// Size -- возвращает длину списка
func (sf *LstSort) Size() int {
	sf.chSizeIn <- 1
	return <-sf.chSizeOut
}

// Get -- возвращает по индексу
func (sf *LstSort) Get(ind int) ICtxValue {
	Hassert(ind >= 0, "LstSort.Get(): ind(%v)<0", ind)
	sf.chGetIn <- ind
	return <-sf.chGetOut
}

func (sf *LstSort) run() {
	for {
		select {
		case val := <-sf.chAddIn:
			sf.lstVal = append(sf.lstVal, val)
			sf.sort()
			sf.chAddOut <- 1
		case val := <-sf.chDelIn:
			sf.del(val)
			sf.chDelOut <- 1
		case <-sf.chListIn:
			sf.chListOut <- sf.list()
		case <-sf.chSizeIn:
			sf.chSizeOut <- len(sf.lstVal)
		case ind := <-sf.chGetIn:
			Hassert(ind < len(sf.lstVal), "LstSort.run(): ind(%v)<len(%v)", ind, len(sf.lstVal))
			sf.chGetOut <- sf.lstVal[ind]
		}
	}
}

// удаляет элемент из списка
func (sf *LstSort) del(val ICtxValue) {
	var (
		ind  int
		_val ICtxValue
	)
	for ind, _val = range sf.lstVal {
		if val == _val {
			break
		}
		_val = nil
	}
	if _val == nil {
		return
	}
	lst0 := sf.lstVal[:ind]
	lst1 := []ICtxValue{}
	if ind < len(sf.lstVal)-1 {
		lst1 = sf.lstVal[ind+1:]
	}
	sf.lstVal = sf.lstVal[:0]
	sf.lstVal = append(sf.lstVal, lst0...)
	sf.lstVal = append(sf.lstVal, lst1...)
	sf.sort()
}

// возвращает сортированный список
func (sf *LstSort) list() []ICtxValue {
	lst := make([]ICtxValue, 0, len(sf.lstVal))
	lst = append(lst, sf.lstVal...)
	return lst
}

// Сортирует элементы в списке
func (sf *LstSort) sort() {
	sort.Sort(sf)
}

// Swap -- НЕ ИСПОЛЬЗОВАТЬ меняет местами два элемента
func (sf *LstSort) Swap(ind1, ind2 int) {
	sf.lstVal[ind1], sf.lstVal[ind2] = sf.lstVal[ind2], sf.lstVal[ind1]
}

// Less -- НЕ ИСПОЛЬЗОВАТЬ сравнивает элементы по индексам
func (sf *LstSort) Less(ind1, ind2 int) bool {
	return sf.lstVal[ind1].Key() < sf.lstVal[ind2].Key()
}

// Len -- НЕ ИСПОЛЬЗОВАТЬ возвращает длину списка
func (sf *LstSort) Len() int {
	return len(sf.lstVal)
}
