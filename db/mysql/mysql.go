package db

type MysqlOption struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}
type MysqlFn func(opt *MysqlOption)

func WithUserName(username string) MysqlFn {
	return func(opt *MysqlOption) {
		opt.Username = username
	}
}
func WithPassword(password string) MysqlFn {
	return func(opt *MysqlOption) {
		opt.Password = password
	}
}
func WithHost(host string) MysqlFn {
	return func(opt *MysqlOption) {
		opt.Host = host
	}
}
func WithPort(port int) MysqlFn {
	return func(opt *MysqlOption) {
		opt.Port = port
	}
}
func WithDatabase(database string) MysqlFn {
	return func(opt *MysqlOption) {
		opt.Database = database
	}
}
