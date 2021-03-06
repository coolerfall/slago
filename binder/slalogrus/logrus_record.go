// Copyright (c) 2019-2021 Vincent Cheung (coolingfall@gmail.com).
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

package slalogrus

import (
	"encoding/hex"
	"sync"
	"time"

	"github.com/coolerfall/slago"
	"github.com/sirupsen/logrus"
)

var (
	recordPool = &sync.Pool{
		New: func() interface{} {
			return &logrusRecord{}
		},
	}
)

type logrusRecord struct {
	entry *logrus.Entry
	level logrus.Level
}

func newLogrusRecord(lvl logrus.Level) *logrusRecord {
	r := recordPool.Get().(*logrusRecord)
	r.entry = logrus.NewEntry(logrus.StandardLogger())
	r.level = lvl

	return r
}

func (r *logrusRecord) Str(key, val string) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Strs(key string, val []string) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Bytes(key string, val []byte) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Hex(key string, val []byte) slago.Record {
	r.Interface(key, hex.EncodeToString(val))
	return r
}

func (r *logrusRecord) Err(err error) slago.Record {
	r.entry = r.entry.WithError(err)
	return r
}

func (r *logrusRecord) Errs(key string, errs []error) slago.Record {
	r.entry = r.entry.WithField(key, errs)
	return r
}

func (r *logrusRecord) Bool(key string, val bool) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Bools(key string, val []bool) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Int(key string, val int) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Ints(key string, val []int) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Int8(key string, val int8) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Ints8(key string, val []int8) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Int16(key string, val int16) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Ints16(key string, val []int16) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Int32(key string, val int32) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Ints32(key string, val []int32) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Int64(key string, val int64) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Ints64(key string, val []int64) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Uint(key string, val uint) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Uints(key string, val []uint) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Uint8(key string, val uint8) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Uints8(key string, val []uint8) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Uint16(key string, val uint16) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Uints16(key string, val []uint16) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Uint32(key string, val uint32) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Uints32(key string, val []uint32) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Uint64(key string, val uint64) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Uints64(key string, val []uint64) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Float32(key string, val float32) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Floats32(key string, val []float32) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Float64(key string, val float64) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Floats64(key string, val []float64) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Time(key string, val time.Time) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Times(key string, val []time.Time) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Dur(key string, val time.Duration) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Durs(key string, val []time.Duration) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Interface(key string, val interface{}) slago.Record {
	r.entry = r.entry.WithField(key, val)
	return r
}

func (r *logrusRecord) Msg(originMsg ...string) {
	var msg string
	if len(originMsg) != 0 {
		msg = originMsg[0]
	}
	r.entry.Log(r.level, msg)
	recordPool.Put(r)
}

func (r *logrusRecord) Msgf(format string, v ...interface{}) {
	r.entry.Logf(r.level, format, v...)
	recordPool.Put(r)
}
