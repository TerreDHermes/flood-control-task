# Тестовое задание для стажера на позицию «Golang-разработчик»
Реализован интерфейс с методом для проверки правил флуд-контроля. Если за последние N секунд вызовов метода Check будет больше K, значит, проверка на флуд-контроль не пройдена.

База данных - PostgreSQL

Язык Go - 1.21

# Ход мыслей
Первая мысль - где хранить информацию о запросах? Кажется, что реляционная база данных слишком для этого тяжела и неудобна (из пушки по воробьям). Но 
в вакансии был один из пунктов - "Необходимо иметь опыт работы с реляционными базами данных (PostgreSQL, MySQL)". Поэтому решил взять PostgreSQL и попытаться как-то максимально оптимизировать работу.

Мы имеем временное окно (K), во время которого счиаем количество запросов. Поэтому, нужно как-то удалять данные, если они устарели. Чем больше данных, тем дольше мы будем отбирать нужные данные, которые укладываются во временное окно. 

Варианты для удаления:
1) Во время вставки будет срабатывать тригер;
   
2) Написать отдельное приложение, которое будет обрабатывать базу данных (чистить старые данные);
   
3) Создать горутину, которая будет удалять данные каждые N секунд. В принципе, это и есть 2 пункт, просто в все в одном.
 
4) Расширение для Postgres pg_cron. Можно выполнять любые действия над таблицей по расписанию.

Минус первого пункта - после каждой вставки происходит выборка,проверка,удаление... и то, если вообще за это время что-то устарело, возможно и нет. Слишком большая нагрузка на базу данных.

Минус 1, 2 и 3 пункта - нет общего консенсуса по времени жизни информации о запросе. Нужно придумывать методы взаиомедействия, чтобы придти к согласию.

Расширение pg_cron позволяет установить общие для всех правила. И выборка,проверка,удаление происходят раз M секунд. Нет постоянной нагрузки на базу данных. 

С данным расширением ни разу не сталкивался. Поэтому пришлось покопаться с документацией. Особые трудности возникли при сборке postres с pg_cron в Docker-контейнере.

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
Возможна конфигурация параметров, которые используются при создании контейнера, в котором располагается наша база данных Postgres. Базовая конфигурация в Makefile:

<div align="center">
  <img src="https://github.com/TerreDHermes/TerreDHermes/blob/main/assets/vk/makefile.png" alt="Описание изображения" style="width: 100%;">
</div>

Возможна конфигурация параметров, которые использует приложение при подключении к базе данных Postgres. Базовая конфигурация в config.yml:

<div align="center">
  <img src="https://github.com/TerreDHermes/TerreDHermes/blob/main/assets/vk/config.yml.png" alt="Описание изображения" style="width: 100%;">
</div>

В config.yml можно настроить параметры N и K, N - время в секундах, K - количество вызовов метода Check

Добавлено расширение pg_cron в Postgres, которое позволяет выполнять функцию по обработке данных в таблице согласно расписанию. Если честно, на пправильное подключение этого плагина ушло очень много времени.

<div align="center">
  <img src="https://github.com/TerreDHermes/TerreDHermes/blob/main/assets/vk/sql_up.png" alt="Описание изображения" style="width: 100%;">
</div>

Базовая конфигурация - каждые 25 секунд запускается функция, которая удаляет из таблицы данные, которые старше 20 секунд. То есть 20 секунд - время жизни кортежа в таблице. При условии, что показатель N < 20, зачем нам хранить данные, которые живут больше 20 секунд. Этот параметр (время жихни кортежа) можно менять в зависимости от N (если N=30, тогда время жизни кортежа меняем ,например, на 40).

При условии, когда много отдельных программ используют нашу базу, выставляется время жизни кортежа  равное большему N из всех N каждого участника. Участник - выявитель флуда со своими параметрами под свои задачи.

# Тестирование

Время было ограничено (сутки), поэтому тесты написать не успел. Слишком много времени потратил на изучение pg_cron.

Есть реализация вызова функции Check в цикле, где учитываются введенные параметры. Располагается этот цикл в main. Он запускатеся сразу после сборки Makefile-а. 

<div align="center">
  <img src="https://github.com/TerreDHermes/TerreDHermes/blob/main/assets/vk/тесты.png" alt="Описание изображения" style="width: 100%;">
</div>

Чтобы протестировать, придется в цикле for менять время ожидания следующего запроса в time.Sleep(1 * time.Second), а также аргумент в floodControl.Check(ctx, 1), который отвечает за id пользователя.
 




# Что нужно сделать

Реализовать интерфейс с методом для проверки правил флуд-контроля. Если за последние N секунд вызовов метода Check будет больше K, значит, проверка на флуд-контроль не пройдена.

- Интерфейс FloodControl располагается в файле main.go.

- Флуд-контроль может быть запущен на нескольких экземплярах приложения одновременно, поэтому нужно предусмотреть общее хранилище данных. Допустимо использовать любое на ваше усмотрение. 

# Необязательно, но было бы круто

Хорошо, если добавите поддержку конфигурации итоговой реализации. Параметры — на ваше усмотрение.
