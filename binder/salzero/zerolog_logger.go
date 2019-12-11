// Copyright (c) 2019 Anbillon Team (anbillonteam@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package salzero

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/anbillon/slago/slago-api"
)

var (
	slagoLvlToZeroLvl = map[slago.Level]zerolog.Level{
		slago.TraceLevel: zerolog.NoLevel,
		slago.DebugLevel: zerolog.DebugLevel,
		slago.InfoLevel:  zerolog.InfoLevel,
		slago.WarnLevel:  zerolog.WarnLevel,
		slago.ErrorLevel: zerolog.ErrorLevel,
		slago.FatalLevel: zerolog.FatalLevel,
		slago.PanicLevel: zerolog.PanicLevel,
	}
)

// zeroLogger is an implementation of SlaLogger.
type zeroLogger struct {
	logger      zerolog.Logger
	multiWriter *slago.MultiWriter
}

// NewZeroLogger creates a new instance of zeroLogger used to be bound to slago.
func NewZeroLogger() *zeroLogger {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = slago.TimestampFormat
	zerolog.LevelFieldName = slago.LevelFieldKey
	zerolog.TimestampFieldName = slago.TimestampFieldKey
	zerolog.MessageFieldName = slago.MessageFieldKey
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return strings.ToUpper(l.String())
	}

	multiWriter := slago.NewMultiWriter()
	logger := zerolog.New(multiWriter).With().Timestamp().Logger()
	log.Logger = logger

	return &zeroLogger{
		logger:      logger,
		multiWriter: multiWriter,
	}
}

func (l *zeroLogger) Name() string {
	return "zerolog"
}

func (l *zeroLogger) AddWriter(w ...slago.Writer) {
	l.multiWriter.AddWriter(w...)
}

func (l *zeroLogger) ResetWriter() {
	l.multiWriter.Reset()
}

func (l *zeroLogger) SetLevel(lvl slago.Level) {
	zeroLevel := slagoLvlToZeroLvl[lvl]
	if zeroLevel == zerolog.NoLevel {
		zeroLevel = zerolog.TraceLevel
	}

	zerolog.SetGlobalLevel(slagoLvlToZeroLvl[lvl])
}

func (l *zeroLogger) Level(lvl slago.Level) slago.Record {
	zeroLevel := slagoLvlToZeroLvl[lvl]
	if zeroLevel == zerolog.NoLevel {
		return l.Trace()
	}

	return newZeroRecord(l.logger.WithLevel(zeroLevel))
}

func (l *zeroLogger) Trace() slago.Record {
	return newZeroRecord(l.logger.WithLevel(zerolog.NoLevel).
		Str(slago.LevelFieldKey, "TRACE"))
}

func (l *zeroLogger) Debug() slago.Record {
	return newZeroRecord(l.logger.Debug())
}

func (l *zeroLogger) Info() slago.Record {
	return newZeroRecord(l.logger.Info())
}

func (l *zeroLogger) Warn() slago.Record {
	return newZeroRecord(l.logger.Warn())
}

func (l *zeroLogger) Error() slago.Record {
	return newZeroRecord(l.logger.Error())
}

func (l *zeroLogger) Fatal() slago.Record {
	return newZeroRecord(l.logger.Fatal())
}

func (l *zeroLogger) Panic() slago.Record {
	return newZeroRecord(l.logger.Panic())
}

func (l *zeroLogger) Print(v ...interface{}) {
	l.logger.Print(v...)
}

func (l *zeroLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *zeroLogger) WriteRaw(p []byte) {
	_, err := l.multiWriter.Write(p)
	if err != nil {
		l.Error().Err(err).Msg("write raw error")
	}
}