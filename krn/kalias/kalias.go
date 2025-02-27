// package kalias -- алиасы типов ядра
package kalias

// ATime -- метка времени
//
// Форматная строка: time.Now().Local().String()
// Вывод: "2025-02-27 10:17:58.905806162 +0300 MSK")
type ATime string

// AStreamName -- имя потока для ожидания
type AStreamName string

// ATopic -- имя топика в шине
type ATopic string

// AHandlerName -- имя функции обработчика
type AHandlerName string

// AModuleName -- уникальное имя модуля
type AModuleName string
