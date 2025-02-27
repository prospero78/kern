// package kserv_http -- встроенный HTTP-сервер
package kserv_http

import (
	"embed"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/safe_bool"
	"github.com/prospero78/kern/krn/kctx"
	. "github.com/prospero78/kern/krn/ktypes"
)

const (
	streamName = "kernel_server_http" // Контрольная строка для ожидателя потока
)

// kServHttp -- встроенный HTTP-сервер
type kServHttp struct {
	ctx      IKernelCtx
	strPort  string // Порт ,на котором слушает HTTP-сервер
	fiberApp *fiber.App
	isWork   ISafeBool
	block    sync.Mutex
}

//go:embed static/*
var embedDirStatic embed.FS

var (
	kernServHttp *kServHttp
	block        sync.Mutex
)

// GetKernelServHttp -- возвращает  встроенный HTTP-сервер
func GetKernelServHttp() IKernelServerHttp {
	log.Println("GetKernelServHttp()")
	block.Lock()
	defer block.Unlock()
	ctx := kctx.GetKernelCtx()
	if kernServHttp != nil {
		return kernServHttp
	}
	strPort := os.Getenv("SERVER_HTTP_PORT")
	Hassert(strPort != "", "GetKernelServHttp(): env SERVER_HTTP_PORT not set")
	confFiber := fiber.Config{
		ServerHeader:      ctx.Get("monolitName").Val().(string),
		UnescapePath:      true,
		ReadTimeout:       time.Second * 15,
		WriteTimeout:      time.Second * 15,
		AppName:           ctx.Get("monolitName").Val().(string),
		Network:           "tcp4",
		EnablePrintRoutes: true,
	}
	sf := &kServHttp{
		ctx:      ctx,
		strPort:  strPort,
		fiberApp: fiber.New(confFiber),
		isWork:   safe_bool.NewSafeBool(),
	}
	sf.fiberApp.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression, // 2
	}))
	sf.fiberApp.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "static",
		Browse:     true,
		MaxAge:     3600 * 24,
	}))
	sf.fiberApp.Get("/monitor", monitor.New(monitor.Config{Title: "KernelHttpServer"}))
	err := sf.ctx.Wg().Add(streamName)
	Hassert(err == nil, "NewKernelServHttp(): in add stream %v, err=\n\t%v", streamName, err)
	ctx.Set("fiberApp", sf.fiberApp, "kServHttp: internal fiber app")
	kernServHttp = sf
	ctx.Set("kServHttp", kernServHttp, "kServHttp")
	return kernServHttp
}

// Fiber -- возвращает объект веб-приложения fiber
func (sf *kServHttp) Fiber() *fiber.App {
	return sf.fiberApp
}

// Run -- запускает сервер в работу (блокирующий вызов)
func (sf *kServHttp) Run() {
	go sf.close()
	sf.isWork.Set()
	err := sf.fiberApp.Listen(":" + sf.strPort)
	if err != nil {
		log.Printf("kServHttp.Run(): in listen, err=\n\t%v\n", err)
		sf.ctx.Cancel()
	}
}

// Ожидает окончания работы
func (sf *kServHttp) close() {
	sf.ctx.Done()
	sf.block.Lock()
	defer sf.block.Unlock()
	if !sf.isWork.Get() {
		return
	}
	sf.isWork.Reset()
	err := sf.fiberApp.Server().Shutdown()
	Assert(err == nil, "kServHttp.close(): in close server, err=\n\t%v", err)
	sf.ctx.Wg().Done(streamName)
	log.Println("kServHttp.close(): end")
}
