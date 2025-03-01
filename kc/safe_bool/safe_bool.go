// package safe_bool -- потокобезопасный булевый признак
package safe_bool

import (
	. "github.com/prospero78/kern/krn/ktypes"
)

// safeBool -- потокобезопасный булевый признак
type safeBool struct {
	val      bool
	chGetIn  chan int
	chGetOut chan bool

	chSetIn  chan int
	chSetOut chan int

	chResetIn  chan int
	chResetOut chan int
}

// NewSafeBool -- возвращает новый потокобезопасный булевый признак
func NewSafeBool() ISafeBool {
	sf := &safeBool{
		chGetIn:  make(chan int, 2),
		chGetOut: make(chan bool, 2),

		chResetIn:  make(chan int, 2),
		chResetOut: make(chan int, 2),

		chSetIn:  make(chan int, 2),
		chSetOut: make(chan int, 2),
	}
	go sf.run()
	return sf
}

// Get -- возвращает хранимый булевый признак
func (sf *safeBool) Get() bool {
	sf.chGetIn <- 1
	return <-sf.chGetOut
}

// Set -- устанавливает булевый признак
func (sf *safeBool) Set() {
	sf.chSetIn <- 1
	<-sf.chSetOut
}

// Reset -- сбрасывает булевый признак
func (sf *safeBool) Reset() {
	sf.chResetIn <- 1
	<-sf.chResetOut
}

func (sf *safeBool) run() {
	for {
		select {
		case <-sf.chGetIn:
			sf.chGetOut <- sf.val
		case <-sf.chResetIn:
			sf.val = false
			sf.chResetOut <- 1
		case <-sf.chSetIn:
			sf.val = true
			sf.chSetOut <- 1
		}
	}
}
