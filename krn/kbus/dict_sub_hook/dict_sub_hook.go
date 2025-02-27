// package dict_sub_hook -- словарь потребителей топика по подписке
package dict_sub_hook

import (
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// dictSubHook -- словарь потребителей топика по подписке
type dictSubHook struct {
	ctx   IKernelCtx
	dict  map[AHandlerName]bool // В качестве ключа -- URL веб-хука
	block sync.RWMutex
}

// NewDictSubHook -- возвращает новый словарь веб-хуков одного топика
func NewDictSubHook() IDictSubHook {
	sf := &dictSubHook{
		ctx:  kctx.GetKernelCtx(),
		dict: map[AHandlerName]bool{},
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
	sf.ctx.Del(string(handlerName))
}

// Subscribe -- добавляет в словарь подписки новый обработчик
func (sf *dictSubHook) Subscribe(handler IBusHandlerSubscribe) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(handler != nil, "dictSubHook.Subscribe(): handler==nil")
	handlerName := handler.Name()
	sf.dict[handlerName] = true
	sf.ctx.Set(string(handlerName), handler)
}

// Read -- вызывает все обработчики словаря подписок
func (sf *dictSubHook) Read(binMsg []byte) {
	sf.block.RLock()
	defer sf.block.RUnlock()
	for handlerName := range sf.dict {
		handler := sf.ctx.Get(string(handlerName)).(IBusHandlerSubscribe)
		go handler.FnBack(binMsg)
	}
}
