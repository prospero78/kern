// package http_api -- различные API для работы веб-морды
package http_api

import (
	"github.com/gofiber/fiber/v2"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/krn/kctx"
)

// HttpApi -- различные API для работы веб-морды
type HttpApi struct{}

// NewHttpApi -- возвращает новое HttpApi
func NewHttpApi() *HttpApi {
	sf := &HttpApi{}
	kCtx := kctx.GetKernelCtx()
	fiberApp := kCtx.Get("fiberApp").Val().(*fiber.App)
	fiberApp.Post("/api_time", sf.postTime)
	return sf
}

// Возвращает текущее время сервера
func (sf *HttpApi) postTime(ctx *fiber.Ctx) error {
	strTime := TimeNowStr()
	return ctx.SendString(string(strTime))
}
