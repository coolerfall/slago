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

package bridge

import (
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/anbillon/slago/slago-api"
)

var (
	logrusLvlToSlagoLvl = map[logrus.Level]slago.Level{
		logrus.TraceLevel: slago.TraceLevel,
		logrus.DebugLevel: slago.DebugLevel,
		logrus.InfoLevel:  slago.InfoLevel,
		logrus.WarnLevel:  slago.WarnLevel,
		logrus.ErrorLevel: slago.ErrorLevel,
		logrus.FatalLevel: slago.FatalLevel,
	}
)

type logrusBridge struct {
}

// NewLogrusBridge creates a new slago bridge for logrus.
func NewLogrusBridge() slago.Bridge {
	bridge := &logrusBridge{}
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: slago.LevelFieldKey,
			logrus.FieldKeyTime:  slago.TimestampFieldKey,
			logrus.FieldKeyMsg:   slago.MessageFieldKey,
		},
	})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(bridge)

	return bridge
}

func (b *logrusBridge) Name() string {
	return "logrus"
}

func (b *logrusBridge) ParseLevel(lvl string) slago.Level {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		slago.Reportf("parse logrus level error: %s", err)
		level = logrus.TraceLevel
	}

	return logrusLvlToSlagoLvl[level]
}

func (b *logrusBridge) Write(p []byte) (int, error) {
	err := slago.BrigeWrite(b, p)
	if err != nil {
		slago.Reportf("logrus bridge write error", err)
	}

	return len(p), err
}