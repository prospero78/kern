// package main -- пускач для демонстратора монолита
package main

import (
	"gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

var app IKernelMonolit

func main() {
	app = kern.NewMonolitLocal("Demo monolit")

	modServHttp := kern.NewModuleServHttp()
	app.Add(modServHttp)

	modKernelCtx := kern.NewModuleKernelCtx()
	app.Add(modKernelCtx)

	modKernKeep := kern.NewModuleKernelKeeper()
	app.Add(modKernKeep)

	app.Run()
	app.Wait()
}
