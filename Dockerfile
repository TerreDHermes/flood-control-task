# Используем официальный образ PostgreSQL с Docker Hub
FROM postgres:latest

# Установка pg_cron
RUN apt-get update && \
    apt-get install -y postgresql-server-dev-16 build-essential git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Установка pg_cron через исходный код
RUN git clone https://github.com/citusdata/pg_cron.git /pg_cron && \
    cd /pg_cron && \
    make && \
    make install

# Копирование SQL-скрипта инициализации базы данных
COPY migration/001_init_up.sql /docker-entrypoint-initdb.d/

# Открытие порта по умолчанию
EXPOSE 5432

# Добавление настроек pg_cron в postgresql.conf
RUN echo "shared_preload_libraries = 'pg_cron'" >> /usr/share/postgresql/postgresql.conf.sample
RUN echo "cron.database_name = 'postgres'" >> /usr/share/postgresql/postgresql.conf.sample
RUN echo "cron.timezone = 'PRC'" >> /usr/share/postgresql/postgresql.conf.sample
