# Image-Resizer

### Запуск и остановка приложения

```bash
make up
make down
```
файл конфигурации приложения
configs/config.yaml
image_quality - качество изображения от 1 до 100

### Тесты

#### Запустить все проверки одной командой

```bash
make check
```

#### Юнит тесты

```bash
make test
```

#### Линтер

```bash
make lint
```

#### Интеграционные тесты

```bash
make test_integration
```

Для интеграционных тестов используется отдельный docker compose файл
`test/docker-compose_test.yml`.

И свой конфиг файл `test/test_config.yaml`.
По умолчанию тестовое приложение работает на 8089 порту.