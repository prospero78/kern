// package dict_topic_sub -- потокобезопасный словарь подписчиков локальной шины
package dict_topic_sub

import (
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kbus/dict_sub_hook"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// dictTopicSub -- потокобезопасный словарь подписчиков
type dictTopicSub struct {
	ctx           IKernelCtx
	dictTopicHook map[ATopic]IDictSubHook
	block         sync.RWMutex
}

// NewDictTopicSub -- возвращает потокобезопасный словарь подписчиков
func NewDictTopicSub() IDictTopicSub {
	sf := &dictTopicSub{
		ctx:           kctx.GetKernelCtx(),
		dictTopicHook: map[ATopic]IDictSubHook{},
	}
	return sf
}

// Read -- вызывает обработчики при поступлении события
func (sf *dictTopicSub) Read(topic ATopic, binMsg []byte) {
	sf.block.RLock()
	defer sf.block.RUnlock()
	Hassert(topic != "", "dictTopicSub.Read(): topic is empty")
	dictHook := sf.dictTopicHook[topic]
	if dictHook == nil {
		return
	}
	dictHook.Read(binMsg)
}

// Subscribe -- подписывает обработчик на топик
func (sf *dictTopicSub) Subscribe(handler IBusHandlerSubscribe) {
	Hassert(handler != nil, "dictTopicSub.Subscribe(): IBusHandlerSubscribe==nil")
	sf.block.Lock()
	defer sf.block.Unlock()
	topic := handler.Topic()
	Hassert(topic != "", "dictTopicSub.Subscribe(): topic is empty")
	dictSubHook := sf.dictTopicHook[topic]
	if dictSubHook == nil {
		dictSubHook = dict_sub_hook.NewDictSubHook()
		sf.dictTopicHook[topic] = dictSubHook
	}
	dictSubHook.Subscribe(handler)
}

// Unsubscribe -- отписывает обработчик
func (sf *dictTopicSub) Unsubscribe(handler IBusHandlerSubscribe) {
	Hassert(handler != nil, "dictTopicSub.Unsubscribe(): IBusHandlerSubscribe==nil")
	sf.block.Lock()
	defer sf.block.Unlock()
	topic := handler.Topic()
	dictHook := sf.dictTopicHook[topic]
	dictHook.Unsubscribe(handler)
}
