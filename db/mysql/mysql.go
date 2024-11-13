package db

type MysqlOption struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}
type MysqlFn func(opt *MysqlOption)

func DataWithUserName(username string) MysqlFn {
	return func(opt *MysqlOption) {
		opt.Username = username
	}
}
func DataWithPassword(password string) MysqlFn {
	return func(opt *MysqlOption) {
		opt.Password = password
	}
}
func DataWithHost(host string) MysqlFn {
	return func(opt *MysqlOption) {
		opt.Host = host
	}
}
func DataWithPort(port int) MysqlFn {
	return func(opt *MysqlOption) {
		opt.Port = port
	}
}
func DataWithDatabase(database string) MysqlFn {
	return func(opt *MysqlOption) {
		opt.Database = database
	}
}
