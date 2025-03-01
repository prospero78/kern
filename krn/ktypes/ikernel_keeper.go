package ktypes

// IKernelKeeper -- сторож ядра
type IKernelKeeper interface {
	// Log -- возвращает лог сторожа
	Log() ILogBuf
}
