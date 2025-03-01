// package page_modules -- страница представления модулей
package page_modules

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// PageModules -- отображает модули модулей
type PageModules struct {
	ctx IKernelCtx
}

var page *PageModules

// GetPageModules -- возвращает страницу модулей
func GetPageModules() *PageModules {
	if page != nil {
		return page
	}
	kCtx := kctx.GetKernelCtx()
	sf := &PageModules{
		ctx: kCtx,
	}
	fiberApp := kCtx.Get("fiberApp").Val().(*fiber.App)
	fiberApp.Post("/modules", sf.postModules)
	page = sf
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
	chLst := mon.Ctx().SortedList()
	strOut := ``
	for _, val := range chLst {
		if !strings.Contains(val.Key(), "module_") {
			continue
		}
		lstKey := strings.Split(val.Key(), "_")
		id := lstKey[1]
		strRow := strModRowBlank
		strRow = strings.ReplaceAll(strRow, "{.id}", id)
		strRow = strings.ReplaceAll(strRow, "{.key}", val.Key())
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
