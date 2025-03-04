package ktypes

// IModuleStat -- статистика модуля
type IModuleStat interface {
	// SvgSec -- возвращает SVG с посекундным графиком за последнюю минуту
	SvgSec() string
	// SvgMin -- возвращает SVG с поминутным графиком за последний час
	SvgMin() string
	// Add -- добавляет секундное показание в статистику
	Add(int)
}
