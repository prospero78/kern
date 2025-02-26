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
	ctx.Set("kernBus", bus)
	fibApp := kernel_serv_http.GetKernelServHttp().Fiber()
	fibApp.Post("/bus/sub", bus.postSub)             // Топик подписки, IN
	fibApp.Post("/bus/unsub", bus.postUnsub)         // Топик отписки, IN
	fibApp.Post("/bus/request", bus.postSendRequest) // Топик входящих запросов, IN
	fibApp.Post("/bus/pub", bus.postPublish)         // Топик публикаций подписки, IN
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
	Name_   string `json:"name"` // Уникальное имя подписки
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
		Name_:   handler.Name(),
	}
	err := sf.Subscribe(handler)
	if err != nil {
		resp.Status_ = fmt.Sprintf("kernelBusHttp.processSubscribe(): err=\n\t%v", err)
		return resp
	}
	return resp
}

// PublishReq -- запрос на публикацию
type PublishReq struct {
	Topic_ ATopic `json:"topic"`
	Uuid_  string `json:"uuid"`
	BinMsg []byte `json:"msg"`
}

// SelfCheck -- проверяет правильность своих полей
func (sf *PublishReq) SelfCheck() {
	Hassert(sf.Topic_ != "", "PublishReq.SelfCheck(): topic is empty")
	Hassert(sf.Uuid_ != "", "PublishReq.SelfCheck(): uuid is empty")
}

// PublishResp -- ответ на запрос публикации
type PublishResp struct {
	Status_ string `json:"status"`
	Uuid_   string `json:"uuid"`
}

// Входящая публикация
func (sf *kernelBusHttp) postPublish(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/html; charset=utf8")
	ctx.Set("Content-type", "text/json")
	ctx.Set("Cache-Control", "no-cache")
	req := &PublishReq{}
	err := ctx.BodyParser(req)
	if err != nil {
		resp := &SubscribeResp{
			Status_: fmt.Sprintf("kernelBusHttp.postPublish(): in parse request, err=\n\t%v\n", err),
			Uuid_:   req.Uuid_,
		}
		ctx.Response().SetStatusCode(http.StatusBadRequest)
		return ctx.JSON(resp)
	}
	resp := sf.processPublish(req)
	return ctx.JSON(resp)
}

// Выполняет процесс публикации
func (sf *kernelBusHttp) processPublish(req *PublishReq) *PublishResp {
	req.SelfCheck()
	err := sf.Publish(req.Topic_, req.BinMsg)
	resp := &PublishResp{
		Status_: "ok",
		Uuid_:   req.Uuid_,
	}
	if err != nil {
		resp.Status_ = fmt.Sprintf("kernelBusHttp.processPublish(): err=\n\t%v", err)
	}
	return resp
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

// UnsubReq -- запрос на отписку от топика
type UnsubReq struct {
	Name_ string `json:"name"` // Уникальная метка подписки
	Uuid_ string `json:"uuid"`
}

// SelfCheck -- проверка запроса на правильность полей
func (sf *UnsubReq) SelfCheck() {
	Hassert(sf.Name_ != "", "UnsubReq.SelfCheck(): name is empty")
	Hassert(sf.Uuid_ != "", "UnsubReq.SelfCheck(): uuid is empty")
}

// UnsubResp -- ответ на запрос отписки
type UnsubResp struct {
	Status_ string `json:"status"`
	Uuid_   string `json:"uuid"`
}

// Входящая отписка от топика по HTTP
func (sf *kernelBusHttp) postUnsub(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/html; charset=utf8")
	ctx.Set("Content-type", "text/json")
	ctx.Set("Cache-Control", "no-cache")
	req := &UnsubReq{}
	err := ctx.BodyParser(req)
	if err != nil {
		resp := &ServeResp{
			Status_: fmt.Sprintf("kernelBusHttp.postSendRequest(): err=\n\t%v", err),
			Uuid_:   req.Uuid_,
		}
		ctx.Response().SetStatusCode(http.StatusBadRequest)
		return ctx.JSON(resp)
	}
	resp := sf.processUnsubRequest(req)
	return ctx.JSON(resp)
}

// Процесс отписки от топика
func (sf *kernelBusHttp) processUnsubRequest(req *UnsubReq) *UnsubResp {
	req.SelfCheck()
	_hand := sf.Ctx_.Get(req.Name_)
	resp := &UnsubResp{
		Status_: "ok",
		Uuid_:   req.Uuid_,
	}
	if _hand == nil {
		resp.Status_ = fmt.Sprintf("kernelBusHttp.processUnsubRequest(): handler name(%v) not exists", req.Name_)
		return resp
	}
	hand := _hand.(IBusHandlerSubscribe)
	sf.Unsubscribe(hand)
	return resp
}
