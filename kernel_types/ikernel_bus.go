package kernel_types

import (
	. "github.com/prospero78/kern/kernel_alias"
)

// IBusBaseHandler -- базовый обработчик обратного вызова
type IBusBaseHandler interface {
	// Topic -- топик подписки обработчика
	Topic() ATopic
	// Name -- уникальное имя обработчика
	Name() AHandlerName
}

// IBusHandlerSubscribe -- объект обработчика подписки
type IBusHandlerSubscribe interface {
	IBusBaseHandler
	// FnBack -- функция обратного вызова
	FnBack([]byte)
}

// IBusHandlerServe -- обработчик входящих запросов
type IBusHandlerServe interface {
	IBusBaseHandler
	// FnBack -- функция обратного вызова
	FnBack(binReq []byte) (binResp []byte, err error)
}

// IDictSubHook -- словарь обработчиков по единственному топик
type IDictSubHook interface {
	// Subscribe -- подписывает обработчик
	Subscribe(IBusHandlerSubscribe)
	// Read -- все локальные обработчики читают сообщение по его приходу
	Read(binMsg []byte)
	// Unsubscribe -- отписывает обработчик
	Unsubscribe(IBusHandlerSubscribe)
}

// IDictTopicSub -- интерфейс к словарю обработчиков подписки на словарь топиков
//
//	При подписке потребителей топика может быть НЕСКОЛЬКО на КАЖДЫЙ топик
type IDictTopicSub interface {
	// Subscribe -- подписывает подписчиков на любой из топиков
	Subscribe(IBusHandlerSubscribe)
	// Read -- читает сообщение для всех обработчиков подписки по приходу на любой из топиков
	Read(topic ATopic, binMsg []byte)
	// Unsubscribe -- отписывает подписчиков от любого из топиков
	Unsubscribe(IBusHandlerSubscribe)
}

// IDictTopicServe -- интерфейс к обработчику входящих запросов на словарь топиков
//
// При обслуживании входящих запросов обработчик может быть только ОДИН на КАЖДЫЙ топик.
// Но обработчик вызывается конкурентно.
type IDictTopicServe interface {
	// Register -- регистрирует единственный обработчик на единственный топик
	Register(IBusHandlerServe)
	// SendRequest -- выполняет запрос по указанному топику
	SendRequest(topic ATopic, binReq []byte) (binResp []byte, errResp error)
	// Unregister -- удаляет единственный обработчик с единственного топика
	Unregister(IBusHandlerServe)
}

// IKernelBus -- шина сообщений ядра
//
//	Публикация и запрос требуют параметров на _передачу_.
//	Подписка и обслуживание входящих запросов требует _обработчиков_.
type IKernelBus interface {
	// Publish -- публикует сообщение в шину
	Publish(topic ATopic, binMsg []byte) error
	// SendRequest -- выполняет запрос по указанному топику
	SendRequest(topic ATopic, binReq []byte) (binResp []byte, errResp error)

	// Subscribe -- подписывает обработчик на топик
	Subscribe(IBusHandlerSubscribe) error
	// Unsubscribe -- отписывается от топика
	Unsubscribe(IBusHandlerSubscribe)
	// RegisterServe -- Регистрирует обработчик на обслуживание входящих запросов
	RegisterServe(IBusHandlerServe)

	// IsWork -- возвращает признак работы шины
	IsWork() bool
}
