# Image-Resizer

Сервис предназначен для создания изображения с новыми размерами.

Сервис представляет собой web-сервер (прокси), загружающий изображения, масштабирующий/обрезающий их до нужного формата и возвращающий пользователю.

Реализовано проксирование заголовков (исходные HTTP-заголовки передаются на удалённый сервер) и **LRU-кэш**, за размер которого отвечает настройка `cache_size` в `configs/config.yaml`



### Быстрый старт

```bash
make up
# или
docker compose up
```

файл конфигурации приложения - 
`configs/config.yaml`,
image_quality - качество изображения от 1 до 100. 
По умолчанию приложение работает на `8080` порту

## Использование

```bash
GET /fill/{width}/{height}/{sourceURL}
```
### Исходное изображение:
![Screen 2](images/goOrigin.jpg)

### 150x300 resize
GET `http://localhost:8080/fill/150/300/https://example.jpg`

![Screen 2](images/150x300.jpeg)

### 100x100 resize
GET `http://localhost:8080/fill/100/100/https://example.jpg`

![Screen 2](images/100x100.jpeg)



## Тесты

#### Юнит тесты

```bash
make test
```


#### Интеграционные тесты с использованием Nginx

```bash
make test_integration
```

Для интеграционных тестов используется отдельный docker compose файл
`test/docker-compose_test.yml` и отдельный файл конфигурации - `test/test_config.yaml`.