package helpers

import (
	"os"
	"testing"
)

type tester struct {
	t *testing.T
}

func TestHelpers(t *testing.T) {
	sf := &tester{
		t: t,
	}
	sf.assert()
	sf.hassert()
	sf.init_()
}

// Неизвестное значение STAGE
func (sf *tester) init_() {
	sf.t.Log("init_")
	sf.initBad1()
	_ = os.Unsetenv("STAGE")
	_ = os.Setenv("STAGE", "local")
	init_()
}

func (sf *tester) initBad1() {
	sf.t.Log("initBad1")
	defer func() {
		if panic_ := recover(); panic_ == nil {
			sf.t.Fatalf("initBad1(): panic==nil")
		}
	}()
	_ = os.Unsetenv("STAGE")
	_ = os.Setenv("STAGE", "tra-lala")
	init_()
}

// Проверка мягкого ассерта
func (sf *tester) hassert() {
	sf.t.Log("assert")
	sf.hassertLocal()
	sf.hassertProd()
	sf.hassertProdGood1()
	if strTime := TimeNow(); strTime == "" {
		sf.t.Fatalf("hassert(): strTime==''")
	}
	SleepMs()
}

// Мягкая ТВЁРДАЯ проверка на ок
func (sf *tester) hassertProdGood1() {
	sf.t.Log("hassertProdGood1")
	Hassert(true, "tra-la-la")
}

// Твёрдая проверка
func (sf *tester) hassertProd() {
	sf.t.Log("hassertProd")
	err := os.Setenv("STAGE", "prod")
	if err != nil {
		sf.t.Fatalf("hassertProd(): err=%v", err)
	}
	init_()
	defer func() {
		if panic_ := recover(); panic_ == nil {
			sf.t.Fatalf("assertLocal(): panic==nil")
		}
	}()
	Hassert(false, "tra-la-la")
}

// ТВЁРДАЯ проверка
func (sf *tester) hassertLocal() {
	sf.t.Log("hassertLocal")
	defer func() {
		if panic_ := recover(); panic_ == nil {
			sf.t.Fatalf("assertLocal(): panic==nil")
		}
	}()
	Hassert(false, "tra-la-la")
}

//========================================================================

// Проверка мягкого ассерта
func (sf *tester) assert() {
	sf.t.Log("assert")
	sf.assertLocal()
	sf.assertProd()
	sf.assertProdGood1()
}

// Мягкая ТВЁРДАЯ проверка на ок
func (sf *tester) assertProdGood1() {
	sf.t.Log("assertProdGood1")
	Assert(true, "tra-la-la")
}

// Мягкая мягка проверка (на проде)
func (sf *tester) assertProd() {
	sf.t.Log("assertProd")
	err := os.Setenv("STAGE", "prod")
	if err != nil {
		sf.t.Fatalf("assertProd(): err=%v", err)
	}
	init_()
	Assert(false, "tra-la-la")
}

// Мягкая ТВЁРДАЯ локальная проверка (локально)
func (sf *tester) assertLocal() {
	sf.t.Log("assertLocal")
	defer func() {
		if panic_ := recover(); panic_ == nil {
			sf.t.Fatalf("assertLocal(): panic==nil")
		}
	}()
	Assert(false, "tra-la-la")
}
