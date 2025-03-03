// package mod_stat -- статистика модуля
//
// Подробная статистика по периодам:
//   60 сек -- первая минута
//   60 минут -- первый час
//   48 получасов -- первые сутки
//   4 часа -- первые 14 суток

package mod_stat

import (
	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kmodule/mod_stat/mod_stat_sec"
)

// ModStat -- статистика модуля
type ModStat struct {
	statSec *mod_stat_sec.ModStatSec // Объект статистики 60 секунд
	name    AModuleName
}

// NewModStat -- возвращает новую статистику модуля
func NewModStat(name AModuleName) *ModStat {
	Hassert(name != "", "NewModuleStat(): name module is empty")
	sf := &ModStat{
		statSec: mod_stat_sec.NewModStatSec(),
		name:    name,
	}
	return sf
}

// Add -- добавляет значение в статистику
func (sf *ModStat) Add(val int) {
	sf.statSec.Add(val)
}

// SvgSec -- возвращает посекундную SVG за последнюю минуту
func (sf *ModStat) SvgSec() string {
	return sf.statSec.Svg()
}
