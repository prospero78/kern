// package handler_http_sub -- обработчик подписки по HTTP
package handler_http_sub

import (
	"bytes"
	"crypto/rand"
	"log"
	"net/http"

	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	. "github.com/prospero78/kern/kernel_types"
)

// handlerHttpSub -- обработчик подписки по HTTP
type handlerHttpSub struct {
	name    AHandlerName // Уникальное имя обработчика
	topic   ATopic       // Имя топика, на который подписан обработчик
	webHook string       // Куда обращаться при запросах
}

// NewHandlerHttpSub -- возвращает новый обработчик подписки по HTTP
func NewHandlerHttpSub(topic ATopic, webHook string) IBusHandlerSubscribe {
	Hassert(topic != "", "NewHandlerHttpSub(): topic is empty")
	Hassert(webHook != "", "NewHandlerHttpSub(): webHook is empty")
	sf := &handlerHttpSub{
		topic:   topic,
		name:    AHandlerName(webHook + "_" + rand.Text()),
		webHook: webHook,
	}
	return sf
}

// Topic -- возвращает имя топика, на который подписан обработчик
func (sf *handlerHttpSub) Topic() ATopic {
	return sf.topic
}

// Name -- возвращает уникальное имя обработчика
func (sf *handlerHttpSub) Name() AHandlerName {
	return sf.name
}

// FnBack -- обратный вызов по приходу сообщения
func (sf *handlerHttpSub) FnBack(binMsg []byte) {
	body := bytes.NewBuffer(binMsg)
	resp, err := http.Post(sf.webHook, "application/json", body)
	if err != nil {
		log.Printf("handlerHttpSub.FnBack(): topic='%v', in make request, err=\n\t%v\n", sf.webHook, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("handlerHttpSub.FnBack(): topic='%v', code=%v, status=%v\n", sf.webHook, resp.StatusCode, resp.Status)
	}
}
