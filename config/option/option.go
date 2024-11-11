package option

import (
	"github.com/zngue/zng_app/config"
	"github.com/zngue/zng_app/config/nacos"
)

type Option struct {
	GroupName    string
	NaFns        []nacos.Fn
	CFns         []config.Fn
	RegisterNaFn Fn
}
type Fn func(fn *nacos.CenterOptions) (fnErr error)

func NewOption(cfg any, option *Option) (err error) {
	var options *nacos.CenterOptions
	options, err = nacos.NewCenterOptions(nacos.NewOption(option.NaFns...))
	err = config.NewConfig(
		options,
		config.NewOption(&config.OptionConfig{
			Group: option.GroupName,
			Fns:   option.CFns,
		}),
	).Scan(&cfg)
	if err != nil {
		return
	}
	if option.RegisterNaFn != nil {
		err = option.RegisterNaFn(options)
	}
	return
}
