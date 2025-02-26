// package dict_sub_hook -- словарь потребителей топика по подписке
package dict_sub_hook

import (
	"sync"

	. "github.com/prospero78/kern/helpers"
	"github.com/prospero78/kern/kernel_ctx"
	. "github.com/prospero78/kern/kernel_types"
)

// dictSubHook -- словарь потребителей топика по подписке
type dictSubHook struct {
	ctx   IKernelCtx
	dict  map[string]bool // В качестве ключа -- URL веб-хука
	block sync.RWMutex
}

// NewDictSubHook -- возвращает новый словарь веб-хуков одного топика
func NewDictSubHook() IDictSubHook {
	sf := &dictSubHook{
		ctx:  kernel_ctx.GetKernelCtx(),
		dict: map[string]bool{},
	}
	return sf
}

// Unsubscribe -- удаляет из словаря подписки обработчик
func (sf *dictSubHook) Unsubscribe(handler IBusHandlerSubscribe) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(handler != nil, "dictSubHook.Unsubscribe(): handler==nil")
	handlerName := handler.Name()
	delete(sf.dict, handlerName)
	sf.ctx.Del(handlerName)
}

// Subscribe -- добавляет в словарь подписки новый обработчик
func (sf *dictSubHook) Subscribe(handler IBusHandlerSubscribe) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(handler != nil, "dictSubHook.Subscribe(): handler==nil")
	handlerName := handler.Name()
	sf.dict[handlerName] = true
	sf.ctx.Add(handlerName, handler)
}

// Call -- вызывает все обработчики словаря подписок
func (sf *dictSubHook) Call(binMsg []byte) {
	sf.block.RLock()
	defer sf.block.RUnlock()
	for key := range sf.dict {
		handler := sf.ctx.Get(key).(IBusHandlerSubscribe)
		go handler.FnBack(binMsg)
	}
}

// Del -- удаляет из словаря веб-хук
func (sf *dictSubHook) Del(webHook string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	delete(sf.dict, webHook)
}
