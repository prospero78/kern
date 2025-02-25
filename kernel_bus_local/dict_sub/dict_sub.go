// package dict_sub -- потокобезопасный словарь подписчиков
package dict_sub

import (
	"sync"

	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	. "github.com/prospero78/kern/kernel_types"
)

// DictSub -- потокобезопасный словарь подписчиков
type DictSub struct {
	dictSub map[ATopic]IBusHandlerSubscribe
	block   sync.RWMutex
}

// NewDictSub -- возвращает потокобезопасный словарь подписчиков
func NewDictSub() *DictSub {
	sf := &DictSub{
		dictSub: map[ATopic]IBusHandlerSubscribe{},
	}
	return sf
}

// Call -- вызывает обработчики при поступлении события
func (sf *DictSub) Call(topic ATopic, binMsg []byte) {
	sf.block.RLock()
	defer sf.block.RUnlock()
	Hassert(topic != "", "DictSub.Call(): topic is empty")
	handler, isOk := sf.dictSub[topic]
	if !isOk {
		return
	}
	go handler.FnBack(binMsg)
}

// Add -- добавляет обработчик подписки в словарь
func (sf *DictSub) Add(handler IBusHandlerSubscribe) {
	Hassert(handler != nil, "DictSub.Add(): IBusHandlerSubscribe==nil")
	sf.block.Lock()
	defer sf.block.Unlock()
	_, isOk := sf.dictSub[handler.Topic()]
	if isOk {
		Hassert(false, "DictSub.Add(): handler of topic (%v) already exists", handler.Topic())
	}
	sf.dictSub[handler.Topic()] = handler
}

// Del -- удаляет подписчика из словаря
func (sf *DictSub) Del(handler IBusHandlerSubscribe) {
	Hassert(handler != nil, "DictSub.Del(): IBusHandlerSubscribe==nil")
	sf.block.Lock()
	defer sf.block.Unlock()
	_, isOk := sf.dictSub[handler.Topic()]
	if !isOk {
		return
	}
	delete(sf.dictSub, handler.Topic())
}
