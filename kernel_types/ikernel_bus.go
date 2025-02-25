package kernel_types

import (
	"context"

	. "github.com/prospero78/kern/kernel_alias"
)

// iBusBaseHandler -- базовый обработчик обратного вызова
type iBusBaseHandler interface {
	// Topic -- топик подписки обработчика
	Topic() ATopic
}

// IBusHandlerSubscribe -- объект обработчика подписки
type IBusHandlerSubscribe interface {
	iBusBaseHandler
	// FnBack -- функция обратного вызова
	FnBack([]byte)
}

// IBusHandlerServe -- обработчик входящих запросов
type IBusHandlerServe interface {
	iBusBaseHandler
	// FnBack -- функция обратного вызова
	FnBack(binReq []byte) (binResp []byte, err error)
}

// IKernelBus -- шина сообщений ядра
type IKernelBus interface {
	// Publish -- публикует сообщение в шину
	Publish(ctx context.Context, topic ATopic, binMsg []byte) error
	// Subscribe -- подписывает обработчик на топик
	Subscribe(IBusHandlerSubscribe) error
	// Unsubscribe -- отписывается от топика
	Unsubscribe(IBusHandlerSubscribe)
	// Serve -- обслуживает входящие запросы
	Serve(IBusHandlerServe)
	// Request -- выполняет запрос по указанному топику
	Request(ctx context.Context, topic ATopic, binReq []byte) (binResp []byte, errResp error)
	// IsWork -- возвращает признак работы шины
	IsWork() bool
}
