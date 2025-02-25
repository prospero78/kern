// package mock_env -- обеспечивает мок-окружение для тестов
package mock_env

import (
	"os"
	"strings"

	. "github.com/prospero78/kern/helpers"
)

// MockEnv -- объект мок-окружения
type MockEnv struct {
	dictEnv  map[string]string // Словарь переменных окружения
	pwd      string            // Текущий путь процесса
	fileName string            // Полный путь к .env файлу
	strEnv   string            // Фактическое содержимое .env файла
}

// MakeEnv -- возвращает новое мок-окружение
func MakeEnv() *MockEnv {
	sf := &MockEnv{
		dictEnv: map[string]string{},
	}
	sf.setPath()
	sf.readEnv()
	err := sf.parseEnv()
	Hassert(err == nil, "MakeEnv(): parseEnv, err=\n\t%v", err)
	err = sf.SetAll()
	Hassert(err == nil, "MakeEnv(): SetAll, err=\n\t%v", err)
	return sf
}

// Парсит полученный текст .env для установки переменных окружения
func (sf *MockEnv) parseEnv() error {
	lstStr := strings.Split(sf.strEnv, "\n")
	for _, val := range lstStr {
		pair := strings.TrimSuffix(val, "\r")
		if pair == "" {
			continue
		}
		if strings.HasPrefix(val, "#") {
			continue
		}
		err := sf.parseEnvStr(pair)
		Hassert(err == nil, "MockEnv.parseEnv(): in parse str(%v), err=\n\t%w", pair, err)
	}
	return nil
}

// Парсит строку окружения
func (sf *MockEnv) parseEnvStr(pair string) error {
	lstPair := strings.Split(pair, "=")
	Hassert(len(lstPair) >= 2, "MockEnv.parseEnvStr(): pair(%v) not have =", pair)
	key := lstPair[0]
	Hassert(key != "", "MockEnv.parseEnvStr(): in pair(%v) key is empty", pair)
	lstVal := lstPair[1:]

	val := ""
	for _, val_ := range lstVal {
		val += val_ + "="
	}
	val = strings.TrimSuffix(val, "=")

	Hassert(val != "", "MockEnv.parseEnvStr(): in pair(%v) val is empty", pair)
	sf.dictEnv[key] = val
	return nil
}

// Читает файл переменных окружения
func (sf *MockEnv) readEnv() {
	binEnv, err := os.ReadFile(sf.fileName)
	Hassert(err == nil, "MockEnv.readEnv(): in read %v, err=\n\t%v", sf.fileName, err)
	sf.strEnv = string(binEnv)
}

// Устанавливает текущий путь относительно проекта в dev
func (sf *MockEnv) setPath() {
	dir, err := os.Getwd()
	Hassert(err == nil, "MockEnv.setPath(): in get PWD, err=\n\t%v", err)
	lstDir := strings.Split(dir, "/")
	pwd := ""
	fileName := ""
	isFind := false
	for _, point := range lstDir {
		pwd += point + "/"
		fileName = pwd + ".env"
		_, err := os.ReadFile(fileName)
		if err == nil {
			isFind = true
			break
		}
	}
	Hassert(isFind, "MockEnv.setPath(): not find file .env")
	sf.pwd = pwd + "/bin_dev"
	err = os.MkdirAll(sf.pwd, 0750)
	Hassert(err == nil, "MockEnv.setPath(): in create PWD, err=\n\t%v", err)
	err = os.Chdir(sf.pwd)
	Hassert(err == nil, "MockEnv.setPath(): in set cwd(%v), err=\n\t%v", sf.pwd, err)
	sf.fileName = fileName
}

// Pwd -- возвращает текущий рабочий каталог процесса
func (sf *MockEnv) Pwd() string {
	return sf.pwd
}

// SetAll -- поднимает все переменные окружения
func (sf *MockEnv) SetAll() error {
	sf.ResetAll()
	for key, val := range sf.dictEnv {
		sf.Reset("key")
		err := os.Setenv(key, val)
		Hassert(err == nil, "MockEnv.SetAll(): in set env %v on key %v, err=\n\t%v", key, val, err)
	}
	return nil
}

// ResetAll -- сбрасывает все переменные окружения
func (sf *MockEnv) ResetAll() {
	for key := range sf.dictEnv {
		sf.Reset(key)
	}
}

// Reset -- сбрасывает конкретный ключ
func (sf *MockEnv) Reset(key string) {
	Hassert(key != "", "MockEnv.Reset(): key is empty")
	err := os.Unsetenv(key)
	Hassert(err == nil, "MockEnv.Reset(): in reset env %v, err=\n\t%v", key, err)
}
