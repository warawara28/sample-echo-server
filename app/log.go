package app

import (
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Logger logger用インターフェイス
type Logger interface {
	Fatal(format string, v ...interface{})
	Error(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

// Zerolog zerologを使用したlogger実装
type Zerolog struct {
	l zerolog.Logger
}

// NewZerolog Zerologのインスタンスを生成
func NewZerolog() *Zerolog {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		switch l {
		case zerolog.DebugLevel:
			return "DEBUG"
		case zerolog.InfoLevel:
			return "INFO"
		case zerolog.WarnLevel:
			return "WARNING"
		case zerolog.ErrorLevel:
			return "ERROR"
		case zerolog.FatalLevel:
			return "FATAL"
		case zerolog.PanicLevel:
			return "PANIC"
		case zerolog.TraceLevel:
			return "TRACE"
		case zerolog.Disabled:
			return "DISABLED"
		default:
			return "UNKNOWN"
		}
	}
	return &Zerolog{l: zerolog.New(os.Stdout).With().Timestamp().Logger()}
}

// Fatal .
func (z *Zerolog) Fatal(format string, v ...interface{}) {
	z.l.Fatal().Msgf(format, v...)
}

// Error .
func (z *Zerolog) Error(format string, v ...interface{}) {
	z.l.Error().Msgf(format, v...)
}

// Warning .
func (z *Zerolog) Warning(format string, v ...interface{}) {
	z.l.Warn().Msgf(format, v...)
}

// Info .
func (z *Zerolog) Info(format string, v ...interface{}) {
	z.l.Info().Msgf(format, v...)
}

// Debug .
func (z *Zerolog) Debug(format string, v ...interface{}) {
	z.l.Debug().Msgf(format, v...)
}

func (z *Zerolog) getLogger() *zerolog.Logger {
	return &z.l
}

const (
	FieldIP       = "ip"
	FieldURI      = "uri"
	FieldHost     = "host"
	FieldMethod   = "method"
	FieldPath     = "path"
	FieldProtocol = "protocol"
	FieldReferer  = "referer"
	FieldUa       = "ua"
	FieldStatus   = "status"
	FieldElapsed  = "elapsed"
)

var defaultFields = []string{
	FieldIP,
	FieldURI,
	FieldHost,
	FieldMethod,
	FieldPath,
	FieldProtocol,
	FieldReferer,
	FieldUa,
	FieldStatus,
	FieldElapsed,
}

type fetchFieldFunc = func(ctx echo.Context) string

var fetchFieldFunctions = map[string]func(ctx echo.Context) string{
	FieldIP:       func(ctx echo.Context) string { return ctx.RealIP() },
	FieldURI:      func(ctx echo.Context) string { return ctx.Request().RequestURI },
	FieldHost:     func(ctx echo.Context) string { return ctx.Request().Host },
	FieldMethod:   func(ctx echo.Context) string { return ctx.Request().Method },
	FieldPath:     func(ctx echo.Context) string { return ctx.Request().URL.Path },
	FieldProtocol: func(ctx echo.Context) string { return ctx.Request().Proto },
	FieldReferer:  func(ctx echo.Context) string { return ctx.Request().Referer() },
	FieldUa:       func(ctx echo.Context) string { return ctx.Request().UserAgent() },
	FieldStatus:   func(ctx echo.Context) string { return strconv.Itoa(ctx.Response().Status) },
	FieldElapsed:  func(ctx echo.Context) string { return strconv.Itoa(int(ctx.Get(FieldElapsed).(int64))) + " ns" },
}

// NewZerologMiddleware ログ出力のデフォルト設定
func NewZerologMiddleware(logger *Zerolog) echo.MiddlewareFunc {
	return NewZerologMiddlewareWithFields(logger, defaultFields)
}

// NewZerologMiddlewareWithFields ログ出力のデフォルト設定、出力項目を設定可能
func NewZerologMiddlewareWithFields(logger *Zerolog, fields []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			start := time.Now()

			if err = next(ctx); err != nil {
				ctx.Error(err)
			}

			elapsedNano := time.Now().Sub(start).Nanoseconds()
			ctx.Set(FieldElapsed, elapsedNano)

			dict := zerolog.Dict()
			for _, key := range fields {
				fetchValue, ok := fetchFieldFunctions[key]
				if ok {
					dict.Str(key, fetchValue(ctx))
				}
			}
			logger.getLogger().Info().Dict("request", dict).Send()

			return
		}
	}
}
