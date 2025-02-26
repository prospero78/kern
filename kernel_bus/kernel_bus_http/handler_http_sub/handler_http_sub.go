// package handler_http_sub -- обработчик подписки по HTTP
package handler_http_sub

import (
	"bytes"
	"log"
	"net/http"

	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	. "github.com/prospero78/kern/kernel_types"
)

// handlerHttpSub -- обработчик подписки по HTTP
type handlerHttpSub struct {
	name  string // Уникальное имя обработчика (webHook)
	topic ATopic // Имя топика, на который подписан обработчик
}

// NewHandlerHttpSub -- возвращает новый обработчик подписки по HTTP
func NewHandlerHttpSub(topic ATopic, webHook string) IBusHandlerSubscribe {
	Hassert(topic != "", "NewHandlerHttpSub(): topic is empty")
	Hassert(webHook != "", "NewHandlerHttpSub(): webHook is empty")
	sf := &handlerHttpSub{
		topic: topic,
		name:  webHook,
	}
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
func (sf *handlerHttpSub) FnBack(binMsg []byte) {
	body := bytes.NewBuffer(binMsg)
	resp, err := http.Post(sf.name, "application/json", body)
	if err != nil {
		log.Printf("handlerHttpSub.FnBack(): topic='%v', in make request, err=\n\t%v\n", sf.name, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("handlerHttpSub.FnBack(): topic='%v', code=%v, status=%v\n", sf.name, resp.StatusCode, resp.Status)
	}
}
