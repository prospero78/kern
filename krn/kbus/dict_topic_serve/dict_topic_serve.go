// package dict_topic_serve -- словарь топиков обработчиков запросов
package dict_topic_serve

import (
	"context"
	"fmt"
	"time"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_msg/msg_serve"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// dictServe -- потокобезопасный словарь обработчиков запросов
//
// Допускается только один обработчик запросов на один топик
type dictServe struct {
	ctx IKernelCtx

	chUnregisterIn chan IBusHandlerServe
	dictServe      map[ATopic]IBusHandlerServe

	chSendRequestIn  chan *msg_serve.ServeReq
	chSendRequestOut chan *serveResp

	chRegisterIn  chan IBusHandlerServe
	chRegisterOut chan bool
}

// NewDictServe -- возвращает потокобезопасный словарь обработчиков запросов
func NewDictServe() IDictTopicServe {
	sf := &dictServe{
		ctx: kctx.GetKernelCtx(),

		chUnregisterIn: make(chan IBusHandlerServe, 5),
		dictServe:      map[ATopic]IBusHandlerServe{},

		chSendRequestIn:  make(chan *msg_serve.ServeReq, 5),
		chSendRequestOut: make(chan *serveResp, 5),

		chRegisterIn:  make(chan IBusHandlerServe, 5),
		chRegisterOut: make(chan bool, 5),
	}
	go sf.run()
	return sf
}

// Register -- регистрирует обработчик запросов
func (sf *dictServe) Register(handler IBusHandlerServe) {
	Hassert(handler != nil, "dictServe.Register(): IBusHandlerSubscribe==nil")
	topic := handler.Topic()
	Hassert(topic != "", "dictServe.Register(): empty topic of handler")
	sf.chRegisterIn <- handler
	isTwinRegister := <-sf.chRegisterOut
	Hassert(!isTwinRegister, "dictServe.Register(): handler of topic (%v) already register", handler.Topic())
}

// Unregister -- удаляет обработчик запросов из словаря
func (sf *dictServe) Unregister(handler IBusHandlerServe) {
	Hassert(handler != nil, "dictServe.Unregister(): IBusHandlerSubscribe==nil")
	sf.chUnregisterIn <- handler
}

type serveResp struct {
	binResp []byte
	err     error
}

// SendRequest -- вызывает обработчик при поступлении запроса
func (sf *dictServe) SendRequest(topic ATopic, binReq []byte) ([]byte, error) {
	req := &msg_serve.ServeReq{
		Topic_:  topic,
		BinReq_: binReq,
	}
	sf.chSendRequestIn <- req
	resp := <-sf.chSendRequestOut
	return resp.binResp, resp.err
}

func (sf *dictServe) run() {
	for {
		select {
		case handler := <-sf.chUnregisterIn:
			delete(sf.dictServe, handler.Topic())
		case reqServe := <-sf.chSendRequestIn:
			binResp, err := sf.sendRequest(reqServe)
			resp := &serveResp{
				err:     err,
				binResp: binResp,
			}
			sf.chSendRequestOut <- resp
		case handler := <-sf.chRegisterIn:
			sf.chRegisterOut <- sf.register(handler)
		}
	}
}

var TimeoutDefault = 15000

// вызывает обработчик при поступлении запроса
func (sf *dictServe) sendRequest(req *msg_serve.ServeReq) ([]byte, error) {
	handler, isOk := sf.dictServe[req.Topic_]
	if !isOk {
		return nil, fmt.Errorf("dictServe.sendRequest(): handler for topic (%v) not exists", req.Topic_)
	}
	var (
		chErr  = make(chan error, 2)
		binRes []byte
	)
	ctx, fnCancel := context.WithTimeout(sf.ctx.BaseCtx(), time.Millisecond*time.Duration(TimeoutDefault))
	defer fnCancel()
	fnCall := func() {
		defer close(chErr)
		var err error
		binRes, err = handler.FnBack(req.BinReq_)
		if err != nil {
			chErr <- err
		}
	}
	go fnCall()
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("dictServe.sendRequest(): in call for topic (%v), err=\n\t%w", req.Topic_, ctx.Err())
	case err := <-chErr:
		if err != nil {
			return nil, fmt.Errorf("dictServe.sendRequest(): error in call for topic (%v), err=\n\t%w", req.Topic_, err)
		}
	}
	return binRes, nil
}

// регистрирует обработчик запросов
func (sf *dictServe) register(handler IBusHandlerServe) bool {
	topic := handler.Topic()
	_, isOk := sf.dictServe[topic]
	if isOk {
		return true
	}
	sf.dictServe[topic] = handler
	return false
}
