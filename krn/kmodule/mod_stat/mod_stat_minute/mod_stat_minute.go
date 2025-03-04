// package mod_stat_minute -- статистика модуля за первые 60 мин
package mod_stat_minute

import (
	"bytes"
	"fmt"
	"math"

	"sync"
	"time"

	svg "github.com/ajstarks/svgo"
)

// ModStatMinutes -- статистика модуля за первые 60 мин
type ModStatMinutes struct {
	sync.Mutex
	lst    []int // Список значений за последние 60 мин
	bufSvg *bytes.Buffer
}

// NewModStatMinute -- возвращает новую статистику модуля за последние 60 мин
func NewModStatMinute() *ModStatMinutes {
	sf := &ModStatMinutes{
		lst:    []int{},
		bufSvg: bytes.NewBufferString(""),
	}
	return sf
}

// Add -- добавляет значение в часовой срез
func (sf *ModStatMinutes) Add(val int) {
	sf.Lock()
	defer sf.Unlock()

	for len(sf.lst) < 60 {
		sf.lst = append(sf.lst, math.MinInt64)
	}

	sf.lst = append(sf.lst, val)
	if len(sf.lst) > 60 {
		sf.lst = sf.lst[1:]
	}
}

// Svg -- возвращает сгенерированный SVG по часовому срезу
func (sf *ModStatMinutes) Svg() string {
	sf.Lock()
	defer sf.Unlock()
	sf.bufSvg.Reset()
	cnv := svg.New(sf.bufSvg)
	cnv.Start(480, 320)
	cnv.Title("Last 60 minute")
	cnv.Desc("Graphic of last 60 minute")
	cnv.Text(20, 20, "Last 60 minute", "")
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
		x1 := int(float32(i)*6) + 42
		y1 := int(240 * float32(valMax) / float32(val))
		cnv.Rect(x1, 280-y1, 5, y1, "fill:true;stroke:red;")
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
		timeSub60 := timeNow.Add(-60 * time.Minute).Format("04")
		cnv.Text(40, 295, timeSub60, "")

		timeSub45 := timeNow.Add(-45 * time.Minute).Format("04")
		cnv.Text(128, 295, timeSub45, "")

		timeSub30 := timeNow.Add(-30 * time.Minute).Format("04")
		cnv.Text(216, 295, timeSub30, "")

		timeSub15 := timeNow.Add(-15 * time.Minute).Format("04")
		cnv.Text(304, 295, timeSub15, "")

		timeSub0 := time.Now().Format("04")
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
