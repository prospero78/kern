// package dict_topic_sub -- потокобезопасный словарь подписчиков локальной шины
package dict_topic_sub

import (
	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kbus/dict_sub_hook"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

type tReadReq struct {
	topic  ATopic
	binMsg []byte
}

// dictTopicSub -- потокобезопасный словарь подписчиков
type dictTopicSub struct {
	ctx           IKernelCtx
	dictTopicHook map[ATopic]IDictSubHook

	chReadIn      chan *tReadReq
	chSubscribeIn chan IBusHandlerSubscribe
	chUnsubIn     chan IBusHandlerSubscribe
}

// NewDictTopicSub -- возвращает потокобезопасный словарь подписчиков
func NewDictTopicSub() IDictTopicSub {
	sf := &dictTopicSub{
		ctx:           kctx.GetKernelCtx(),
		chReadIn:      make(chan *tReadReq, 2),
		chSubscribeIn: make(chan IBusHandlerSubscribe, 2),
		chUnsubIn:     make(chan IBusHandlerSubscribe, 2),
		dictTopicHook: map[ATopic]IDictSubHook{},
	}
	go sf.run()
	return sf
}

// Read -- вызывает обработчики при поступлении события
func (sf *dictTopicSub) Read(topic ATopic, binMsg []byte) {
	Hassert(topic != "", "dictTopicSub.Read(): topic is empty")
	msg := &tReadReq{
		topic:  topic,
		binMsg: binMsg,
	}
	sf.chReadIn <- msg
}

// Subscribe -- подписывает обработчик на топик
func (sf *dictTopicSub) Subscribe(handler IBusHandlerSubscribe) {
	Hassert(handler != nil, "dictTopicSub.Subscribe(): IBusHandlerSubscribe==nil")
	topic := handler.Topic()
	Hassert(topic != "", "dictTopicSub.Subscribe(): topic is empty")
	sf.chSubscribeIn <- handler
}

// Unsubscribe -- отписывает обработчик
func (sf *dictTopicSub) Unsubscribe(handler IBusHandlerSubscribe) {
	if handler == nil {
		return
	}
	sf.chUnsubIn <- handler
}

func (sf *dictTopicSub) run() {
	for {
		select {
		case msg := <-sf.chReadIn:
			sf.read(msg)
		case handler := <-sf.chSubscribeIn:
			sf.subscribe(handler)
		case handler := <-sf.chUnsubIn:
			sf.unsub(handler)
		}
	}
}

// вызывает обработчики при поступлении события
func (sf *dictTopicSub) read(msg *tReadReq) {
	dictHook := sf.dictTopicHook[msg.topic]
	if dictHook == nil {
		return
	}
	dictHook.Read(msg.binMsg)
}

// подписывает обработчик на топик
func (sf *dictTopicSub) subscribe(handler IBusHandlerSubscribe) {
	topic := handler.Topic()
	dictSubHook := sf.dictTopicHook[topic]
	if dictSubHook == nil {
		dictSubHook = dict_sub_hook.NewDictSubHook()
		sf.dictTopicHook[topic] = dictSubHook
	}
	dictSubHook.Subscribe(handler)
}

// отписывает обработчик
func (sf *dictTopicSub) unsub(handler IBusHandlerSubscribe) {
	topic := handler.Topic()
	dictSubHook := sf.dictTopicHook[topic]
	if dictSubHook == nil {
		return
	}
	dictSubHook.Unsubscribe(handler)
}
