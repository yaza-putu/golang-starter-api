package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/yaza-putu/golang-starter-api/internal/config"
)

const (
	DEBUG Lvl = iota + 1
	INFO
	WARN
	ERROR
	OFF
	PANIC
	FATAL
)

type (
	optFunc func(opts2 *opts)
	Lvl     uint8
	opts    struct {
		Write bool
		Type  Lvl
	}
)

func defaultOpts() opts {
	return opts{
		Write: true,
		Type:  ERROR,
	}
}

func IsWrite(r bool) optFunc {
	return func(o *opts) {
		o.Write = r
	}
}

func SetType(t Lvl) optFunc {
	return func(o *opts) {
		o.Type = t
	}
}

func New(err error, opts ...optFunc) {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	if err != nil {
		if o.Write == true && config.App().Status != "test" {
			writeError(err, o.Type)
		}
		if config.App().Debug == true {
			fmt.Printf("%s : %s", getlabel(o.Type), err.Error())
		}
	}
}

func getlabel(l Lvl) string {
	switch l {
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case FATAL:
		return "FATAL"
	case PANIC:
		return "PANIC"
	default:
		return "ERROR"
	}
}

func writeError(err error, l Lvl) {
	cwd, _ := os.Getwd()
	filename := fmt.Sprintf("%s-error.log", time.Now().Format("2006-01-02"))
	logFile, er := os.OpenFile(fmt.Sprintf("logs/%s", filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if er != nil {
		_ = os.Mkdir("logs", os.ModePerm)
		path := filepath.Join(cwd, "logs", filename)
		newFilePath := filepath.FromSlash(path)
		logFile, er = os.Create(newFilePath)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	switch l {
	case INFO:
		log.Info(err.Error())
		break
	case DEBUG:
		log.Debug(err.Error())
		break
	case FATAL:
		log.Fatal(err.Error())
		break
	case PANIC:
		log.Panic(err.Error())
		break
	default:
		log.Error(err.Error())
		break
	}
}
