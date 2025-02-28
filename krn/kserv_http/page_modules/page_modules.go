// package page_modules -- страница представления модулей
package page_modules

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/prospero78/kern/kc/local_ctx"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// PageModules -- отображает модули модулей
type PageModules struct {
	ctx IKernelCtx
}

// NewPageModules -- возвращает новую страницу модулей
func NewPageModules() *PageModules {
	kCtx := kctx.GetKernelCtx()
	sf := &PageModules{
		ctx: kCtx,
	}
	fiberApp := kCtx.Get("fiberApp").Val().(*fiber.App)
	fiberApp.Post("/modules", sf.postModules)
	return sf
}

//go:embed mod_row_block.html
var strModRowBlock string

//go:embed mod_row_val.html
var strModRowBlank string

// Индексная страница модулей
func (sf *PageModules) postModules(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/html; charset=utf8;\n\n")
	mon := sf.ctx.Get("monolit").Val().(IKernelMonolit)
	ctxMon := mon.Ctx().(*local_ctx.LocalCtx)
	dictVal := ctxMon.DictVal_
	strOut := ``
	for key, val := range dictVal {
		if !strings.Contains(key, "module/") {
			continue
		}
		strRow := strModRowBlank
		strRow = strings.ReplaceAll(strRow, "{.key}", key)
		moduleName := string(val.Val().(IKernelModule).Name())
		strRow = strings.ReplaceAll(strRow, "{.name}", moduleName)
		type_ := fmt.Sprintf("%#T", val.Val())
		type_ = strings.ReplaceAll(type_, ".", ".<br>")
		strRow = strings.ReplaceAll(strRow, "{.type}", type_)
		strRow = strings.ReplaceAll(strRow, "{.createAt}", string(val.CreateAt()))
		strRow = strings.ReplaceAll(strRow, "{.updateAt}", string(val.UpdateAt()))
		strRow = strings.ReplaceAll(strRow, "{.comment}", val.Comment())
		strOut += strRow
	}
	strOut = strings.ReplaceAll(strModRowBlock, "{.mod_block}", strOut)
	return ctx.SendString(strOut)
}
