package redis

import "time"

/*
 * @abstract Redis's config
 * @mail neo532@126.com
 * @date 2024-10-19
 */

type ConnectConfig struct {
	Name     string `yaml:"name"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int32  `yaml:"db"`
}

type Config struct {
	MaxSlowtime time.Duration    `yaml:"max_slowtime"`
	Default     []*ConnectConfig `yaml:"default"`
	Shadow      []*ConnectConfig `yaml:"shadow"`
	Gray        []*ConnectConfig `yaml:"gray"`
}
