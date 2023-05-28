package model

type AppConf struct {
	Loggo *LoggoConf `yaml:"loggo"`
}

type LoggoConf struct {
	LogFilePrefix string `yaml:"log_file_prefix"`
	RotationSize  int64  `yaml:"rotation_size"`
	RotationTime  int64  `yaml:"rotation_time"`
	MaxAge        int64  `yaml:"max_age"`
}
