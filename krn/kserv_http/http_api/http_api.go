// package http_api -- различные API для работы веб-морды
package http_api

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/prospero78/kern/krn/kctx"
)

// HttpApi -- различные API для работы веб-морды
type HttpApi struct{}

// NewHttpApi -- возвращает новое HttpApi
func NewHttpApi() *HttpApi {
	sf := &HttpApi{}
	kCtx := kctx.GetKernelCtx()
	fiberApp := kCtx.Get("fiberApp").Val().(*fiber.App)
	fiberApp.Post("/api/time", sf.postTime)
	return sf
}

// Возвращает текущее время сервера
func (sf *HttpApi) postTime(ctx *fiber.Ctx) error {
	strTime := time.Now().Local().Format("2006-01-02 15:04:05.000 -07 MST")
	return ctx.SendString(strTime)
}
