// package dict_topic_sub -- потокобезопасный словарь подписчиков локальной шины
package dict_topic_sub

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/dict_sub_hook"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

type tReadReq struct {
	topic  ATopic
	binMsg []byte
}

// dictTopicSub -- потокобезопасный словарь подписчиков
type dictTopicSub struct {
	sync.RWMutex
	ctx           IKernelCtx
	dictTopicHook map[ATopic]IDictSubHook
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
	sf.RLock()
	defer sf.RUnlock()
	Hassert(topic != "", "dictTopicSub.Read(): topic is empty")
	msg := &tReadReq{
		topic:  topic,
		binMsg: binMsg,
	}
	dictHook := sf.dictTopicHook[msg.topic]
	if dictHook == nil {
		return
	}
	dictHook.Read(msg.binMsg)
}

// Subscribe -- подписывает обработчик на топик
func (sf *dictTopicSub) Subscribe(handler IBusHandlerSubscribe) {
	sf.Lock()
	defer sf.Unlock()
	Hassert(handler != nil, "dictTopicSub.Subscribe(): IBusHandlerSubscribe==nil")
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
	sf.Lock()
	defer sf.Unlock()
	if handler == nil {
		return
	}
	topic := handler.Topic()
	dictSubHook := sf.dictTopicHook[topic]
	if dictSubHook == nil {
		return
	}
	dictSubHook.Unsubscribe(handler)
}
