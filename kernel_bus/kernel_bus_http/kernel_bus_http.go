// package kernel_bus_http -- шина сообщений поверх HTTP
package kernel_bus_http

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	. "github.com/prospero78/kern/helpers"
	. "github.com/prospero78/kern/kernel_alias"
	"github.com/prospero78/kern/kernel_bus/kernel_bus_base"
	"github.com/prospero78/kern/kernel_bus/kernel_bus_http/handler_http_sub"
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
	fibApp.Post("/bus/request", bus.postSendRequest)
	fibApp.Post("/bus/pub", bus.postPublish)
	return bus
}

// SubscribeReq -- входящий запрос на подписку
type SubscribeReq struct {
	Topic_   ATopic `json:"topic"` // Топик, на который надо подписаться
	Uuid_    string `json:"uuid"`
	WebHook_ string `json:"web_hook"` // Веб-хук для обратного вызова
}

// SelfCheck -- проверяет поля на правильность
func (sf *SubscribeReq) SelfCheck() {
	Hassert(sf.Topic_ != "", "SubscribeReq.SelfCheck(): topic is empty")
	Hassert(sf.Uuid_ != "", "SubscribeReq.SelfCheck(): uuid is empty")
	Hassert(sf.WebHook_ != "", "SubscribeReq.SelfCheck(): WebHook_ is empty")
}

// SubscribeResp -- ответ на запрос подписки
type SubscribeResp struct {
	Status_ string `json:"status"`
	Uuid_   string `json:"uuid"`
}

// Входящий запрос HTTP на подписку
func (sf *kernelBusHttp) postSub(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/html; charset=utf8")
	ctx.Set("Content-type", "text/json")
	ctx.Set("Cache-Control", "no-cache")
	req := &SubscribeReq{}
	err := ctx.BodyParser(req)
	if err != nil {
		resp := &SubscribeResp{
			Status_: fmt.Sprintf("kernelBusHttp.postSub(): in parse request, err=\n\t%v\n", err),
			Uuid_:   req.Uuid_,
		}
		ctx.Response().SetStatusCode(http.StatusBadRequest)
		return ctx.JSON(resp)
	}
	resp := sf.processSubscribe(req)
	return ctx.JSON(resp)
}

// Процесс подписки веб-хука
func (sf *kernelBusHttp) processSubscribe(req *SubscribeReq) *SubscribeResp {
	req.SelfCheck()
	handler := handler_http_sub.NewHandlerHttpSub(req.Topic_, req.WebHook_)
	resp := &SubscribeResp{
		Status_: "ok",
		Uuid_:   req.Uuid_,
	}
	err := sf.Subscribe(handler)
	if err != nil {
		resp.Status_ = fmt.Sprintf("kernelBusHttp.processSubscribe(): err=\n\t%v", err)
		return resp
	}
	return resp
}

// Входящая публикация
func (sf *kernelBusHttp) postPublish(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusInternalServerError)
}

// ServeReq -- входящий запрос на обслуживание
type ServeReq struct {
	Topic_  ATopic `json:"topic"`
	Uuid_   string `json:"uuid"`
	BinReq_ []byte `json:"req"`
}

// SelfCheck -- проверяет структуру на правильность полей
func (sf *ServeReq) SelfCheck() {
	Hassert(sf.Topic_ != "", "ServeReq.SelfCheck(): topic is empty")
	Hassert(sf.Uuid_ != "", "ServeReq.SelfCheck(): uuid is empty")
}

// ServeResp -- ответ на входящий запрос
type ServeResp struct {
	Status_  string `json:"status"`
	Uuid_    string `json:"uuid"`
	BinResp_ []byte `json:"resp"`
}

// Входящий запрос
func (sf *kernelBusHttp) postSendRequest(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/html; charset=utf8")
	ctx.Set("Content-type", "text/json")
	ctx.Set("Cache-Control", "no-cache")
	req := &ServeReq{}
	err := ctx.BodyParser(req)
	if err != nil {
		resp := &ServeResp{
			Status_: fmt.Sprintf("kernelBusHttp.postSendRequest(): err=\n\t%v", err),
			Uuid_:   req.Uuid_,
		}
		ctx.Response().SetStatusCode(http.StatusBadRequest)
		return ctx.JSON(resp)
	}
	resp := sf.processSendRequest(req)
	return ctx.JSON(resp)
}

// Обрабатывает входящий запрос
func (sf *kernelBusHttp) processSendRequest(req *ServeReq) *ServeResp {
	req.SelfCheck()
	binResp, err := sf.SendRequest(req.Topic_, req.BinReq_)
	resp := &ServeResp{
		Status_:  "ok",
		Uuid_:    req.Uuid_,
		BinResp_: binResp,
	}
	if err != nil {
		resp.Status_ = fmt.Sprintf("kernelBusHttp.processSendRequest(): err=\n\t%v", err)
	}

	return resp
}

// Входящая отписка от топика по HTTP
func (sf *kernelBusHttp) postUnsub(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusInternalServerError)
}
