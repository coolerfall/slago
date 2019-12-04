// Copyright (c) 2019 Anbillon Team (anbillonteam@gmail.com).

package zeroslago

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/anbillon/slago/slago-api"
)

var (
	zapLvlToSlagoLvl = map[zerolog.Level]slago.Level{
		zerolog.NoLevel:    slago.TraceLevel,
		zerolog.DebugLevel: slago.DebugLevel,
		zerolog.InfoLevel:  slago.InfoLevel,
		zerolog.WarnLevel:  slago.WarnLevel,
		zerolog.ErrorLevel: slago.ErrorLevel,
		zerolog.FatalLevel: slago.FatalLevel,
	}
)

type zerologBridge struct {
}

// NewZerologBridge creates a new slago bridge for zerolog.
func NewZerologBridge() slago.Bridge {
	bridge := &zerologBridge{}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = slago.TimestampFormat
	zerolog.LevelFieldName = slago.LevelFieldKey
	zerolog.TimestampFieldName = slago.TimestampFieldKey
	zerolog.MessageFieldName = slago.MessageFieldKey
	logger := zerolog.New(bridge).With().Timestamp().Logger()
	log.Logger = logger

	return bridge
}

func (b *zerologBridge) Name() string {
	return "zerolog"
}

func (b *zerologBridge) ParseLevel(lvl string) slago.Level {
	level, err := zerolog.ParseLevel(lvl)
	if err != nil {
		level = zerolog.NoLevel
		slago.Reportf("parse zerolog level error: %s", err)
	}

	return zapLvlToSlagoLvl[level]
}

func (b *zerologBridge) Write(p []byte) (int, error) {
	err := slago.BrigeWrite(b, p)
	if err != nil {
		slago.Reportf("zerolog bridge write error", err)
	}

	return len(p), err
}
