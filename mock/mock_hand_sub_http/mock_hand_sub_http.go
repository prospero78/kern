// package handler_http_sub -- обработчик подписки по HTTP
package mock_hand_sub_http

import (
	// "bytes"
	"crypto/rand"
	"log"
	"sync"

	// "net/http"

	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
	. "github.com/prospero78/kern/krn/ktypes"
)

// MockHandSubHttp -- обработчик подписки по HTTP
type MockHandSubHttp struct {
	Name_    AHandlerName // Уникальное имя обработчика
	Topic_   ATopic       // Имя топика, на который подписан обработчик
	WebHook_ string       // Куда обращаться при запросах
	BinMsg_  []byte       // Последнее бинарное сообщение
	block    sync.RWMutex
}

// NewMockHandSubHttp -- возвращает новый обработчик подписки по HTTP
func NewMockHandSubHttp(topic ATopic, webHook string) IBusHandlerSubscribe {
	Hassert(topic != "", "NewHandlerHttpSub(): topic is empty")
	Hassert(webHook != "", "NewHandlerHttpSub(): webHook is empty")
	sf := &MockHandSubHttp{
		Topic_:   topic,
		Name_:    AHandlerName(webHook + "_" + rand.Text()),
		WebHook_: webHook,
	}
	return sf
}

// Topic -- возвращает имя топика, на который подписан обработчик
func (sf *MockHandSubHttp) Topic() ATopic {
	return sf.Topic_
}

// SetName -- устанавливает имя обработчика
func (sf *MockHandSubHttp) SetName(name AHandlerName) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(name != "", "HandlerHttpSub.SetName(): name is empty")
	sf.Name_ = name
}

// Name -- возвращает уникальное имя обработчика
func (sf *MockHandSubHttp) Name() AHandlerName {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.Name_
}

// FnBack -- обратный вызов по приходу сообщения
func (sf *MockHandSubHttp) FnBack(binMsg []byte) {
	sf.block.Lock()
	defer sf.block.Unlock()
	/*
		body := bytes.NewBuffer(binMsg)
		resp, err := http.Post(sf.WebHook_, "application/json", body)
		if err != nil {
			log.Printf("handlerHttpSub.FnBack(): topic='%v', in make request, err=\n\t%v\n", sf.WebHook_, err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Printf("handlerHttpSub.FnBack(): topic='%v', code=%v, status=%v\n", sf.WebHook_, resp.StatusCode, resp.Status)
		}
	*/
	log.Printf("HandlerHttpSub.FnBack(): msg=%v\n", string(binMsg))
	sf.BinMsg_ = binMsg
}

// Msg -- возвращает строковое представление хранимого сообщения
func (sf *MockHandSubHttp) Msg() string {
	sf.block.Lock()
	defer sf.block.Unlock()
	return string(sf.BinMsg_)
}
