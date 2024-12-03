package log

func Errorf(s string, i ...any) {
	Default().Sugar().Errorf(s, i...)
}
func Error(i ...any) {
	Default().Sugar().Error(i...)
}
func Warnf(s string, i ...any) {
	Default().Sugar().Warnf(s, i...)
}
func Warn(i ...any) {
	Default().Sugar().Warn(i...)
}
func Infof(s string, i ...any) {
	Default().Sugar().Infof(s, i...)
}
func Info(i ...any) {
	Default().Sugar().Info(i...)
}
func Debugf(s string, i ...any) {
	Default().Sugar().Debugf(s, i...)
}
func Debug(i ...any) {
	Default().Sugar().Debug(i...)
}
