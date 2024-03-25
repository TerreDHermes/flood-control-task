package main

import (
	"context"
	"log"
	"os"
	"time"

	dB "vk/internal/db"
	fl "vk/internal/floodcontrol"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Чтение конфигурационного файла
	configFile, err := os.Open("configs/config.yml")
	if err != nil {
		logrus.Fatalf("error open configs/config.yml: %s", err.Error())
	}
	defer configFile.Close()

	config := fl.Config{}
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		logrus.Fatalf("failed to decode config file: %s", err.Error())
	}

	// подключение к бд
	db, err := dB.NewPosrgresDB(config)

	if config.Delete.IntervalDelete != 0 && config.Delete.PeriodDelete != 0 {
		if err = fl.UpdateDeleteSQL(db, config.Delete.IntervalDelete, config.Delete.PeriodDelete); err != nil {
			logrus.Errorf("failed to update delete mode (cron): %s", err.Error())
		} else {
			logrus.Infof("Update succesful. PeriodDelete = %d, IntervalDelete = %d", config.Delete.PeriodDelete, config.Delete.IntervalDelete)
		}
	}

	var floodControl FloodControl

	floodControl = fl.NewFloodControl(db, time.Second*time.Duration(config.Flood.N), config.Flood.K)

	// Создание контекста для управления жизненным циклом горутин
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Пример использования:
	// Создание 100 запросов от пользователя с интервалом в 5 секунд
	for i := 0; i < 100; i++ {
		// Проверка флуда для пользователя с ID 1
		isAllowed, err := floodControl.Check(ctx, 1)
		if err != nil {
			logrus.Fatalf("Ошибка при проверке флуда: %s", err.Error())
		}
		if !isAllowed {
			log.Println("Флуд обнаружен!")
		} else {
			log.Println("Флуд не обнаружен.")
		}

		// Ожидание 5 секунд перед следующим запросом
		time.Sleep(1 * time.Second)
	}
}
