// package kernel_bus_http -- шина сообщений поверх HTTP
package kernel_bus_http

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	. "github.com/prospero78/kern/kernel_alias"
	"github.com/prospero78/kern/kernel_bus/kernel_bus_base"
	"github.com/prospero78/kern/kernel_ctx"
	"github.com/prospero78/kern/kernel_serv_http"
	. "github.com/prospero78/kern/kernel_types"
)

// kernelBusHttp -- шина данных поверх HTTP
type kernelBusHttp struct {
	*kernel_bus_base.KernelBusBase
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
	bus = &kernelBusHttp{
		KernelBusBase: kernel_bus_base.GetKernelBusBase(),
	}
	ctx.Add("kernBus", bus)
	fibApp := kernel_serv_http.GetKernelServHttp().Fiber()
	fibApp.Post("/bus/sub", bus.postSub)
	fibApp.Post("/bus/unsub", bus.postUnsub)
	fibApp.Post("/bus/request", bus.postRequest)
	fibApp.Post("/bus/pub", bus.postPublish)
	return bus
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
	err = sf.Subscribe(handler)
	if err != nil {
		dict := map[string]string{
			"error": fmt.Sprintf("kernelBusHttp.postSub(): in subscribe request, err=\n\t%v\n", err),
		}
		ctx.Response().SetStatusCode(http.StatusInternalServerError)
		return ctx.JSON(dict)
	}
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
