package internal

import "fmt"

type Config struct {
	Port       int
	FileToTail string           `yaml:"file_to_tail"`
	ParseKeys  []string         `yaml:"parse_keys"`
	LogLevels  map[string]level `yaml:"log_levels"`
}

func (c Config) String() string {
	return fmt.Sprintf("port: %v, file_to_tail: %v, parse_keys: %v, log_levels: %v", c.Port, c.FileToTail,
		c.ParseKeys, c.LogLevels)
}

type level struct {
	Key   string `yaml:"key"`
	Color string `yaml:"color"`
}

func (l level) String() string {
	return fmt.Sprintf("key: %v, color: %v", l.Key, l.Color)
}
