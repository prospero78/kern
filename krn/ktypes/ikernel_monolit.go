package ktypes

// IKernelMonolit -- интерфейс к монолиту на основе ядра
type IKernelMonolit interface {
	// Name -- возвращает имя монолита
	Name() string
	// IsLocal -- возвращает признак локальной шины
	IsLocal() bool
	// IsWork -- возвращает признак работы монолита
	IsWork() bool
	// Run -- запускает монолит в работу
	Run()
	// Wait -- ожидание окончания работы
	Wait()
	// Add -- добавляет модуль в монолит
	Add(IKernelModule)
	// Log -- возвращает лог монолита
	Log() ILogBuf
	// Ctx -- возвращает контекст монолита
	Ctx() ILocalCtx
}
