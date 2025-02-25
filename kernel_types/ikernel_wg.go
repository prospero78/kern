package kernel_types

import (
	. "github.com/prospero78/kern/kernel_alias"
)

// IKernelWg -- интерфейс к именованному ожидателю потоков
type IKernelWg interface {
	// Add -- добавляет имя потока в ожидатель потоков
	Add(AStreamName) error
	// Done -- удаляет имя потока из ожидания
	Done(AStreamName)
	// Wait -- ожидает завершения работы всех потоков
	Wait()
	// IsWork -- признак работы ожидателя потоков (после закрытия добавлять нельзя)
	IsWork() bool
	// List -- возвращает список имён потоков на ожидании
	List() []AStreamName
	// Len -- возвращает размер списка потоков ожидания
	Len() int
}
