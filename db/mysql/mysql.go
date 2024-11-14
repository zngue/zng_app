package mysql

type Option struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}
type Fn func(opt *Option)

func DataWithUserName(username string) Fn {
	return func(opt *Option) {
		opt.Username = username
	}
}
func DataWithPassword(password string) Fn {
	return func(opt *Option) {
		opt.Password = password
	}
}
func DataWithHost(host string) Fn {
	return func(opt *Option) {
		opt.Host = host
	}
}
func DataWithPort(port int) Fn {
	return func(opt *Option) {
		opt.Port = port
	}
}
func DataWithDatabase(database string) Fn {
	return func(opt *Option) {
		opt.Database = database
	}
}
