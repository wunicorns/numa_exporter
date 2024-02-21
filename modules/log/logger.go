package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const LOG_FILENAME = "/ttm-datareceiver.log"

type Config struct {
	// log file
	LogDir        string `toml:"LOG_DIR"`
	LogMaxSize    int    `toml:"LOG_MAX_SIZE"`
	LogMaxBackups int    `toml:"LOG_MAX_BACKUPS"`
	LogMaxAge     int    `toml:"LOG_MAX_AGE"`
	LogCompress   bool   `toml:"LOG_COMPRESS"`
	DebugLevel    string `toml:"DEBUG_LEVEL"`
}

// Log is the library wide logger. Setting to nil disables logging.
// var Log Logger = &logger{logLevel: ErrorLevel, prefix: "ttmBatch"}

// Logger defines interface for logging
type Logger interface {
	// Writes formatted debug message if debug logLevel is enabled.
	Debugf(format string, v ...interface{})
	// Writes debug message if debug is enabled.
	Debug(msg string)
	// Writes formatted info message if info logLevel is enabled.
	Infof(format string, v ...interface{})
	// Writes info message if info logLevel is enabled
	Info(msg string)
	// Writes formatted warning message if warning logLevel is enabled.
	Warnf(format string, v ...interface{})
	// Writes warning message if warning logLevel is enabled.
	Warn(msg string)
	// Writes formatted error message
	Errorf(format string, v ...interface{})
	// Writes error message
	Error(msg string)
	// SetLogLevel sets allowed logging level.
	SetLogLevel(logLevel uint)
	// LogLevel retrieves current logging level
	LogLevel() uint
	// SetPrefix sets logging prefix.
	SetPrefix(prefix string)
}

// logger provides default implementation for Logger. It logs using Go log API
// mutex is needed in cases when multiple clients run concurrently
// type logger struct {
// 	prefix   string
// 	logLevel uint
// 	lock     sync.Mutex
// }

var (
	logger *logrus.Logger
)

// logrus hook for rotate
type hook struct {
	rotateCh    chan struct{}
	currentDate string
	recentDate  string
}

func (h *hook) Fire(entry *logrus.Entry) error {
	h.currentDate = entry.Time.Format("2006-01-02")
	if h.currentDate != h.recentDate {
		h.rotateCh <- struct{}{}
	}
	h.recentDate = h.currentDate
	return nil
}
func (h *hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// logrus formatter
type fomatter struct {
	logrus.TextFormatter
}

func (f *fomatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s [%s] - %s\n", entry.Time.Format(f.TimestampFormat), strings.ToUpper(entry.Level.String()), entry.Message)), nil
}

func Init(config Config) error {

	// lumberjack
	lum := &lumberjack.Logger{
		Filename:   config.LogDir + LOG_FILENAME,
		MaxSize:    config.LogMaxSize,    // log파일의 최대 사이즈 (MB)
		MaxBackups: config.LogMaxBackups, // 보존 할 최대 이전 로그 파일 수
		MaxAge:     config.LogMaxAge,     // 타임 스탬프를 기준으로 오래된 로그 파일을 보관할 수있는 최대 일수
		Compress:   config.LogCompress,   // 압축 여부 (default: false)
	}

	// logrus
	logger = &logrus.Logger{
		Out:   io.MultiWriter(os.Stdout, lum),
		Hooks: make(logrus.LevelHooks),
		Level: logrus.InfoLevel,
		Formatter: &fomatter{logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05.000",
			ForceColors:            false,
			DisableLevelTruncation: true,
		},
		},
	}

	// make channel and go routine for rotate log file.
	rotateCh := make(chan struct{})
	go func() {
		for {
			<-rotateCh
			lum.Rotate()
		}
	}()

	// add hook that determines whether rotate log file or not with channel
	logger.Hooks.Add(&hook{
		rotateCh:   rotateCh,
		recentDate: time.Now().Format("2006-01-02"),
	})

	SetLogLevel(config.DebugLevel)

	return nil
}

func Debugf(v ...interface{}) {
	Debug(v)
}

func Infof(v ...interface{}) {
	Info(v)
}

func Warnf(v ...interface{}) {
	Warn(v)
}

func Errorf(v ...interface{}) {
	Error(v)
}

func Debug(v ...interface{}) {
	v = addSpaceSep(v...)
	logger.Debug(v...)
}

func Info(v ...interface{}) {
	v = addSpaceSep(v...)
	logger.Info(v...)
}

func Warn(v ...interface{}) {
	v = addSpaceSep(v...)
	logger.Warn(v...)
}

func Error(v ...interface{}) {
	caller := getCaller()
	s := []interface{}{}
	s = append(s, caller)
	s = append(s, v...)

	v = addSpaceSep(s...)
	logger.Error(v...)
}

func addSpaceSep(val ...interface{}) []interface{} {
	s := []interface{}{}
	for _, v := range val {
		if len(s) > 0 {
			s = append(s, " ")
		}
		s = append(s, v)
	}
	return s
}

func SetLogLevel(lv string) {
	var v logrus.Level
	switch strings.ToUpper(lv) {
	case "DEBUG":
		v = logrus.DebugLevel
	case "INFO":
		v = logrus.InfoLevel
	case "WARN":
		v = logrus.WarnLevel
	case "ERROR":
		v = logrus.ErrorLevel
	default:
		v = logrus.InfoLevel
	}
	logger.SetLevel(v)
}

func IsDebugLevel() bool {
	if logger.GetLevel() == logrus.DebugLevel {
		return true
	}
	return false
}

func getCaller() string {
	pc, file, line, ok := runtime.Caller(2)

	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".")
	}

	return fmt.Sprintf("[%s:%d %s]", filepath.Base(file), line, fnName)
}
