
# Инструкции по сборке idealcore

## Вариант 1: Pure-Go (рекомендуется для начала)

### Преимущества:
- ✅ Простая сборка (не нужен компилятор C++)
- ✅ Кроссплатформенность
- ✅ Быстрая компиляция
- ✅ Легче отладка

### Сборка:
```bash
# Просто собери без CGO
CGO_ENABLED=0 go build -o idealcore ./cmd/server/main.go

# Или с CGO=0 явно
go build -tags '!cgo' -o idealcore ./cmd/server/main.go

### Запуск:

```bash

./idealcore

Вариант 2: CGO + HNSW (для продакшена с большими данными)
Преимущества:
✅ Быстрый поиск (O(log n) vs O(n))
✅ Масштабируемость (сотни тысяч векторов)
✅ Меньше памяти на индекс
Требования:

```bash

# Ubuntu/Debian
sudo apt-get install build-essential

# Установи hnswlib
cd /tmp
git clone https://github.com/nmslib/hnswlib.git
cd hnswlib
sudo cp -r hnswlib /usr/local/include/

Сборка:

```bash

# С включенным CGO
CGO_ENABLED=1 go build -tags cgo -o idealcore ./cmd/server/main.go

# Или просто
go build -o idealcore ./cmd/server/main.go

Запуск:

```bash

./idealcore


Проверка какой вариант собран


```bash

# Pure-Go
file idealcore
# Вывод: ELF 64-bit LSB executable, x86-64

# CGO
ldd idealcore
# Вывод: покажет динамические библиотеки

Переключение между вариантами
Из go.mod (рекомендуется):

```bash

# Pure-Go
echo 'build: !cgo' > .build_flags

# CGO
echo 'build: cgo' > .build_flags

Через переменные окружения:

```bash

# Pure-Go
export CGO_ENABLED=0
go build -o idealcore ./cmd/server/main.go

# CGO
export CGO_ENABLED=1
go build -o idealcore ./cmd/server/main.go


Бенчмарки
Pure-Go:
1000 векторов: ~5ms поиск
10000 векторов: ~50ms поиск
100000 векторов: ~500ms поиск
CGO + HNSW:
1000 векторов: ~0.5ms поиск
10000 векторов: ~1ms поиск
100000 векторов: ~5ms поиск

Рекомендации

Сценарий	                            Вариант
Личное использование (<1000 записей)	    Pure-Go
Небольшой проект (<10000 записей)	    Pure-Go
Продакшен (>10000 записей)	            CGO + HNSW
Разработка/тестирование	                    Pure-Go
Деплой на сервер	                    CGO + HNSW

