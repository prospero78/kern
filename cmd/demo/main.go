// package main -- пускач для демонстратора монолита
package main

import (
	"github.com/prospero78/kern"
)

func main() {
	app := kern.NewMonolitLocal("Demo monolit")

	modServHttp := kern.NewModuleServHttp()
	app.Add(modServHttp)

	modKernelCtx := kern.NewModuleKernelCtx()
	app.Add(modKernelCtx)

	modKernKeep := kern.NewModuleKernelKeeper()
	app.Add(modKernKeep)

	app.Run()
	app.Wait()
}
