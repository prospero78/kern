// package log_msg -- сообщение логгера
package log_msg

import (
	"fmt"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

const (
	DEBUG = -3
	INFO  = -2
	WARN  = -1
	ERROR = 0
)

// logMsg -- сообщение логгера
type logMsg struct {
	level    string
	createAt ATime
	msg      string
}

// NewLogMsg -- возвращает новое сообщение логгера
func NewLogMsg(level int, msg string) ILogMsg {
	sf := &logMsg{
		createAt: TimeNowStr(),
		msg:      msg,
	}
	sf.check(level)
	return sf
}

// String -- возвращает форматированное сообщение лога
func (sf *logMsg) String() string {
	strOut := fmt.Sprintf("%v   %v  %v", sf.level, sf.createAt, sf.msg)
	return strOut
}

// Msg -- возвращает хранимое сообщение
func (sf *logMsg) Msg() string {
	return sf.msg
}

// Level -- возвращает уровень сообщения
func (sf *logMsg) Level() string {
	return sf.level
}

// CreateAt -- когда сообщение создано
func (sf *logMsg) CreateAt() ATime {
	return sf.createAt
}

// Проверяет правильность своего состава
func (sf *logMsg) check(level int) {
	switch level {
	case DEBUG:
		sf.level = "DEBU"
	case INFO:
		sf.level = "INFO"
	case WARN:
		sf.level = "WARN"
	case ERROR:
		sf.level = "ERRO"
	default:
		Hassert(false, "logMsg.check(): unknown level(%v)", level)
	}
}
