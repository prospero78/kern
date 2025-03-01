// package ctx_value -- потокобезопасное значение локального контекста
package ctx_value

import (
	. "github.com/prospero78/kern/kc/helpers"
	. "github.com/prospero78/kern/krn/kalias"
	. "github.com/prospero78/kern/krn/ktypes"
)

type pairValue struct {
	val      any
	updateAt ATime
	comment  string
}

// ctxValue -- потокобезопасное значение локального контекста
type ctxValue struct {
	key      string
	createAt ATime

	chUpdateIn  chan pairValue
	chUpdateOut chan int

	chValIn  chan int
	chValOut chan any

	chUpdateAtIn  chan int
	chUpdateAtOut chan ATime

	chCommentIn  chan int
	chCommentOut chan string

	pair pairValue
}

// NewCtxValue -- возвращает новое потокобезопасное значение локального контекста
func NewCtxValue(key string, val any, comment string) ICtxValue {
	Hassert(key != "", "NewCtxValue(): key is empty")
	sf := &ctxValue{
		key:      key,
		createAt: TimeNow(),

		chUpdateIn:  make(chan pairValue, 2),
		chUpdateOut: make(chan int, 2),

		chValIn:  make(chan int, 2),
		chValOut: make(chan any),

		chUpdateAtIn:  make(chan int, 2),
		chUpdateAtOut: make(chan ATime, 2),

		chCommentIn:  make(chan int, 2),
		chCommentOut: make(chan string, 2),

		pair: pairValue{
			val:     val,
			comment: comment,
		},
	}
	go sf.run()
	return sf
}

// Update -- обновляет хранимое значение
func (sf *ctxValue) Update(val any, comment string) {
	pair := pairValue{
		val:      val,
		comment:  comment,
		updateAt: TimeNow(),
	}
	sf.chUpdateIn <- pair
	<-sf.chUpdateOut
}

// Key -- возвращает ключ значения
func (sf *ctxValue) Key() string {
	return sf.key
}

// Val -- возвращает хранимое значение
func (sf *ctxValue) Val() any {
	sf.chValIn <- 1
	return <-sf.chValOut
}

// UpdateAt -- возвращает время обновления значения
func (sf *ctxValue) UpdateAt() ATime {
	sf.chUpdateAtIn <- 1
	return <-sf.chUpdateAtOut
}

// CreateAt -- возвращает время создания значения
func (sf *ctxValue) CreateAt() ATime {
	return sf.createAt
}

// Comment -- возвращает комментарий значения
func (sf *ctxValue) Comment() string {
	sf.chCommentIn <- 1
	return <-sf.chCommentOut
}

// Работает в отдельном потоке
func (sf *ctxValue) run() {
	for {
		select {
		case pair := <-sf.chUpdateIn:
			sf.pair = pair
			sf.chUpdateOut <- 1
		case <-sf.chCommentIn:
			sf.chCommentOut <- sf.pair.comment
		case <-sf.chUpdateAtIn:
			sf.chUpdateAtOut <- sf.pair.updateAt
		case <-sf.chValIn:
			sf.chValOut <- sf.pair.val
		}
	}
}
