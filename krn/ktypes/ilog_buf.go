package ktypes

import (
	. "github.com/prospero78/kern/krn/kalias"
)

// ILogMsg -- сообщение лога
type ILogMsg interface {
	// CreateAt -- когда создано
	CreateAt() ATime
	// Level -- уровень сообщения
	Level() string
	// Msg -- текст сообщения
	Msg() string
	// String -- форматированная строка
	String() string
}

// ILogBuf -- буферизованный лог для диагностики
//
//	Буфер для Error -- отдельный
type ILogBuf interface {
	// Debug -- сообщение отладки
	Debug(fMsg string, args ...any)
	// Info -- информационные сообщения
	Info(fMsg string, args ...any)
	// Warn -- предупреждающие сообщения
	Warn(fMsg string, args ...any)
	// Err -- сообщения об ошибках
	Err(fMsg string, args ...any)
	// Get -- возвращает сообщение по номеру (0..99)
	Get(num int) ILogMsg
	// GetErr -- возвращает сообщение ошибки по номеру (0..99)
	GetErr(num int) ILogMsg
}
