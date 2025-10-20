package config

type Cache struct {
	URL          string `yaml:"url" env:"LP_CACHE_URL"`
	PoolSize     int    `yaml:"pool_size" env:"LP_CACHE_POOL_SIZE"`
	ReadTimeout  int    `yaml:"read_timeout" env:"LP_CACHE_READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"LP_CACHE_WRITE_TIMEOUT"`
	MinIdleConns int    `yaml:"min_idle_conns" env:"LP_CACHE_MIN_IDLE_CONNS"`
	UseTLS       bool   `yaml:"use_tls" env:"LP_CACHE_USE_TLS"`
}
