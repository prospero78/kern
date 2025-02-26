// package kernel_bus_http -- шина сообщений поверх HTTP
package kernel_bus_http

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	. "github.com/prospero78/kern/kernel_alias"
	"github.com/prospero78/kern/kernel_bus/dict_topic_serve"
	"github.com/prospero78/kern/kernel_bus/dict_topic_sub"
	"github.com/prospero78/kern/kernel_ctx"
	. "github.com/prospero78/kern/kernel_types"
	"github.com/prospero78/kern/safe_bool"
)

// kernelBusHttp -- шина данных поверх HTTP
type kernelBusHttp struct {
	ctx       IKernelCtx
	isWork    ISafeBool
	block     sync.RWMutex
	kernServ  IKernelServerHttp
	dictSub   IDictTopicSub
	dictServe IDictTopicServe
}

var (
	bus *kernelBusHttp
)

// GetKernelBusHttp -- возвращает шину HTTP
func GetKernelBusHttp() IKernelBus {
	if bus != nil {
		return bus
	}
	ctx := kernel_ctx.GetKernelCtx()
	bus := &kernelBusHttp{
		ctx:       ctx,
		kernServ:  ctx.Get("kernServHttp").(IKernelServerHttp),
		dictSub:   dict_topic_sub.NewDictTopicSub(),
		dictServe: dict_topic_serve.NewDictServe(),
		isWork:    safe_bool.NewSafeBool(),
	}
	ctx.Add("kernBus", bus)
	bus.kernServ.Fiber().Post("/bus/sub", bus.postSub)
	bus.kernServ.Fiber().Post("/bus/unsub", bus.postUnsub)
	bus.kernServ.Fiber().Post("/bus/request", bus.postRequest)
	bus.kernServ.Fiber().Post("/bus/pub", bus.postPublish)
	return bus
}

// Request -- обрабатывает входящий запрос
func (sf *kernelBusHttp) Request(topic ATopic, binMsg []byte) ([]byte, error) {
	return sf.dictServe.Request(topic, binMsg)
}

// RegisterServe -- регистрирует обработчик входящих запросов
func (sf *kernelBusHttp) RegisterServe(handler IBusHandlerServe) {

}

// Publish -- публикует сообщение в шину HTTP
func (sf *kernelBusHttp) Publish(topic ATopic, binMsg []byte) error {
	sf.dictSub.Read(topic, binMsg)
	return nil
}

// RequestSubscribe -- входящий запрос на подписку
type RequestSubscribe struct {
	Topic_   ATopic `json:"topic"`    // Топик, на который надо подписаться
	WebHook_ string `json:"web_hook"` // Веб-хук для обратного вызова
}

// Входящий запрос HTTP на подписку
func (sf *kernelBusHttp) postSub(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/html; charset=utf8")
	ctx.Set("Content-type", "text/json")
	ctx.Set("Cache-Control", "no-cache")
	req := &RequestSubscribe{}
	err := ctx.ParamsParser(req)
	if err != nil {
		dict := map[string]string{
			"error": fmt.Sprintf("kernelBusHttp.postSub(): in parse request, err=\n\t%v\n", err),
		}
		ctx.Response().SetStatusCode(http.StatusBadRequest)
		return ctx.JSON(dict)
	}
	var handler IBusHandlerSubscribe
	sf.dictSub.Subscribe(handler)
	dict := map[string]string{
		"status": "ok",
	}
	return ctx.JSON(dict)
}

// Входящий запрос на публикацию
func (sf *kernelBusHttp) postPublish(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusInternalServerError)
}

// Входящие запрос на обслуживание
func (sf *kernelBusHttp) postRequest(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusInternalServerError)
}

// Входящий запрос HTTP на отписку от топика
func (sf *kernelBusHttp) postUnsub(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusInternalServerError)
}

// Unsubscribe -- отписывает подписчика от топика
func (sf *kernelBusHttp) Unsubscribe(handler IBusHandlerSubscribe) {
	sf.dictSub.Unsubscribe(handler)
}

// Subscribe -- подписывает на топик обработчик
func (sf *kernelBusHttp) Subscribe(handler IBusHandlerSubscribe) error {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.dictSub.Subscribe(handler)
	return fmt.Errorf("надо доделать")
}

// IsWork -- возвращает признак работы шины HTTP
func (sf *kernelBusHttp) IsWork() bool {
	return sf.isWork.Get()
}
