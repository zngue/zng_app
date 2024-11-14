package redis

type Option struct {
	Host     string
	Password string
	Port     int
	Database int
}
type Fn func(*Option)

func DataWithHost(host string) Fn {
	return func(opt *Option) {
		opt.Host = host
	}
}
func DataWithPassword(password string) Fn {
	return func(option *Option) {
		option.Password = password
	}
}
func DataWithPort(port int) Fn {
	return func(option *Option) {
		option.Port = port
	}
}
func DataWithDatabase(database int) Fn {
	return func(option *Option) {
		option.Database = database
	}
}
