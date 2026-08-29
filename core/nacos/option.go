package nacos

type Option struct {
	NamespaceId         string
	LogDir              string
	CacheDir            string
	LogLevel            LogType
	AppendToStdout      bool
	Username            string
	Password            string
	Host                string
	Port                int
	GrpcPort            int
	TimeoutMs           uint64
	NotLoadCacheAtStart bool
}

type OptionFunc func(*Option)

func WithNamespace(namespaceId string) OptionFunc {
	return func(opt *Option) {
		opt.NamespaceId = namespaceId
	}
}

func WithHost(host string) OptionFunc {
	return func(opt *Option) {
		opt.Host = host
	}
}

func WithPort(port int) OptionFunc {
	return func(opt *Option) {
		opt.Port = port
	}
}

func WithGrpcPort(port int) OptionFunc {
	return func(opt *Option) {
		opt.GrpcPort = port
	}
}

func WithTimeoutMs(ms uint64) OptionFunc {
	return func(opt *Option) {
		opt.TimeoutMs = ms
	}
}

func WithUserName(userName string) OptionFunc {
	return func(opt *Option) {
		opt.Username = userName
	}
}

func WithPassword(password string) OptionFunc {
	return func(opt *Option) {
		opt.Password = password
	}
}

func WithLogDir(s string) OptionFunc {
	return func(opt *Option) {
		opt.LogDir = s
	}
}

func WithCacheDir(s string) OptionFunc {
	return func(opt *Option) {
		opt.CacheDir = s
	}
}

func WithLogLevel(s LogType) OptionFunc {
	return func(opt *Option) {
		opt.LogLevel = s
	}
}

func WithAppendToStdout(b bool) OptionFunc {
	return func(opt *Option) {
		opt.AppendToStdout = b
	}
}

func WithNotLoadCacheAtStart(b bool) OptionFunc {
	return func(opt *Option) {
		opt.NotLoadCacheAtStart = b
	}
}

func NewOption(fns ...OptionFunc) *Option {
	opt := &Option{
		NamespaceId:    "develop",
		LogDir:         "nacos/log",
		CacheDir:       "nacos/cache",
		TimeoutMs:      5000,
		AppendToStdout: false,
		Port:           8848,
		GrpcPort:       9848,
	}
	for _, fn := range fns {
		fn(opt)
	}
	return opt
}
