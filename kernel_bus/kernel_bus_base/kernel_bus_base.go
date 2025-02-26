// package kernel_bus_base -- базовая часть шины данных
package kernel_bus_base

import (
	"fmt"
	"log"
	"sync"

	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	"github.com/prospero78/kern/kernel_bus/dict_topic_serve"
	"github.com/prospero78/kern/kernel_bus/dict_topic_sub"
	"github.com/prospero78/kern/kernel_ctx"
	. "github.com/prospero78/kern/kernel_types"
	"github.com/prospero78/kern/safe_bool"
)

const (
	strBusBaseStream = "bus_base"
)

// KernelBusBase -- базовая часть шины данных
type KernelBusBase struct {
	ctx       IKernelCtx
	isWork    ISafeBool
	dictSub   IDictTopicSub
	dictServe IDictTopicServe
	block     sync.Mutex
}

var (
	bus *KernelBusBase
)

// GetKernelBusBase -- возвращает базовую шину сообщений
func GetKernelBusBase() *KernelBusBase {
	if bus != nil {
		return bus
	}
	ctx := kernel_ctx.GetKernelCtx()
	bus = &KernelBusBase{
		ctx:       ctx,
		isWork:    safe_bool.NewSafeBool(),
		dictSub:   dict_topic_sub.NewDictTopicSub(),
		dictServe: dict_topic_serve.NewDictServe(),
	}
	go bus.close()
	err := bus.ctx.Wg().Add(strBusBaseStream)
	Hassert(err == nil, "GetKernelBusBase(): in add name stream '%v' Wg, err=%v", strBusBaseStream, err)
	bus.isWork.Set()
	ctx.Add("kernBusBase", bus)
	_ = IKernelBus(bus)
	return bus
}

// Unsubscribe -- отписывает обработчик от топика
func (sf *KernelBusBase) Unsubscribe(handler IBusHandlerSubscribe) {
	sf.dictSub.Unsubscribe(handler)
}

// Subscribe -- подписывает обработчик на топик
func (sf *KernelBusBase) Subscribe(handler IBusHandlerSubscribe) error {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return fmt.Errorf("KernelBusBase.Subscribe(): bus already closed")
	}
	sf.dictSub.Subscribe(handler)
	return nil
}

// SendRequest -- отправляет запрос в шину данных
func (sf *KernelBusBase) SendRequest(topic ATopic, binReq []byte) ([]byte, error) {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return nil, fmt.Errorf("KernelBusBase.SendRequest(): bus already closed")
	}
	binResp, err := sf.dictServe.SendRequest(topic, binReq)
	if err != nil {
		return nil, fmt.Errorf("KernelBusBase.SendRequest(): topic='%v', err=\n\t%w", topic, err)
	}
	return binResp, nil
}

// RegisterServe -- регистрирует обработчики входящих запросов
func (sf *KernelBusBase) RegisterServe(handler IBusHandlerServe) {
	Hassert(handler != nil, "KernelBusBase.Subscribe(): IBusHandlerSubscribe==nil")
	sf.dictServe.Register(handler)
}

// Publish -- публикует сообщение в шину
func (sf *KernelBusBase) Publish(topic ATopic, binMsg []byte) (err error) {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return fmt.Errorf("KernelBusBase.Publish(): bus already closed")
	}
	// Асинхронный запуск чтения
	go sf.dictSub.Read(topic, binMsg)
	return nil
}

// IsWork -- возвращает признак работы шины
func (sf *KernelBusBase) IsWork() bool {
	return sf.isWork.Get()
}

// Ожидает закрытия шины в отдельном потоке
func (sf *KernelBusBase) close() {
	<-sf.ctx.Ctx().Done()
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return
	}
	sf.isWork.Reset()
	sf.ctx.Wg().Done(strBusBaseStream)
	log.Println("kernelBusLocal.close(): done")
}
