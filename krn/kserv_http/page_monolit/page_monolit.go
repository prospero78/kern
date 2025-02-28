// package page_monolit -- страница показа монолита
package page_monolit

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/prospero78/kern/kc/local_ctx"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// PageMonolit -- страница показа монолита
type PageMonolit struct {
	ctx IKernelCtx
}

// NewPageMonolit -- возвращает новую страницу монолита
func NewPageMonolit() *PageMonolit {
	kCtx := kctx.GetKernelCtx()
	sf := &PageMonolit{
		ctx: kCtx,
	}
	fiberApp := kCtx.Get("fiberApp").Val().(*fiber.App)
	fiberApp.Get("/monolit", sf.getMonolit)
	fiberApp.Post("/monolit_state", sf.postMonolitState)
	fiberApp.Post("/monolit_ctx", sf.postMonolitCtx)
	fiberApp.Post("/monolit_log", sf.postMonolitLog)
	return sf
}

//go:embed log_block.html
var strLogBlock string

// Возвращает страницу лога монолита
func (sf *PageMonolit) postMonolitLog(ctx *fiber.Ctx) error {
	mon := sf.ctx.Get("monolit").Val().(IKernelMonolit)
	_log := mon.Log()
	strOut := ""
	for i := range 100 {
		msg := _log.Get(i).String()
		if strings.Contains(msg, "*no msg*") {
			continue
		}
		strOut += msg + "\n"
	}
	strOut = strings.ReplaceAll(strLogBlock, "{.log}", strOut)
	return ctx.SendString(strOut)
}

//go:embed ctx_row_val.html
var strCtxRowVal string

//go:embed ctx_row_block.html
var strCtxRowBlock string

// Возвращает блок контекста монолита
func (sf *PageMonolit) postMonolitCtx(ctx *fiber.Ctx) error {
	mon := sf.ctx.Get("monolit").Val().(IKernelMonolit)
	ctxMon := mon.Ctx().(*local_ctx.LocalCtx)
	dictVal := ctxMon.DictVal_
	strOut := ``
	for key, val := range dictVal {
		strRow := strCtxRowVal
		strRow = strings.ReplaceAll(strRow, "{.key}", key)
		strRow = strings.ReplaceAll(strRow, "{.value}", fmt.Sprint(val.Val()))
		strRow = strings.ReplaceAll(strRow, "{.type}", fmt.Sprintf("%#T", val.Val()))
		strRow = strings.ReplaceAll(strRow, "{.createAt}", string(val.CreateAt()))
		strRow = strings.ReplaceAll(strRow, "{.updateAt}", string(val.UpdateAt()))
		strRow = strings.ReplaceAll(strRow, "{.comment}", val.Comment())
		strOut += strRow
	}
	strOut = strings.ReplaceAll(strCtxRowBlock, "{.ctx_block}", strOut)
	return ctx.SendString(strOut)
}

//go:embed monolit_state.html
var strStateMonolit string

// Показывает состояние монолита
func (sf *PageMonolit) postMonolitState(ctx *fiber.Ctx) error {
	dictState := map[string]any{}
	mon := sf.ctx.Get("monolit").Val().(IKernelMonolit)
	dictState["{.name}"] = mon.Name()
	monVal := sf.ctx.Get("monolit")
	dictState["{.createAt}"] = monVal.CreateAt()
	dictState["{.updateAt}"] = monVal.UpdateAt()
	dictState["{.comment}"] = monVal.Comment()
	strOut := strStateMonolit
	for key, val := range dictState {
		strOut = strings.ReplaceAll(strOut, key, fmt.Sprint(val))
	}
	return ctx.SendString(strOut)
}

//go:embed page_monolit.html
var strPageMonolit string

// Индексная страница монолита
func (sf *PageMonolit) getMonolit(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/html; charset=utf8;\n\n")
	return ctx.SendString(strPageMonolit)
}
