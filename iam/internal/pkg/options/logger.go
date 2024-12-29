package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

type LogOption struct {
	Level string `json:"level"` // debug error warning
}

func NewLogOption() *LogOption {
	return &LogOption{
		Level: "debug",
	}
}

func (log *LogOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&log.Level, "log.level", log.Level,
		"选择日志等级进行打印，可选 err || error , warn || warning, info, debug , 默认为debug")
}

func (log *LogOption) Validate() []error {

	var err []error

	switch log.Level {
	case "err", "error", "warn", "warning", "info", "debug":

	default:
		err = append(err, fmt.Errorf("unknown log level type"))
	}

	return err
}
