package logtail

type Config struct {
	Port      int
	ParseKeys []string         `yaml:"parse_keys"`
	LogLevels map[string]level `yaml:"log_levels"`
}

type level struct {
	Key   string `yaml:"key"`
	Color string `yaml:"color"`
}
