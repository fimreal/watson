package crond

import (
	"github.com/fimreal/goutils/ezap"
	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func init() {
	l := logger{ezap.New()}
	c = cron.New(cron.WithLogger(cron.VerbosePrintfLogger(l)))
}

type logger struct {
	*ezap.Logger
}

func (l logger) Printf(msg string, keysAndValues ...interface{}) {
	l.Logger.Logger.Infof(msg, keysAndValues...)
}

func Run(spec string, f func()) (cron.EntryID, error) {
	id, err := c.AddFunc(spec, func() { f() })
	if err != nil {
		return 0, err
	}
	c.Start()
	return id, nil
}
