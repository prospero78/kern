// package handler_http_sub -- обработчик подписки по HTTP
package handler_http_sub

import (
	. "github.com/prospero78/kern/kernel_alias"
	. "github.com/prospero78/kern/kernel_types"
)

// handlerHttpSub -- обработчик подписки по HTTP
type handlerHttpSub struct {
	name  string // Уникальное имя обработчика
	topic ATopic // Имя топика, на который подписан обработчик
}

// NewHandlerHttpSub -- возвращает новый обработчик подписки по HTTP
func NewHandlerHttpSub() IBusHandlerSubscribe {
	sf := &handlerHttpSub{}
	return sf
}

// Topic -- возвращает имя топика, на который подписан обработчик
func (sf *handlerHttpSub) Topic() ATopic {
	return sf.topic
}

// Name -- возвращает уникальное имя обработчика
func (sf *handlerHttpSub) Name() string {
	return sf.name
}

// FnBack -- обратный вызов по приходу сообщения
func (sf *handlerHttpSub) FnBack(binMsg []byte) {}
