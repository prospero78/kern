// package lst_sort -- сортированный список значений контекста
package lst_sort

import (
	"sort"
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/ktypes"
)

// LstSort -- сортированный список значений контекста
type LstSort struct {
	lstVal []ICtxValue // Сортированный список значений
	block  sync.RWMutex
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
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.lstVal = append(sf.lstVal, val)
	sf.sort()
}

// Del -- удаляет элемент из списка
func (sf *LstSort) Del(val ICtxValue) {
	sf.block.Lock()
	defer sf.block.Unlock()
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
func (sf *LstSort) List() <-chan ICtxValue {
	sf.block.RLock()
	defer sf.block.RUnlock()
	chList := make(chan ICtxValue, len(sf.lstVal)+2)
	defer close(chList)
	for _, val := range sf.lstVal {
		chList <- val
	}
	return chList
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

// Len -- возвращает длину списка
func (sf *LstSort) Size() int {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return len(sf.lstVal)
}

// Get -- возвращает по индексу
func (sf *LstSort) Get(ind int) ICtxValue {
	sf.block.RLock()
	defer sf.block.RUnlock()
	Hassert(ind >= 0, "LstSort.Get(): ind(%v)<0", ind)
	if ind < len(sf.lstVal) {
		return sf.lstVal[ind]
	}
	return nil
}
