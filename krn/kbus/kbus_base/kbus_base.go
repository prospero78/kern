// package kbus_base -- базовая часть шины данных
package kbus_base

import (
	"fmt"
	"log"
	"sync"

	. "github.com/prospero78/kern/kc/helpers"
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
	Ctx_      IKernelCtx
	IsWork_   ISafeBool
	dictSub   IDictTopicSub
	dictServe IDictTopicServe
	block     sync.Mutex
}

var (
	Bus_ *KBusBase
)

// GetKernelBusBase -- возвращает базовую шину сообщений
func GetKernelBusBase() *KBusBase {
	if Bus_ != nil {
		return Bus_
	}
	ctx := kctx.GetKernelCtx()
	Bus_ = &KBusBase{
		Ctx_:      ctx,
		IsWork_:   safe_bool.NewSafeBool(),
		dictSub:   dict_topic_sub.NewDictTopicSub(),
		dictServe: dict_topic_serve.NewDictServe(),
	}
	go Bus_.close()
	err := Bus_.Ctx_.Wg().Add(strBusBaseStream)
	Hassert(err == nil, "GetKernelBusBase(): in add name stream '%v' Wg, err=%v", strBusBaseStream, err)
	Bus_.IsWork_.Set()
	ctx.Set("kernBusBase", Bus_, "base of data bus")
	_ = IKernelBus(Bus_)
	return Bus_
}

// Unsubscribe -- отписывает обработчик от топика
func (sf *KBusBase) Unsubscribe(handler IBusHandlerSubscribe) {
	sf.dictSub.Unsubscribe(handler)
}

// Subscribe -- подписывает обработчик на топик
func (sf *KBusBase) Subscribe(handler IBusHandlerSubscribe) error {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.IsWork_.Get() {
		return fmt.Errorf("KernelBusBase.Subscribe(): bus already closed")
	}
	sf.dictSub.Subscribe(handler)
	return nil
}

// SendRequest -- отправляет запрос в шину данных
func (sf *KBusBase) SendRequest(topic ATopic, binReq []byte) ([]byte, error) {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.IsWork_.Get() {
		return nil, fmt.Errorf("KernelBusBase.SendRequest(): bus already closed")
	}
	binResp, err := sf.dictServe.SendRequest(topic, binReq)
	if err != nil {
		return nil, fmt.Errorf("KernelBusBase.SendRequest(): topic='%v', err=\n\t%w", topic, err)
	}
	return binResp, nil
}

// RegisterServe -- регистрирует обработчики входящих запросов
func (sf *KBusBase) RegisterServe(handler IBusHandlerServe) {
	Hassert(handler != nil, "KernelBusBase.Subscribe(): IBusHandlerSubscribe==nil")
	sf.dictServe.Register(handler)
}

// Publish -- публикует сообщение в шину
func (sf *KBusBase) Publish(topic ATopic, binMsg []byte) (err error) {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.IsWork_.Get() {
		return fmt.Errorf("KernelBusBase.Publish(): bus already closed")
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
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.IsWork_.Get() {
		return
	}
	sf.IsWork_.Reset()
	sf.Ctx_.Wg().Done(strBusBaseStream)
	log.Println("kernelBusLocal.close(): done")
}
