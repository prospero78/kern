// package kbus_http -- шина сообщений поверх HTTP
package kbus_http

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/prospero78/kern/krn/kbus/kbus_base"
	"github.com/prospero78/kern/krn/kbus/kbus_http/handler_http_sub"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_pub"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_serve"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_sub"
	"github.com/prospero78/kern/krn/kbus/kbus_msg/msg_unsub"
	"github.com/prospero78/kern/krn/kctx"
	"github.com/prospero78/kern/krn/kserv_http"
	. "github.com/prospero78/kern/krn/ktypes"
)

// kBusHttp -- шина данных поверх HTTP
type kBusHttp struct {
	*kbus_base.KBusBase
	log ILogBuf
}

var (
	bus *kBusHttp
)

// GetKernelBusHttp -- возвращает шину HTTP
func GetKernelBusHttp() IKernelBus {
	if bus != nil {
		return bus
	}
	ctx := kctx.GetKernelCtx()
	bus = &kBusHttp{
		KBusBase: kbus_base.GetKernelBusBase(),
	}
	bus.log = bus.Log()
	ctx.Set("kernBus", bus, "http data bus")
	fibApp := kserv_http.GetKernelServHttp().Fiber()
	fibApp.Post("/bus/sub", bus.postSub)             // Топик подписки, IN
	fibApp.Post("/bus/unsub", bus.postUnsub)         // Топик отписки, IN
	fibApp.Post("/bus/request", bus.postSendRequest) // Топик входящих запросов, IN
	fibApp.Post("/bus/pub", bus.postPublish)         // Топик публикаций подписки, IN
	return bus
}

// Входящий запрос HTTP на подписку
func (sf *kBusHttp) postSub(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/html; charset=utf8")
	ctx.Set("Content-type", "text/json")
	ctx.Set("Cache-Control", "no-cache")
	sf.log.Debug("kBusHttp.postSub()")
	req := &msg_sub.SubscribeReq{}
	err := ctx.BodyParser(req)
	if err != nil {
		resp := &msg_sub.SubscribeResp{
			Status_: fmt.Sprintf("kernelBusHttp.postSub(): in parse request, err=\n\t%v\n", err),
			Uuid_:   req.Uuid_,
		}
		resp.SelfCheck()
		ctx.Response().SetStatusCode(http.StatusBadRequest)
		sf.log.Err(resp.Status_)
		return ctx.JSON(resp)
	}
	resp := sf.processSubscribe(req)
	resp.SelfCheck()
	return ctx.JSON(resp)
}

// Процесс подписки веб-хука
func (sf *kBusHttp) processSubscribe(req *msg_sub.SubscribeReq) *msg_sub.SubscribeResp {
	req.SelfCheck()
	handler := handler_http_sub.NewHandlerHttpSub(req.Topic_, req.WebHook_)
	resp := &msg_sub.SubscribeResp{
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

// Входящая публикация
func (sf *kBusHttp) postPublish(ctx *fiber.Ctx) error {
	sf.log.Debug("kBusHttp.postPublish()")
	ctx.Set("Content-type", "text/html; charset=utf8")
	ctx.Set("Content-type", "text/json")
	ctx.Set("Cache-Control", "no-cache")
	req := &msg_pub.PublishReq{}
	err := ctx.BodyParser(req)
	if err != nil {
		resp := &msg_pub.PublishResp{
			Status_: fmt.Sprintf("kernelBusHttp.postPublish(): in parse request, err=\n\t%v\n", err),
			Uuid_:   req.Uuid_,
		}
		resp.SelfCheck()
		ctx.Response().SetStatusCode(http.StatusBadRequest)
		sf.log.Err(resp.Status_)
		return ctx.JSON(resp)
	}
	resp := sf.processPublish(req)
	resp.SelfCheck()
	return ctx.JSON(resp)
}

// Выполняет процесс публикации
func (sf *kBusHttp) processPublish(req *msg_pub.PublishReq) *msg_pub.PublishResp {
	req.SelfCheck()
	err := sf.Publish(req.Topic_, req.BinMsg_)
	resp := &msg_pub.PublishResp{
		Status_: "ok",
		Uuid_:   req.Uuid_,
	}
	if err != nil {
		resp.Status_ = fmt.Sprintf("kernelBusHttp.processPublish(): err=\n\t%v", err)
	}
	return resp
}

// Входящий запрос
func (sf *kBusHttp) postSendRequest(ctx *fiber.Ctx) error {
	sf.log.Debug("kBusHttp.postSendRequest()")
	ctx.Set("Content-type", "text/html; charset=utf8")
	ctx.Set("Content-type", "text/json")
	ctx.Set("Cache-Control", "no-cache")
	req := &msg_serve.ServeReq{}
	err := ctx.BodyParser(req)
	if err != nil {
		resp := &msg_serve.ServeResp{
			Status_: fmt.Sprintf("kernelBusHttp.postSendRequest(): err=\n\t%v", err),
			Uuid_:   req.Uuid_,
		}
		resp.SelfCheck()
		ctx.Response().SetStatusCode(http.StatusBadRequest)
		sf.log.Err(resp.Status_)
		return ctx.JSON(resp)
	}
	resp := sf.processSendRequest(req)
	resp.SelfCheck()
	return ctx.JSON(resp)
}

// Обрабатывает входящий запрос
func (sf *kBusHttp) processSendRequest(req *msg_serve.ServeReq) *msg_serve.ServeResp {
	req.SelfCheck()
	binResp, err := sf.SendRequest(req.Topic_, req.BinReq_)
	resp := &msg_serve.ServeResp{
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
func (sf *kBusHttp) postUnsub(ctx *fiber.Ctx) error {
	sf.log.Debug("kBusHttp.postUnsub()")
	ctx.Set("Content-type", "text/html; charset=utf8")
	ctx.Set("Content-type", "text/json")
	ctx.Set("Cache-Control", "no-cache")
	req := &msg_unsub.UnsubReq{}
	err := ctx.BodyParser(req)
	if err != nil {
		resp := &msg_serve.ServeResp{
			Status_: fmt.Sprintf("kernelBusHttp.postSendRequest(): err=\n\t%v", err),
			Uuid_:   req.Uuid_,
		}
		resp.SelfCheck()
		ctx.Response().SetStatusCode(http.StatusBadRequest)
		sf.log.Err(resp.Status_)
		return ctx.JSON(resp)
	}
	resp := sf.processUnsubRequest(req)
	resp.SelfCheck()
	return ctx.JSON(resp)
}

// Процесс отписки от топика
func (sf *kBusHttp) processUnsubRequest(req *msg_unsub.UnsubReq) *msg_unsub.UnsubResp {
	req.SelfCheck()
	_hand := sf.Ctx_.Get(string(req.Name_))
	resp := &msg_unsub.UnsubResp{
		Status_: "ok",
		Uuid_:   req.Uuid_,
	}
	if _hand == nil {
		resp.Status_ = fmt.Sprintf("kernelBusHttp.processUnsubRequest(): handler name(%v) not exists", req.Name_)
		return resp
	}
	hand := _hand.Val().(IBusHandlerSubscribe)
	sf.Unsubscribe(hand)
	return resp
}
