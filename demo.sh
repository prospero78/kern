# Предыдущую пустую строку НЕ УДАЛЯТЬ!!! НУЖНА ДЛЯ ТЕСТОВ!!!
# Переменная окружения [local, prod]
export STAGE=local

# URL для локального HTTP-сервера
export LOCAL_HTTP_URL="http://localhost:18200/"

# Путь для локального хранилища (нужен, если локальное хранилище используется)
export LOCAL_STORE_PATH=/store

cd ./bin_dev && \
./demo