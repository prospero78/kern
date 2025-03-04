// package mod_stat -- статистика модуля
//
// Подробная статистика по периодам:
//   60 сек -- первая минута
//   60 минут -- первый час
//   48 получасов -- первые сутки
//   4 часа -- первые 14 суток

package mod_stat

import (
	"time"

	. "github.com/prospero78/kern/kc/helpers"
	"github.com/prospero78/kern/kc/safe_int"
	. "github.com/prospero78/kern/krn/kalias"
	"github.com/prospero78/kern/krn/kmodule/mod_stat/mod_stat_minute"
	"github.com/prospero78/kern/krn/kmodule/mod_stat/mod_stat_sec"
	. "github.com/prospero78/kern/krn/ktypes"
)

// ModStat -- статистика модуля
type ModStat struct {
	statSec    *mod_stat_sec.ModStatSec        // Объект статистики 60 секунд
	timeMinute ISafeInt                        // Интервал ожидания минутного таймера, мсек
	statMin    *mod_stat_minute.ModStatMinutes // Объект статистики 60 минут
	name       AModuleName
}

// NewModStat -- возвращает новую статистику модуля
func NewModStat(name AModuleName) *ModStat {
	Hassert(name != "", "NewModuleStat(): name module is empty")
	sf := &ModStat{
		statSec:    mod_stat_sec.NewModStatSec(),
		statMin:    mod_stat_minute.NewModStatMinute(),
		timeMinute: safe_int.NewSafeInt(),
		name:       name,
	}
	sf.timeMinute.Set(60 * 1000)
	go sf.eventMinute()
	return sf
}

// Срабатывает раз в минуту
func (sf *ModStat) eventMinute() {
	for {
		time.Sleep(time.Millisecond * time.Duration(sf.timeMinute.Get()))
		sum := sf.statSec.Sum()
		sf.statMin.Add(sum)
	}
}

// Add -- добавляет значение в статистику
func (sf *ModStat) Add(val int) {
	sf.statSec.Add(val)
}

// SvgSec -- возвращает посекундную SVG за последнюю минуту
func (sf *ModStat) SvgSec() string {
	return sf.statSec.Svg()
}

// SvgMin -- возвращает поминутную SVG за последнюю минуту
func (sf *ModStat) SvgMin() string {
	return sf.statMin.Svg()
}
