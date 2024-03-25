package floodcontrol

type Config struct {
	DB     DBConfig    `yaml:"db"`
	Flood  FloodConfig `yaml:"flood"`
	Delete DeleteMode  `yaml:"delete_mode"`
}

type DBConfig struct {
	LocalPort  string `yaml:"local_port"`
	DockerPort string `yaml:"docker_port"`
	Password   string `yaml:"password"`
	Username   string `yaml:"username"`
	Host       string `yaml:"host"`
	SSLMode    string `yaml:"sslmode"`
	Database   string `yaml:"database"`
}

type FloodConfig struct {
	N int `yaml:"N"`
	K int `yaml:"K"`
}

type DeleteMode struct {
	IntervalDelete int `yaml:"interval_delete"`
	PeriodDelete   int `yaml:"period_delete"`
}
