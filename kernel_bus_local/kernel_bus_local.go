// package kernel_bus_local -- реализация локальной шины сообщений
package kernel_bus_local

import (
	"context"
	"fmt"
	"log"
	"sync"

	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	"github.com/prospero78/kern/kernel_bus_local/dict_serve"
	"github.com/prospero78/kern/kernel_bus_local/dict_sub"
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
	dictSub   *dict_sub.DictSub
	dictServe *dict_serve.DictServe
	block     sync.Mutex
}

var (
	bus *kernelBusLocal
)

// GetKernelBusLocal -- возвращает локальную шину сообщений
func GetKernelBusLocal(ctx IKernelCtx) IKernelBus {
	if bus != nil {
		return bus
	}
	Hassert(ctx != nil, "GetKernelBusLocal(): IKernelCtx==nil")
	bus = &kernelBusLocal{
		ctx:       ctx,
		isWork:    safe_bool.NewSafeBool(),
		dictSub:   dict_sub.NewDictSub(),
		dictServe: dict_serve.NewDictServe(ctx),
	}
	go bus.close()
	err := bus.ctx.Wg().Add(strBusLocalStream)
	Hassert(err == nil, "GetKernelBusLocal(): in add name stream '%v' Wg, err=%v", strBusLocalStream, err)
	bus.isWork.Set()
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
	sf.dictSub.Add(handler)
	return nil
}

// Unsubscribe -- отписывает обработчик от топика
func (sf *kernelBusLocal) Unsubscribe(handler IBusHandlerSubscribe) {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.dictSub.Del(handler)
}

// Serve -- выполняет обслуживание входящих запросов
func (sf *kernelBusLocal) Serve(handler IBusHandlerServe) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(handler != nil, "kernelBusLocal.Subscribe(): IBusHandlerSubscribe==nil")
	sf.dictServe.Add(handler)
}

// Request -- выполняет запрос в шину данных
func (sf *kernelBusLocal) Request(ctx context.Context, topic ATopic, binReq []byte) ([]byte, error) {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return nil, fmt.Errorf("kernelBusLocal.Request(): bus already closed")
	}
	binResp, err := sf.dictServe.Call(ctx, topic, binReq)
	if err != nil {
		return nil, fmt.Errorf("kernelBusLocal.Request(): topic='%v', err=\n\t%w", topic, err)
	}
	return binResp, nil
}

// Publish -- публикует сообщение в шину
func (sf *kernelBusLocal) Publish(ctx context.Context, topic ATopic, binMsg []byte) (err error) {
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return fmt.Errorf("kernelBusLocal.Publish(): bus already closed")
	}
	go sf.dictSub.Call(topic, binMsg)
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
