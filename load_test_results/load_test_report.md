# Отчет по нагрузочному тестированию и балансировке нагрузки

**Дата тестирования:** Чт 20 ноя 2025 17:18:36 MSK  
**Инструмент:** ApacheBench 2.3  
**Конечная точка:** /health/

## 1. Методология тестирования

### Конфигурация системы
- **Балансировщик:** Nginx 1.24.0
- **Бэкенды:** 3 инстанса Go приложения
- **Схема балансировки:** Weighted Round Robin (2:1:1)
- **Upstream конфигурация:**
```nginx
upstream backend_all {
    server api:8080 weight=2;
    server api-readonly-1:8081 weight=1;
    server api-readonly-2:8082 weight=1;
}

upstream backend_primary {
    server api:8080;
}

```

## 2. Результаты нагрузочного тестирования

### 2.1 Тест 1: Базовая нагрузка (2000 запросов, 10 параллельных соединений)
```
This is ApacheBench, Version 2.3 <$Revision: 1903618 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


Server Software:        nginx/1.24.0
Server Hostname:        localhost
Server Port:            80

Document Path:          /health/
Document Length:        42 bytes

Concurrency Level:      10
Time taken for tests:   0.637 seconds
Complete requests:      2000
Failed requests:        0
Non-2xx responses:      2000
Requests per second:    3140.95 [#/sec] (mean)
Time per request:       3.184 [ms] (mean)
Time per request:       0.318 [ms] (mean, across all concurrent requests)
Transfer rate:          1049.03 [Kbytes/sec] received
```

### 2.2 Тест 2: Средняя нагрузка (3000 запросов, 20 параллельных соединений)
```
This is ApacheBench, Version 2.3 <$Revision: 1903618 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 300 requests
Completed 600 requests
Completed 900 requests
Completed 1200 requests
Completed 1500 requests
Completed 1800 requests
Completed 2100 requests
Completed 2400 requests
Completed 2700 requests
Completed 3000 requests
Finished 3000 requests


Server Software:        nginx/1.24.0
Server Hostname:        localhost
Server Port:            80

Document Path:          /health/
Document Length:        42 bytes

Concurrency Level:      20
Time taken for tests:   0.838 seconds
Complete requests:      3000
Failed requests:        0
Non-2xx responses:      3000
Requests per second:    3578.00 [#/sec] (mean)
Time per request:       5.590 [ms] (mean)
Time per request:       0.279 [ms] (mean, across all concurrent requests)
Transfer rate:          1195.00 [Kbytes/sec] received
```

### 2.3 Тест 3: Высокая нагрузка (5000 запросов, 30 параллельных соединений)
```
This is ApacheBench, Version 2.3 <$Revision: 1903618 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 500 requests
Completed 1000 requests
Completed 1500 requests
Completed 2000 requests
Completed 2500 requests
Completed 3000 requests
Completed 3500 requests
Completed 4000 requests
Completed 4500 requests
Completed 5000 requests
Finished 5000 requests


Server Software:        nginx/1.24.0
Server Hostname:        localhost
Server Port:            80

Document Path:          /health/
Document Length:        42 bytes

Concurrency Level:      30
Time taken for tests:   1.316 seconds
Complete requests:      5000
Failed requests:        0
Non-2xx responses:      5000
Requests per second:    3798.29 [#/sec] (mean)
Time per request:       7.898 [ms] (mean)
Time per request:       0.263 [ms] (mean, across all concurrent requests)
Transfer rate:          1268.57 [Kbytes/sec] received
```

## 3. Сводка производительности

| Уровень нагрузки | Параллельных соединений | Запросов/сек | Время ответа (мс) | Неудачных запросов |
|------------------|------------------------|--------------|-------------------|-------------------|
| Базовая | 10 | 3140.95 | 3.184 | 0 |
| Средняя | 20 | 3578.00 | 5.590 | 0 |
| Высокая | 30 | 3798.29 | 7.898 | 0 |

## 5. Доказательства работы балансировки
Был написан скрипт для тестирования балансировки, который отправляет 100 запросов с помощью ```curl```, затем анализируется сколько запросов попало на какой порт.
```
=== Results ===
Backend 8080 (weight 2): 50 requests
Backend 8081 (weight 1): 25 requests
Backend 8082 (weight 1): 25 requests
Unknown: 0 requests

Expected ratio: 2:1:1 (50%:25%:25%)
Actual ratio: 50% : 25% : 25%
```