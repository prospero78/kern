package ktypes

import "github.com/gofiber/fiber/v2"

// IKernelServerHttp -- интерфейс к серверу HTTP в ядре
type IKernelServerHttp interface {
	// Run -- запускает сервер в работу (блокирующий вызов; надо для добавления роутов)
	Run()
	// Fiber -- возвращает объект веб-приложения fiber
	Fiber() *fiber.App
}
