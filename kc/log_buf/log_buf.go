// package log_buf -- потокобезопасный буфер лога
package log_buf

import (
	"fmt"

	"github.com/prospero78/kern/kc/log_buf/log_msg"
	. "github.com/prospero78/kern/krn/ktypes"
)

// logBuf -- потокобезопасный буфер лога
type logBuf struct {
	chGetIn  chan int
	chGetOut chan ILogMsg

	chDebugIn  chan tMsg
	chDebugOut chan int

	chInfoIn  chan tMsg
	chInfoOut chan int

	chWarnIn  chan tMsg
	chWarnOut chan int

	chErrorIn chan tMsg

	lst []ILogMsg

	chGetErrIn  chan int
	chGetErrOut chan ILogMsg
	lstErr      []ILogMsg

	chSizeIn  chan int
	chSizeOut chan int
}

// NewLogBuf -- возвращает новый потокобезопасный буфер лога
func NewLogBuf() ILogBuf {
	sf := &logBuf{
		chGetIn:  make(chan int, 2),
		chGetOut: make(chan ILogMsg, 2),

		chDebugIn:  make(chan tMsg, 2),
		chDebugOut: make(chan int, 2),

		chInfoIn:  make(chan tMsg, 2),
		chInfoOut: make(chan int, 2),

		chWarnIn:  make(chan tMsg, 2),
		chWarnOut: make(chan int, 2),

		chErrorIn: make(chan tMsg, 2),

		lst: []ILogMsg{},

		chGetErrIn:  make(chan int, 2),
		chGetErrOut: make(chan ILogMsg, 2),
		lstErr:      []ILogMsg{},

		chSizeIn:  make(chan int, 2),
		chSizeOut: make(chan int, 2),
	}
	go sf.run()
	return sf
}

// GetErr -- возвращает сообщение ошибки по номеру
func (sf *logBuf) GetErr(num int) ILogMsg {
	sf.chGetErrIn <- num
	return <-sf.chGetErrOut
}

// Get -- возвращает сообщение по номеру
func (sf *logBuf) Get(num int) ILogMsg {
	sf.chGetIn <- num
	return <-sf.chGetOut
}

type tMsg struct {
	text string
	args []any
}

// Debug -- сообщение отладки
func (sf *logBuf) Debug(fMsg string, args ...any) {
	msg := tMsg{
		text: fMsg,
		args: args,
	}
	sf.chDebugIn <- msg
	<-sf.chDebugOut
}

// Info -- информационные сообщения
func (sf *logBuf) Info(fMsg string, args ...any) {
	msg := tMsg{
		text: fMsg,
		args: args,
	}
	sf.chInfoIn <- msg
	<-sf.chInfoOut
}

// Warn -- предупреждающие сообщения
func (sf *logBuf) Warn(fMsg string, args ...any) {
	msg := tMsg{
		text: fMsg,
		args: args,
	}
	sf.chWarnIn <- msg
	<-sf.chWarnOut
}

// Err -- сообщения об ошибках
func (sf *logBuf) Err(fMsg string, args ...any) {
	msg := tMsg{
		text: fMsg,
		args: args,
	}
	sf.chErrorIn <- msg
}

// Size -- возвращает размер буфера
func (sf *logBuf) Size() int {
	sf.chSizeIn <- 1
	return <-sf.chSizeOut
}

func (sf *logBuf) run() {
	for {
		select {
		case num := <-sf.chGetErrIn:
			sf.chGetErrOut <- sf.getErr(num)
		case num := <-sf.chGetIn:
			sf.chGetOut <- sf.get(num)
		case msg := <-sf.chDebugIn:
			sf.debug(msg)
			sf.chDebugOut <- 1
		case msg := <-sf.chInfoIn:
			sf.info(msg)
			sf.chInfoOut <- 1
		case msg := <-sf.chWarnIn:
			sf.warn(msg)
			sf.chWarnOut <- 1
		case msg := <-sf.chErrorIn:
			sf.err(msg)
		case <-sf.chSizeIn:
			sf.chSizeOut <- len(sf.lst)
		}
	}
}

// Возвращает сообщение ошибки по номеру
func (sf *logBuf) getErr(num int) ILogMsg {
	if len(sf.lstErr) == 0 {
		return log_msg.NewLogMsg(log_msg.DEBUG, "not error msg")
	}
	if num >= len(sf.lstErr) {
		return sf.lstErr[len(sf.lstErr)-1]
	}
	if num <= 0 {
		return sf.lstErr[0]
	}
	return sf.lstErr[num]
}

// возвращает сообщение по номеру
func (sf *logBuf) get(num int) ILogMsg {
	if len(sf.lst) == 0 {
		return log_msg.NewLogMsg(log_msg.DEBUG, "*no msg*")
	}
	if num >= len(sf.lst) {
		return log_msg.NewLogMsg(log_msg.DEBUG, "*no msg*")
	}
	if num <= 0 {
		return log_msg.NewLogMsg(log_msg.DEBUG, "*no msg*")
	}
	return sf.lst[num]
}

// сообщение отладки
func (sf *logBuf) debug(msg tMsg) {
	strMsg := fmt.Sprintf(msg.text, msg.args...)
	_msg := log_msg.NewLogMsg(log_msg.DEBUG, strMsg)
	sf.lst = append(sf.lst, _msg)
	sf.checkLen()
}

// информационные сообщения
func (sf *logBuf) info(msg tMsg) {
	strMsg := fmt.Sprintf(msg.text, msg.args...)
	_msg := log_msg.NewLogMsg(log_msg.INFO, strMsg)
	sf.lst = append(sf.lst, _msg)
	sf.checkLen()
}

// предупреждающие сообщения
func (sf *logBuf) warn(msg tMsg) {
	strMsg := fmt.Sprintf(msg.text, msg.args...)
	_msg := log_msg.NewLogMsg(log_msg.WARN, strMsg)
	sf.lst = append(sf.lst, _msg)
	sf.checkLen()
}

// сообщения об ошибках
func (sf *logBuf) err(msg tMsg) {
	strMsg := fmt.Sprintf(msg.text, msg.args...)
	_msg := log_msg.NewLogMsg(log_msg.ERROR, strMsg)
	sf.lst = append(sf.lst, _msg)
	sf.lstErr = append(sf.lstErr, _msg)
	sf.checkLen()
	sf.checkLenErr()
}

// Проверяет длину общую лога
func (sf *logBuf) checkLen() {
	for len(sf.lst) > 100 {
		sf.lst = sf.lst[1:]
	}
}

// Проверяет длину лога ошибок
func (sf *logBuf) checkLenErr() {
	for len(sf.lstErr) > 100 {
		sf.lstErr = sf.lstErr[1:]
	}
}
