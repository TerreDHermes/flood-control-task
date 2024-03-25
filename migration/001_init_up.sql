-- Создание таблицы flood_control
CREATE TABLE IF NOT EXISTS flood_control (
    user_id BIGINT NOT NULL,
    call_time TIMESTAMP NOT NULL
);

-- Установка расширения pg_cron
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Создание расписания для вызова функции delete_old_records() каждые 10 секунд
SELECT cron.schedule('delete_old_records_schedule', '25 seconds', 'DELETE FROM flood_control WHERE call_time < NOW() - INTERVAL ''20 seconds''');
