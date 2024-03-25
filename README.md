# Тестовое задание для стажера на позицию «Golang-разработчик»
Реализован интерфейс с методом для проверки правил флуд-контроля. Если за последние N секунд вызовов метода Check будет больше K, значит, проверка на флуд-контроль не пройдена.

База данных - PostgreSQL

Язык Go - 1.21

# Запуск
Запуск Makefile из корневой папки:

```bash
make
```

Если будет Ping-ошибка, значит база данных не успела полностью запуститься. Подождем пару секунд и введём команду для запуска main.go

```bash
go run cmd/main.go
```

# Возможности в конфигурации

<div align="center">
  <img src="https://github.com/TerreDHermes/TerreDHermes/blob/main/assets/vk/makefile.png" alt="Описание изображения" style="width: 70%;">
</div>



# Что нужно сделать

Реализовать интерфейс с методом для проверки правил флуд-контроля. Если за последние N секунд вызовов метода Check будет больше K, значит, проверка на флуд-контроль не пройдена.

- Интерфейс FloodControl располагается в файле main.go.

- Флуд-контроль может быть запущен на нескольких экземплярах приложения одновременно, поэтому нужно предусмотреть общее хранилище данных. Допустимо использовать любое на ваше усмотрение. 

# Необязательно, но было бы круто

Хорошо, если добавите поддержку конфигурации итоговой реализации. Параметры — на ваше усмотрение.
