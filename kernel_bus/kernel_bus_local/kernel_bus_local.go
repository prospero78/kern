// package kernel_bus_local -- реализация локальной шины сообщений
package kernel_bus_local

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
	strBusLocalStream = "bus_local"
)

// Локальная шина данных
type kernelBusLocal struct {
	ctx       IKernelCtx
	isWork    ISafeBool
	dictSub   IDictTopicSub
	dictServe IDictTopicServe
	block     sync.Mutex
}

var (
	bus *kernelBusLocal
)

// GetKernelBusLocal -- возвращает локальную шину сообщений
func GetKernelBusLocal() IKernelBus {
	if bus != nil {
		return bus
	}
	ctx := kernel_ctx.GetKernelCtx()
	bus = &kernelBusLocal{
		ctx:       ctx,
		isWork:    safe_bool.NewSafeBool(),
		dictSub:   dict_topic_sub.NewDictTopicSub(),
		dictServe: dict_topic_serve.NewDictServe(),
	}
	go bus.close()
	err := bus.ctx.Wg().Add(strBusLocalStream)
	Hassert(err == nil, "GetKernelBusLocal(): in add name stream '%v' Wg, err=%v", strBusLocalStream, err)
	bus.isWork.Set()
	ctx.Add("kernBus", bus)
	return bus
}

// IsWork -- возвращает признак работы шины
func (sf *kernelBusLocal) IsWork() bool {
	return sf.isWork.Get()
}

// Subscribe -- подписывает обработчик на топик
func (sf *kernelBusLocal) Subscribe(handler IBusHandlerSubscribe) error {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return fmt.Errorf("kernelBusLocal.Subscribe(): bus already closed")
	}
	sf.dictSub.Subscribe(handler)
	return nil
}

// Unsubscribe -- отписывает обработчик от топика
func (sf *kernelBusLocal) Unsubscribe(handler IBusHandlerSubscribe) {
	sf.dictSub.Unsubscribe(handler)
}

// RegisterServe -- регистрирует обработчики входящих запросов
func (sf *kernelBusLocal) RegisterServe(handler IBusHandlerServe) {
	Hassert(handler != nil, "kernelBusLocal.Subscribe(): IBusHandlerSubscribe==nil")
	sf.dictServe.Register(handler)
}

// Request -- выполняет запрос в шину данных
func (sf *kernelBusLocal) Request(topic ATopic, binReq []byte) ([]byte, error) {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return nil, fmt.Errorf("kernelBusLocal.Request(): bus already closed")
	}
	binResp, err := sf.dictServe.Request(topic, binReq)
	if err != nil {
		return nil, fmt.Errorf("kernelBusLocal.Request(): topic='%v', err=\n\t%w", topic, err)
	}
	return binResp, nil
}

// Publish -- публикует сообщение в шину
func (sf *kernelBusLocal) Publish(topic ATopic, binMsg []byte) (err error) {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return fmt.Errorf("kernelBusLocal.Publish(): bus already closed")
	}
	// Асинхронный запуск чтения
	go sf.dictSub.Read(topic, binMsg)
	return nil
}

// Ожидает закрытия шины в отдельном потоке
func (sf *kernelBusLocal) close() {
	<-sf.ctx.Ctx().Done()
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return
	}
	sf.isWork.Reset()
	sf.ctx.Wg().Done(strBusLocalStream)
	log.Println("kernelBusLocal.close(): done")
}
