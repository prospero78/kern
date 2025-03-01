// package page_module -- страница показа модуля
package page_module

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

// PageModule -- страница показа модуля
type PageModule struct {
	ctx IKernelCtx
}

var page *PageModule

// GetPageModule -- возвращает страницу модуля
func GetPageModule() *PageModule {
	if page != nil {
		return page
	}
	kCtx := kctx.GetKernelCtx()
	sf := &PageModule{
		ctx: kCtx,
	}
	fiberApp := kCtx.Get("fiberApp").Val().(*fiber.App)
	fiberApp.Post("/module/:id", sf.postModuleState)
	fiberApp.Post("/module_ctx/:id", sf.postModuleCtx)
	fiberApp.Post("/module_log/:id", sf.postModuleLog)
	page = sf
	return sf
}

//go:embed log_block.html
var strLogBlock string

// Возвращает страницу лога модуля
func (sf *PageModule) postModuleLog(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "1")
	module, _ := sf.getModule(id)
	if module == nil {
		strOut := strings.ReplaceAll(strLogBlock, "{.id}", id)
		strOut = strings.ReplaceAll(strOut, "{.log}", strOut)
		strOut = strings.ReplaceAll(strOut, "{.name}", "not found")
		return ctx.SendString(strOut)
	}
	_log := module.Log()
	strOut := ""
	for i := range 100 {
		msg := _log.Get(i).String()
		if strings.Contains(msg, "*no msg*") {
			continue
		}
		strOut += msg + "\n"
	}
	strOut = strings.ReplaceAll(strLogBlock, "{.log}", strOut)
	strOut = strings.ReplaceAll(strOut, "{.name}", string(module.Name()))
	strOut = strings.ReplaceAll(strOut, "{.id}", id)
	return ctx.SendString(strOut)
}

//go:embed ctx_row_val.html
var strCtxRowVal string

//go:embed ctx_row_block.html
var strCtxRowBlock string

// Возвращает блок контекста монолита
func (sf *PageModule) postModuleCtx(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "1")
	module, _ := sf.getModule(id)
	if module == nil {
		strOut := strings.ReplaceAll(strCtxRowBlock, "{.id}", id)
		strOut = strings.ReplaceAll(strOut, "{.name}", "not found")
		strOut = strings.ReplaceAll(strOut, "{.ctx_block}", "")
		return ctx.SendString(strOut)
	}
	lst := module.Ctx().SortedList()
	strOut := ""
	for _, val := range lst {
		strRow := strCtxRowVal
		strRow = strings.ReplaceAll(strRow, "{.key}", val.Key())
		strRow = strings.ReplaceAll(strRow, "{.value}", fmt.Sprint(val.Val()))
		strRow = strings.ReplaceAll(strRow, "{.type}", fmt.Sprintf("%#T", val.Val()))
		strRow = strings.ReplaceAll(strRow, "{.createAt}", string(val.CreateAt()))
		strRow = strings.ReplaceAll(strRow, "{.updateAt}", string(val.UpdateAt()))
		strRow = strings.ReplaceAll(strRow, "{.comment}", val.Comment())
		strOut += strRow
	}
	strOut = strings.ReplaceAll(strCtxRowBlock, "{.ctx_block}", strOut)
	strOut = strings.ReplaceAll(strOut, "{.id}", id)
	strOut = strings.ReplaceAll(strOut, "{.name}", string(module.Name()))
	return ctx.SendString(strOut)
}

//go:embed module_state.html
var strStateModule string

// Показывает состояние модуля
func (sf *PageModule) postModuleState(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "1")
	module, modVal := sf.getModule(id)
	if module == nil {
		return ctx.SendString(strStateModule)
	}
	dictState := map[string]any{}
	dictState["{.name}"] = module.Name()
	dictState["{.createAt}"] = modVal.CreateAt()
	dictState["{.updateAt}"] = modVal.UpdateAt()
	dictState["{.comment}"] = modVal.Comment()
	dictState["{.id}"] = id
	strOut := strStateModule
	for key, val := range dictState {
		strOut = strings.ReplaceAll(strOut, key, fmt.Sprint(val))
	}
	return ctx.SendString(strOut)
}

// Возвращает модуль
func (sf *PageModule) getModule(id string) (IKernelModule, ICtxValue) {
	mon := sf.ctx.Get("monolit").Val().(IKernelMonolit)
	lst := mon.Ctx().SortedList()
	var (
		val    ICtxValue
		isFind bool
	)
	for _, val = range lst {
		name := "module_" + id
		if name == val.Key() {
			isFind = true
			break
		}
	}
	if !isFind {
		return nil, nil
	}
	mod := val.Val().(IKernelModule)
	return mod, val
}
