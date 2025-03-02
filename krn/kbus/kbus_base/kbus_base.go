// package kbus_base -- базовая часть шины данных
package kbus_base

import (
	"fmt"
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/local_ctx"
	"github.com/prospero78/kern/kc/safe_bool"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kbus/dict_topic_serve"
	"github.com/prospero78/kern/krn/kbus/dict_topic_sub"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

const (
	strBusBaseStream = "bus_base"
)

// KBusBase -- базовая часть шины данных
type KBusBase struct {
	sync.Mutex
	Ctx_      IKernelCtx
	IsWork_   ISafeBool
	ctx       ILocalCtx
	log       ILogBuf
	dictSub   IDictTopicSub
	dictServe IDictTopicServe
}

var (
	Bus_  *KBusBase
	block sync.Mutex
)

// GetKernelBusBase -- возвращает базовую шину сообщений
func GetKernelBusBase() *KBusBase {
	block.Lock()
	defer block.Unlock()
	if Bus_ != nil {
		return Bus_
	}
	ctx := kctx.GetKernelCtx()
	Bus_ = &KBusBase{
		Ctx_:      ctx,
		IsWork_:   safe_bool.NewSafeBool(),
		dictSub:   dict_topic_sub.NewDictTopicSub(),
		dictServe: dict_topic_serve.NewDictServe(),
		ctx:       local_ctx.NewLocalCtx(ctx.BaseCtx()),
	}
	Bus_.log = Bus_.ctx.Log()
	go Bus_.close()
	go Bus_.run()
	Bus_.IsWork_.Set()
	err := Bus_.Ctx_.Wg().Add(strBusBaseStream)
	Hassert(err == nil, "GetKernelBusBase(): in add name stream(%v), err=\n\t%v", strBusBaseStream, err)
	ctx.Set("kernBusBase", Bus_, "base of data bus")
	_ = IKernelBus(Bus_)
	return Bus_
}

// Log -- возвращает лог шины
func (sf *KBusBase) Log() ILogBuf {
	return sf.log
}

func (sf *KBusBase) run() {
	sf.log.Debug("KBusBase.run()")
	for {
		break
	}
}

// Unsubscribe -- отписывает обработчик от топика
func (sf *KBusBase) Unsubscribe(handler IBusHandlerSubscribe) {
	sf.log.Debug("KBusBase.Unsubscribe(): handler='%v'", handler.Name())
	sf.dictSub.Unsubscribe(handler)
}

// Subscribe -- подписывает обработчик на топик
func (sf *KBusBase) Subscribe(handler IBusHandlerSubscribe) error {
	sf.Lock()
	defer sf.Unlock()
	sf.log.Debug("KBusBase.Subscribe(): handler='%v'", handler.Name())
	if !sf.IsWork_.Get() {
		err := fmt.Errorf("KBusBase.Subscribe():  handler='%v', bus already closed", handler.Name())
		sf.log.Err(err.Error())
		return err
	}
	sf.dictSub.Subscribe(handler)
	return nil
}

// SendRequest -- отправляет запрос в шину данных
func (sf *KBusBase) SendRequest(topic ATopic, binReq []byte) ([]byte, error) {
	sf.Lock()
	defer sf.Unlock()
	sf.log.Debug("KBusBase.SendRequest(): topic='%v'", topic)
	if !sf.IsWork_.Get() {
		err := fmt.Errorf("KBusBase.SendRequest():  topic='%v', bus already closed", topic)
		sf.log.Err(err.Error())
		return nil, err
	}
	binResp, err := sf.dictServe.SendRequest(topic, binReq)
	if err != nil {
		err := fmt.Errorf("KBusBase.SendRequest(): topic='%v', err=\n\t%w", topic, err)
		sf.log.Err(err.Error())
		return nil, err
	}
	return binResp, nil
}

// RegisterServe -- регистрирует обработчики входящих запросов
func (sf *KBusBase) RegisterServe(handler IBusHandlerServe) {
	Hassert(handler != nil, "KBusBase.RegisterServe(): IBusHandlerSubscribe==nil")
	sf.log.Debug("KBusBase.RegisterServe(): handler='%v'", handler.Name())
	sf.dictServe.Register(handler)
}

// Publish -- публикует сообщение в шину
func (sf *KBusBase) Publish(topic ATopic, binMsg []byte) (err error) {
	sf.Lock()
	defer sf.Unlock()
	sf.log.Debug("KBusBase.Publish(): topic='%v'", topic)
	if !sf.IsWork_.Get() {
		err := fmt.Errorf("KBusBase.Publish(): topic='%v',bus already closed", topic)
		sf.log.Err(err.Error())
		return err
	}
	// Асинхронный запуск чтения
	go sf.dictSub.Read(topic, binMsg)
	return nil
}

// IsWork -- возвращает признак работы шины
func (sf *KBusBase) IsWork() bool {
	return sf.IsWork_.Get()
}

// Ожидает закрытия шины в отдельном потоке
func (sf *KBusBase) close() {
	sf.Ctx_.Done()
	sf.Ctx_.Lock()
	defer sf.Ctx_.Unlock()
	if !sf.IsWork_.Get() {
		return
	}
	sf.IsWork_.Reset()
	sf.Ctx_.Wg().Done(strBusBaseStream)
	sf.log.Debug("KBusBase.close(): done")
}
