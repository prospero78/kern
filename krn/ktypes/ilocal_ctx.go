package ktypes

import (
	. "github.com/prospero78/kern/krn/kalias"
)

// ICtxValue -- интерфейс к значению локального контекста
type ICtxValue interface {
	// Key -- возвращает ключ значения
	Key() string
	// Val -- возвращает хранимое значение
	Val() any
	// CreateAt -- возвращает метку времени создания
	CreateAt() ATime
	// UpdateAt -- возвращает метку времени обновления
	UpdateAt() ATime
	// Comment -- возвращает комментарий значения
	Comment() string
	// Update -- обновляет хранимое значение
	Update(val any, comment string)
}

// ILocalCtx -- локальный контекст
type ILocalCtx interface {
	// Get -- извлекает значение из контекста
	Get(key string) ICtxValue
	// Del -- удаляет значение из контекста
	Del(key string)
	// Set -- добавляет значение в контекст
	Set(key string, val any, comment string)
	// Size -- возвращает размер словаря контекста
	Size() int
	// SortedList -- возвращает сортированный список объектов контекста
	SortedList() []ICtxValue
	// Cancel -- отменяет контекст
	Cancel()
	// Done -- ожидает отмены контекста
	Done()
	// Log -- возвращает буфер сообщений
	Log() ILogBuf
}
