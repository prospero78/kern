// package main -- пускач для демонстратора монолита
package main

import (
	"github.com/prospero78/kern"
)

func main() {
	app := kern.NewMonolitLocal("Demo monolit")
	modServHttp := kern.NewModuleServHttp()
	app.Add(modServHttp)
	app.Run()
	app.Wait()
}
