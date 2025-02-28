// package helpers -- содержит всякие полезняшки
//
// Пакет импортировать где нужно в нотации `. "gitlab.c2g.pw/back/uaj-abstract-client/pkg/helpers"`
package helpers

import (
	"fmt"
	"os"
	"time"

	. "github.com/prospero78/kern/krn/kalias"
)

var (
	IsStageLocal bool
	IsStageProd  bool
)

// Assert -- проверка на правильность утверждения с падением в панику на локальном стенде (soft assert)
func Assert(isCond bool, msgFormat string, args ...interface{}) {
	if isCond {
		return
	}
	msg := fmt.Sprintf("SOFT ASSERT "+msgFormat+"\n", args...)
	if IsStageLocal {
		panic(msg)
	}
	fmt.Print(msg)
}

// Hassert -- проверка на правильность утверждения с безусловным падением в панику (hard assert)
func Hassert(isCond bool, msgFormat string, args ...interface{}) {
	if isCond {
		return
	}
	msg := fmt.Sprintf("HARD ASSERT "+msgFormat+"\n", args...)
	panic(msg)
}

// TimeNow -- возвращает стандартную строку сейчас-времени "2006-01-02 15:04:05.000 -07 MST"
func TimeNow() ATime {
	strTime := time.Now().Local().Format("2006-01-02 15:04:05.000 -07 MST")
	return ATime(strTime)
}

func init_() {
	strStage := os.Getenv("STAGE")
	switch strStage {
	case "local":
		IsStageLocal = true
		IsStageProd = false
	case "prod":
		IsStageProd = true
		IsStageLocal = false
	case "":
		IsStageLocal = true
		IsStageProd = false
	default:
		panic(fmt.Sprintf("lepers.init_(): unknown env STAGE (%v)\n", strStage))
	}
}

func init() {
	init_()
}
