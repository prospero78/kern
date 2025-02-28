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
	return sf
}

const (
	strRowVal = `
<div class="row">
    <div class="col-2 border">
        {.key}
    </div>
    <div class="col-2 border">
        {.value}
    </div>
	<div class="col-1 border">
        {.type}
    </div>
    <div class="col-2 border">
        {.createAt}
    </div>
    <div class="col-2 border">
        {.updateAt}
    </div>
    <div class="col border">
        {.comment}
    </div>
</div>
`
)

// Возвращает блок контекста монолита
func (sf *PageMonolit) postMonolitCtx(ctx *fiber.Ctx) error {
	mon := sf.ctx.Get("monolit").Val().(IKernelMonolit)
	ctxMon := mon.Ctx().(*local_ctx.LocalCtx)
	dictVal := ctxMon.DictVal_
	strOut := `
<div class="container border rounded m-3 text-center">
    <h2>Monolit</h2>
</div>

<p></p>

<div class="container border rounded m-3 text-center">
<span class="btn btn-primary" hx-post="/monolit_state" hx-target="#main">Monolit</span>
    <span class="btn btn-primary" hx-post="/monolit_ctx" hx-target="#main">ctx</span>
    <span class="btn btn-primary" hx-post="/monolit_log" hx-target="#main">log</span>
</div>

<div class="container">

<div class="row">
    <div class="col-2 border bg-info">
        <b>Key</b>
    </div>
    <div class="col-2 border bg-info">
        <b>Value</b>
    </div>
	<div class="col-1 border bg-info">
        <b>Type</b>
    </div>
    <div class="col-2 border bg-info">
        <b>CreateAt</b>
    </div>
    <div class="col-2 border bg-info">
        <b>UpdateAt</b>
    </div>
    <div class="col border bg-info">
        <b>Comment</b>
    </div>
</div>`

	for key, val := range dictVal {
		strRow := strRowVal
		strRow = strings.ReplaceAll(strRow, "{.key}", key)
		strRow = strings.ReplaceAll(strRow, "{.value}", fmt.Sprint(val.Val()))
		strRow = strings.ReplaceAll(strRow, "{.type}", fmt.Sprintf("%#T", val.Val()))
		strRow = strings.ReplaceAll(strRow, "{.createAt}", string(val.CreateAt()))
		strRow = strings.ReplaceAll(strRow, "{.updateAt}", string(val.UpdateAt()))
		strRow = strings.ReplaceAll(strRow, "{.comment}", val.Comment())
		strOut += strRow
	}
	strOut += `
</div>

<span hx-post="/monolit_ctx" hx-trigger="every 2s" hx-target="#main"></span>`
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
