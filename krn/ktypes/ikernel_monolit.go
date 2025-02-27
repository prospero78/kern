package ktypes

// IKernelMonolit -- интерфейс к монолиту на основе ядра
type IKernelMonolit interface {
	// IsLocal -- возвращает признак локальной шины
	IsLocal() bool
	// IsWork -- возвращает признак работы монолита
	IsWork() bool
	// Run -- запускает монолит в работу
	Run()
	// Add -- добавляет модуль в монолит
	Add(IKernelModule)
}
