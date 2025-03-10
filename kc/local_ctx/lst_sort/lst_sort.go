// package lst_sort -- сортированный список значений контекста
package lst_sort

import (
	"sort"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// LstSort -- сортированный список значений контекста
type LstSort struct {
	sync.RWMutex
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
	sf.Lock()
	defer sf.Unlock()
	Hassert(val != nil, "LstSort.Add(): ICtxValue==nil")
	sf.lstVal = append(sf.lstVal, val)
	sf.sort()
}

// Del -- удаляет элемент из списка
func (sf *LstSort) Del(val ICtxValue) {
	sf.Lock()
	defer sf.Unlock()
	if val == nil {
		return
	}
	sf.del(val)
}

// List -- возвращает сортированный список
func (sf *LstSort) List() []ICtxValue {
	sf.RLock()
	defer sf.RUnlock()
	return sf.list()
}

// Size -- возвращает длину списка
func (sf *LstSort) Size() int {
	sf.RLock()
	defer sf.RUnlock()
	return len(sf.lstVal)
}

// Get -- возвращает по индексу
func (sf *LstSort) Get(ind int) ICtxValue {
	sf.RLock()
	defer sf.RUnlock()
	Hassert(ind >= 0, "LstSort.Get(): ind(%v)<0", ind)
	Hassert(ind < len(sf.lstVal), "LstSort.Get(): ind(%v)>=len(%v)", ind, len(sf.lstVal))
	return sf.lstVal[ind]
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
