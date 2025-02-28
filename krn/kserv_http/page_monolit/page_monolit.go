// package page_monolit -- страница показа монолита
package page_monolit

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"

	"github.com/prospero78/kern/krn/kctx"
)

// PageMonolit -- страница показа монолита
type PageMonolit struct {
}

// NewPageMonolit -- возвращает новую страницу монолита
func NewPageMonolit() *PageMonolit {
	sf := &PageMonolit{}
	kCtx := kctx.GetKernelCtx()
	fiberApp := kCtx.Get("fiberApp").Val().(*fiber.App)
	fiberApp.Get("/monolit", sf.getMonolit)
	return sf
}

//go:embed page_monolit.html
var strPageMonolit string

// Индексная страница монолита
func (sf *PageMonolit) getMonolit(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/html; charset=utf8;\n\n")
	return ctx.SendString(strPageMonolit)
}
