# Сервис укорачивания ссылок

Сервис использует хэш функцию crc32 для того, чтобы укорачивать ссылки.     
Полученные короткие ссылки хранятся в БД(формат _короткая ссылка: полная ссылка_).    
Имеется возможность создавать кастомные ссылки. 

Добавлена валидация входных параметров (корректность параметров link и custom). Код снабжен     
тестами. 

###Запуск проекта

Для запуска проекта настроена docker среда. Для запуска достаточно прописать в папке проекта:

    docker-compose up

###Методы доступа к API сервиса:

Простое получение укороченной ссылки(параметр link):

    http://localhost:8181/?link=https://github.com/avito-tech/auto-backend-trainee-assignment
    
Получение кастомной укороченной ссылки(параметры link, custom):

    http://localhost:8181/?link=https://github.com/avito-tech/auto-backend-trainee-assignment&custom=GO
    
Переход по укороченной ссылке

    http://localhost:8181/h2ifrr
    http://localhost:8181/GO
    
### Нагрузочное тестирование

Нагрузочное тестирование проводилось с помощью утилиты apache benchmarks:

    Server Software:
    Server Hostname:        127.0.0.1
    Server Port:            8181
    
    Document Path:          /?link=https://vk.com/feed
    Document Length:        67 bytes
    
    Concurrency Level:      350
    Time taken for tests:   46.534 seconds
    Complete requests:      50000
    Failed requests:        3910
       (Connect: 0, Receive: 0, Length: 3910, Exceptions: 0)
    Non-2xx responses:      3910
    Total transferred:      9129620 bytes
    HTML transferred:       3252250 bytes
    Requests per second:    1074.49 [#/sec] (mean)
    Time per request:       325.737 [ms] (mean)
    Time per request:       0.931 [ms] (mean, across all concurrent requests)
    Transfer rate:          191.59 [Kbytes/sec] received
    
    Connection Times (ms)
                  min  mean[+/-sd] median   max
    Connect:        0    2   4.2      1      77
    Processing:     2  317 564.5    172   15850
    Waiting:        2  314 563.1    171   15849
    Total:          2  319 565.1    174   15851
    
    Percentage of the requests served within a certain time (ms)
      50%    174
      66%    250
      75%    313
      80%    361
      90%    835
      95%   1261
      98%   1563
      99%   3285
     100%  15851 (longest request)