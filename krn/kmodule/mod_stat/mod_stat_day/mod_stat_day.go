// package mod_stat_day -- статистика модуля за первые 24 часа
package mod_stat_day

import (
	"bytes"
	"fmt"
	"math"

	"sync"
	"time"

	svg "github.com/ajstarks/svgo"
)

// ModStatВфн -- статистика модуля за первые 24 часа
type ModStatDay struct {
	sync.RWMutex
	lst    []int // Список значений за первые 24 часа
	bufSvg *bytes.Buffer
}

// NewModStatDay -- возвращает новую статистику модуля за первые 24 часа
func NewModStatDay() *ModStatDay {
	sf := &ModStatDay{
		lst:    []int{},
		bufSvg: bytes.NewBufferString(""),
	}
	return sf
}

// Sum -- возвращает сумму элементов по требованию
func (sf *ModStatDay) Sum() int {
	sf.RLock()
	defer sf.RUnlock()
	sum := 0
	for _, val := range sf.lst {
		sum += val
	}
	return sum
}

// Add -- добавляет значение в часовой срез
func (sf *ModStatDay) Add(val int) {
	sf.Lock()
	defer sf.Unlock()

	for len(sf.lst) < 72 {
		sf.lst = append(sf.lst, math.MinInt64)
	}

	sf.lst = append(sf.lst, val)
	if len(sf.lst) > 72 {
		sf.lst = sf.lst[1:]
	}
}

// Svg -- возвращает сгенерированный SVG по часовому срезу
func (sf *ModStatDay) Svg() string {
	sf.RLock()
	defer sf.RUnlock()
	sf.bufSvg.Reset()
	cnv := svg.New(sf.bufSvg)
	cnv.Start(480, 320)
	cnv.Title("Last 24 hours")
	cnv.Desc("Graphic of last 24 hours")
	cnv.Text(20, 20, "Last 24 hours", "")
	var (
		valMin = math.MaxInt64
		valMax = math.MinInt64
	)
	fnGetMinMax := func() { // Вычисляет максимальное и минимальное значение в графике
		for _, val := range sf.lst {
			if val < valMin {
				valMin = val
			}
			if val > valMax {
				valMax = val
			}
		}
	}
	fnGetMinMax()
	for i, val := range sf.lst {
		x1 := int(float32(i)*5) + 42
		y1 := int(240 * float32(valMax) / float32(val))
		cnv.Rect(x1, 280-y1, 4, y1, "fill:true;stroke:red;")
	}
	fnDrawNet := func() {
		// Метки величины
		if valMin == math.MaxInt64 || valMin == math.MinInt64 {
			valMin = 0
		}
		cnv.Text(25, 285, fmt.Sprint(valMin), "")
		if valMax == math.MinInt64 || valMax == math.MaxInt64 {
			valMax = 1
		}
		cnv.Text(25, 45, fmt.Sprint(valMax), "")
		// Метки времени
		timeNow := time.Now().Local()
		timeSub24 := timeNow.Add(-24 * time.Hour).Format("15")
		cnv.Text(40, 295, timeSub24, "")

		timeSub45 := timeNow.Add(-18 * time.Hour).Local().Format("15")
		cnv.Text(128, 295, timeSub45, "")

		timeSub30 := timeNow.Add(-12 * time.Hour).Local().Format("15")
		cnv.Text(216, 295, timeSub30, "")

		timeSub15 := timeNow.Add(-6 * time.Hour).Local().Format("15")
		cnv.Text(304, 295, timeSub15, "")

		timeSub0 := time.Now().Local().Format("15")
		cnv.Text(392, 295, timeSub0, "")

		cnv.Line(40, 280, 40, 38, "fill:true;stroke:black;stroke-width:4")
		cnv.Line(40, 280, 442, 280, "fill:true;stroke:black;stroke-width:4")
		count := 0
		for x := 40; x < 460; x += 20 {
			cnv.Line(x, 40, x, 280, "fill:true;stroke:gray")
			if count%5 == 0 {
				cnv.Text(x+3, 278, fmt.Sprint(x), "fill:white;stroke:gray")
			}
			count++
		}
		count = 0
		for y := 40; y < 300; y += 20 {
			cnv.Line(40, y, 440, y, "fill:true;stroke:gray")
			if count%5 == 0 {
				cnv.Text(43, y+12, fmt.Sprint(y), "fill:white;stroke:gray")
			}
			count++
		}
	}
	fnDrawNet()
	cnv.End()
	strOut := sf.bufSvg.String()
	return strOut
}
