package floodcontrol

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// FloodControlImpl - реализация интерфейса FloodControl
type FloodControlImpl struct {
	db         *sqlx.DB
	windowSize time.Duration // Окно времени для проверки флуда
	maxCalls   int           // Максимальное количество вызовов в окне времени
}

// NewFloodControl создает новый экземпляр FloodControlImpl
func NewFloodControl(db *sqlx.DB, windowSize time.Duration, maxCalls int) *FloodControlImpl {
	return &FloodControlImpl{
		db:         db,
		windowSize: windowSize,
		maxCalls:   maxCalls,
	}
}

// Check проверяет флуд для данного пользователя
func (fc *FloodControlImpl) Check(ctx context.Context, userID int64) (bool, error) {
	// Определение времени начала окна в формате UTC
	startTime := time.Now().UTC().Add(-fc.windowSize)

	// Вставляем новый вызов в таблицу с использованием времени UTC
	_, err := fc.db.ExecContext(ctx, `
	INSERT INTO flood_control (user_id, call_time)
	VALUES ($1, $2);
	`, userID, time.Now().UTC())
	if err != nil {
		return false, err
	}

	// Выполнение SQL-запроса для получения количества записей за последние N секунд в формате UTC
	var callCount int
	err = fc.db.QueryRowContext(ctx, `
	SELECT COUNT(*)
	FROM flood_control
	WHERE user_id = $1 AND call_time >= $2;
	`, userID, startTime).Scan(&callCount)
	if err != nil {
		return false, err
	}

	// Если количество вызовов превышает максимальное
	if callCount >= fc.maxCalls {
		return false, nil
	}

	return true, nil
}

func UpdateDeleteSQL(db *sqlx.DB, deleteInterval int, deletePeriod int) error {
	// Отмена текущего расписания
	_, err := db.Exec("SELECT cron.unschedule(jobid) FROM cron.job WHERE jobname = 'delete_old_records_schedule'")
	if err != nil {
		return fmt.Errorf("failed to cancel current schedule: %v", err)
	}

	// Создание нового расписания для запуска с вставленными переменными
	query := fmt.Sprintf("SELECT cron.schedule('delete_old_records_schedule', '%d seconds', 'DELETE FROM flood_control WHERE call_time < NOW() - INTERVAL ''%d seconds''')", deleteInterval, deletePeriod)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create new schedule: %v", err)
	}

	return nil
}
