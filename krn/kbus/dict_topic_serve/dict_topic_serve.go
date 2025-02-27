// package dict_topic_serve -- словарь топиков обработчиков запросов
package dict_topic_serve

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// dictServe -- потокобезопасный словарь обработчиков запросов
//
// Допускается только один обработчик запросов на один топик
type dictServe struct {
	ctx       IKernelCtx
	dictServe map[ATopic]IBusHandlerServe
	block     sync.RWMutex
}

// NewDictServe -- возвращает потокобезопасный словарь обработчиков запросов
func NewDictServe() IDictTopicServe {
	sf := &dictServe{
		ctx:       kctx.GetKernelCtx(),
		dictServe: make(map[ATopic]IBusHandlerServe, 0),
	}
	return sf
}

// Unregister -- удаляет обработчик запросов из словаря
func (sf *dictServe) Unregister(handler IBusHandlerServe) {
	Hassert(handler != nil, "dictServe.Unregister(): IBusHandlerSubscribe==nil")
	sf.block.Lock()
	defer sf.block.Unlock()
	delete(sf.dictServe, handler.Topic())
}

var TimeoutDefault = 15000

// SendRequest -- вызывает обработчик при поступлении запроса
func (sf *dictServe) SendRequest(topic ATopic, binReq []byte) ([]byte, error) {
	var handler IBusHandlerServe
	fnExtract := func() bool {
		sf.block.RLock()
		defer sf.block.RUnlock()
		var isOk bool
		handler, isOk = sf.dictServe[topic]
		return isOk
	}
	if !fnExtract() {
		return nil, fmt.Errorf("dictServe.SendRequest(): handler for topic (%v) not exists", topic)
	}
	var (
		chErr  = make(chan error, 2)
		binRes []byte
	)
	ctx, fnCancel := context.WithTimeout(sf.ctx.Ctx(), time.Millisecond*time.Duration(TimeoutDefault))
	defer fnCancel()
	fnCall := func() {
		defer close(chErr)
		var err error
		binRes, err = handler.FnBack(binReq)
		if err != nil {
			chErr <- err
		}
	}
	go fnCall()
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("dictServe.SendRequest(): in call for topic (%v), err=\n\t%w", topic, ctx.Err())
	case err := <-chErr:
		if err != nil {
			return nil, fmt.Errorf("dictServe.SendRequest(): error in call for topic (%v), err=\n\t%w", topic, err)
		}
	}
	return binRes, nil
}

// Register -- регистрирует обработчик запросов
func (sf *dictServe) Register(handler IBusHandlerServe) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(handler != nil, "dictServe.Register(): IBusHandlerSubscribe==nil")
	topic := handler.Topic()
	Hassert(topic != "", "dictServe.Register(): empty topic of handler")
	_, isRegister := sf.dictServe[topic]
	Hassert(!isRegister, "dictServe.Register(): handler of topic (%v) already register", handler.Topic())
	sf.dictServe[topic] = handler
}
