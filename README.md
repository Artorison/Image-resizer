# Image-Resizer

Сервис предназначен для создания изображения с новыми размерами.

Сервис представляет собой web-сервер (прокси), загружающий изображения, масштабирующий/обрезающий их до нужного формата и возвращающий пользователю.

Реализована система кеша "Least Recent Used", при повторном запросе моментально отдается картинка из кэша

за размер кеша отвечает настройка `cache_size` в `configs/config.yaml`


### Быстрый старт

```bash
make up
make down
```
файл конфигурации приложения
`configs/config.yaml`,
image_quality - качество изображения от 1 до 100

по умолчанию приложение работает на `8080` порту

## Использование

```bash
GET localhost:8080/fill/<width>/<height>/<link>
```

### Исходное изображение:
![Screen 2](images/goOrigin.jpg)

### 150*300 resize
GET `http://localhost:8080/fill/150/300/https://example.jpg`

![Screen 2](images/150*300.jpeg)

### 100*100 resize
GET `http://localhost:8080/fill/100/100/https://example.jpg`

![Screen 2](images/100*100.jpeg)



## Тесты

#### Юнит тесты

```bash
make test
```


#### Интеграционные тесты с использованием nginx

```bash
make test_integration
```

Для интеграционных тестов используется отдельный docker compose файл
`test/docker-compose_test.yml`.

И свой конфиг файл `test/test_config.yaml`.
По умолчанию тестовое приложение использует 8089 порт.