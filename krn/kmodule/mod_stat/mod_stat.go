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
	"github.com/prospero78/kern/krn/kmodule/mod_stat/mod_stat_day"
	"github.com/prospero78/kern/krn/kmodule/mod_stat/mod_stat_sec"
	. "github.com/prospero78/kern/krn/ktypes"
)

// ModStat -- статистика модуля
type ModStat struct {
	statSec    *mod_stat_sec.ModStatSec // Объект статистики 60 секунд
	timeMinute ISafeInt                 // Интервал ожидания минутного таймера, мсек
	statMin    *mod_stat_day.ModStatDay // Объект статистики 60 минут
	statDay    *mod_stat_day.ModStatDay // Объект статистики за последние 24 часа
	name       AModuleName
}

// NewModStat -- возвращает новую статистику модуля
func NewModStat(name AModuleName) *ModStat {
	Hassert(name != "", "NewModuleStat(): name module is empty")
	sf := &ModStat{
		statSec:    mod_stat_sec.NewModStatSec(),
		statMin:    mod_stat_day.NewModStatDay(),
		statDay:    mod_stat_day.NewModStatDay(),
		timeMinute: safe_int.NewSafeInt(),
		name:       name,
	}
	sf.timeMinute.Set(60 * 1000)
	go sf.eventMinute()
	return sf
}

// Срабатывает раз в минуту
func (sf *ModStat) eventMinute() {
	countPartHour := 20
	for {
		time.Sleep(time.Millisecond * time.Duration(sf.timeMinute.Get()))
		sum := sf.statSec.Sum()
		sf.statMin.Add(sum)
		countPartHour--
		if countPartHour == 0 {
			sum := sf.statMin.Sum()
			sf.statDay.Add(sum)
			countPartHour = 20
		}
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

// SvgDay -- возвращает SVG за последние сутки по часам
func (sf *ModStat) SvgDay() string {
	return sf.statDay.Svg()
}
