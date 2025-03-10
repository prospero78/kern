// package dict_sub_hook -- словарь потребителей топика по подписке
package dict_sub_hook

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
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
	sf.ctx.Set(string(handlerName), handler, "subscribe handler")
}

// Read -- вызывает все обработчики словаря подписок
func (sf *dictSubHook) Read(binMsg []byte) {
	sf.block.RLock()
	defer sf.block.RUnlock()
	for handlerName := range sf.dict {
		handler := sf.ctx.Get(string(handlerName)).Val().(IBusHandlerSubscribe)
		go handler.FnBack(binMsg)
	}
}
