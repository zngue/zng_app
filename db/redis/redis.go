package db

type RedisOption struct {
	Host     string
	Password string
	Port     int
	Database int
}
type RedisFn func(*RedisOption)

func DataWithHost(host string) RedisFn {
	return func(opt *RedisOption) {
		opt.Host = host
	}
}
func DataWithPassword(password string) RedisFn {
	return func(option *RedisOption) {
		option.Password = password
	}
}
func DataWithPort(port int) RedisFn {
	return func(option *RedisOption) {
		option.Port = port
	}
}
func DataWithDatabase(database int) RedisFn {
	return func(option *RedisOption) {
		option.Database = database
	}
}
