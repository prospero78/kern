package ktypes

// IModuleStat -- статистика модуля
type IModuleStat interface {
	// SvgSec -- возвращает SVG с посекундным графиком за последнюю минуту
	SvgSec() string
	// Add -- добавляет секундное показание в статистику
	Add(int)
}
