# kern

![Coverage](./docs/img/coverage.svg)
![MX Linux](https://img.shields.io/badge/-MX%20Linux-%23000000?style=for-the-badge&logo=MXlinux&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)
![Visual Studio Code](https://img.shields.io/badge/Visual%20Studio%20Code-0078d7.svg?style=for-the-badge&logo=visual-studio-code&logoColor=white)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![Stat](https://starchart.cc/prospero78/kern.svg)

**kern** -- модульные компоненты ядра с высокой надёжностью для любого микросервиса или модульного монолита.

## Контакты

Пишите запросы в соответствующем [разделе](https://github.com/prospero78/kern/issue?status=).

## Статус проекта

Готовность: `100%` .

Покрытие тестами: 100%

Линтеры: ошибок нет; цикломатическая сложность менее 11.

## Состав

Команда вывода дерева:

```bash
tree -I vendor -I bin_dev -d
```

```bash
.
├── kc                # Вспомогательные компоненты
│   ├── helpers       # Жёсткий и мягкий assert
│   ├── local_ctx     # Локальный контекст
│   │   └── ctx_value # Переменная контекста с метаинформацией
│   ├── log_buf       # Буферизованный лог
│   │   └── log_msg   # Сообщение буферизованного лога
│   └── safe_bool     # Потокобезопасная булева переменная
├── krn                          # Компоненты ядра
│   ├── kalias                   # Алиасы типов ядра
│   ├── kbus                     # Шина данных ядра
│   │   ├── dict_sub_hook        # Словарь обработчиков подписок
│   │   ├── dict_topic_serve     # Список топиков для обработчиков входящих запросов
│   │   ├── dict_topic_sub       # Список топиков подписки
│   │   ├── kbus_base            # БазоваЯ шина данных
│   │   ├── kbus_http            # Шина данных поверх HTTP
│   │   │   └── client_bus_http  # Клиент для шины данных поверх HTTP
│   │   ├── kbus_local           # Локальная шина данных
│   │   │   └── client_bus_local # Клиент для локальной шины данных
│   │   └── kbus_msg          # Сообщения для всех сетевых шин
│   │       ├── msg_pub       # Сообщения для публикации
│   │       ├── msg_serve     # Сообщения для запросов
│   │       ├── msg_sub       # Сообщения для подписки
│   │       └── msg_unsub     # Сообщения для отписки
│   ├── kctx              # Контекст ядра
│   │   ├── kernel_keeper # Сторож сигналов ОС
│   │   └── kwg           # Именованный ожидатель групп
│   ├── kmodule   # Компонент модуль ядра
│   ├── kmonolit  # Компонент модульного монолита ядра
│   ├── kserv_http      # Компонент встроенного быстрого HTTP-сервера (fiber)
│   │   └── static      # Встраиваемые статические файлы
│   │       ├── css     # Встраиваемые стили (bootstrap)
│   │       └── js      # Встраиваемые скрипты (htmx, hyperscript, bootstrap)
│   ├── kstore_kv   # Встраиваемое быстрое key:value хранилище (Badger)
│   └── ktypes  # Интерфейсы ядра
├── mds   # Типовые модули ядра
└── mock        # Мок-объекты для тестирования и экспериментов
    ├── mock_env            # Мок-окружение для запуска компонентов ядра
    ├── mock_hand_serve     # Мок-обработчик входящих запросов
    ├── mock_hand_sub_http  # Мок-обработчик подписки через HTTP-шину
    └── mock_hand_sub_local # Мок-обработчик подписки через локальную шину
```

## Версия компилятора

Не ниже `go 1.24.0`

## Лицензия

Код открытый, [лицензия MIT](./LICENSE.txt)

## Команды сборки

```bash
make      # ОБновление зависимостей
Make mod  # -//-
make test # Запуск тестов
make lint # Запуск линтеров
```
