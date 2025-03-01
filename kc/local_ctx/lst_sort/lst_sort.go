// package lst_sort -- сортированный список значений контекста
package lst_sort

import (
	"sort"

	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/ktypes"
)

// LstSort -- сортированный список значений контекста
type LstSort struct {
	lstVal []ICtxValue // Сортированный список значений
}

// NewLstSort -- возвращает новый сортированный список значений контекста
func NewLstSort() *LstSort {
	sf := &LstSort{
		lstVal: []ICtxValue{},
	}
	return sf
}

// Add -- добавляет значение в список
func (sf *LstSort) Add(val ICtxValue) {
	sf.lstVal = append(sf.lstVal, val)
	sf.sort()
}

// Del -- удаляет элемент из списка
func (sf *LstSort) Del(val ICtxValue) {
	Hassert(val != nil, "LstSort.Del(): ICtxValue == nil")
	var (
		ind  int
		_val ICtxValue
	)
	for ind, _val = range sf.lstVal {
		if val == _val {
			break
		}
	}
	Hassert(_val != nil, "LstSort.Del(): key(%v) not found", val.Key())
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

// Get -- возвращает сортированный список
func (sf *LstSort) List() []ICtxValue {
	lst := []ICtxValue{}
	lst = append(lst, sf.lstVal...)
	return lst
}

// Сортирует элементы в списке
func (sf *LstSort) sort() {
	sort.Sort(sf)
}

// Swap -- меняет местами два элемента
func (sf *LstSort) Swap(ind1, ind2 int) {
	sf.lstVal[ind1], sf.lstVal[ind2] = sf.lstVal[ind2], sf.lstVal[ind1]
}

// Less -- сравнивает элементы по индексам
func (sf *LstSort) Less(ind1, ind2 int) bool {
	return sf.lstVal[ind1].Key() < sf.lstVal[ind2].Key()
}

// Len -- возвращает длину списка
func (sf *LstSort) Len() int {
	return len(sf.lstVal)
}
