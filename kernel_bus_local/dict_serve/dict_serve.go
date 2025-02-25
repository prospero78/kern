// package dict_serve -- словарь обработчиков входящих запросов
package dict_serve

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	. "github.com/prospero78/kern/kernel_types"
)

// DictServe -- потокобезопасный словарь обработчиков входящих запросов
type DictServe struct {
	ctx       IKernelCtx
	dictServe map[ATopic]IBusHandlerServe
	block     sync.RWMutex
}

// NewDictServe -- возвращает потокобезопасный словарь подписчиков
func NewDictServe(ctx IKernelCtx) *DictServe {
	Hassert(ctx != nil, "NewDictServe(): IKernelCtx==nil")
	sf := &DictServe{
		ctx:       ctx,
		dictServe: make(map[ATopic]IBusHandlerServe, 0),
	}
	return sf
}

var TimeoutDefault = 15000

// Call -- вызывает обработчики при поступлении события
func (sf *DictServe) Call(ctx context.Context, topic ATopic, binMsg []byte) ([]byte, error) {
	Hassert(ctx != nil, "DictServe.Call(): ctx==nil")
	var handler IBusHandlerServe
	fnExtract := func() bool {
		sf.block.RLock()
		defer sf.block.RUnlock()
		var isOk bool
		handler, isOk = sf.dictServe[topic]
		return isOk
	}
	if !fnExtract() {
		return nil, fmt.Errorf("DictServe.Call(): handler for topic (%v) not exists", topic)
	}
	var (
		chErr  = make(chan error, 2)
		binRes []byte
	)
	ctx, fnCancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(TimeoutDefault))
	defer fnCancel()
	fnCall := func() {
		defer close(chErr)
		var err error
		binRes, err = handler.FnBack(binMsg)
		if err != nil {
			chErr <- err
		}
	}
	go fnCall()
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("DictServe.Call(): in call for topic (%v), err=\n\t%w", topic, ctx.Err())
	case err := <-chErr:
		if err != nil {
			return nil, fmt.Errorf("DictServe.Call(): error in call for topic (%v), err=\n\t%w", topic, err)
		}
	}
	return binRes, nil
}

// Add -- добавляет обработчик подписки в словарь
func (sf *DictServe) Add(handler IBusHandlerServe) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(handler != nil, "DictServe.Add(): IBusHandlerSubscribe==nil")
	topic := handler.Topic()
	Hassert(topic != "", "DictServe.Add(): empty topic of handler")
	_, isOk := sf.dictServe[topic]
	if isOk {
		Hassert(false, "DictServe.Add(): handler of topic (%v) already exists", handler.Topic())
	}
	sf.dictServe[topic] = handler
}

// Del -- удаляет подписчика из словаря
func (sf *DictServe) Del(handler IBusHandlerServe) {
	Hassert(handler != nil, "DictServe.Del(): IBusHandlerSubscribe==nil")
	sf.block.Lock()
	defer sf.block.Unlock()
	_, isOk := sf.dictServe[handler.Topic()]
	if !isOk {
		return
	}
	delete(sf.dictServe, handler.Topic())
}
