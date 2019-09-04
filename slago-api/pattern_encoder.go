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

package slago

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	DefaultLayout = "#color(#date{2006-01-02}){cyan} #color(#level) #message #fields"
)

var (
	colorMap = map[string]int{
		"black":     colorBlack,
		"red":       colorRed,
		"green":     colorGreen,
		"yellow":    colorYellow,
		"blue":      colorBlue,
		"magenta":   colorMagenta,
		"cyan":      colorCyan,
		"white":     colorWhite,
		"blackbr":   colorBrightBlack,
		"redbr":     colorBrightRed,
		"greenbr":   colorBrightGreen,
		"yellowbr":  colorBrightYellow,
		"bluebr":    colorBrightBlue,
		"magentabr": colorBrightMagenta,
		"cyanbr":    colorBrightCyan,
		"whitebr":   colorBrightWhite,
	}

	levelColorMap = map[string]int{
		"TRACE": colorWhite,
		"DEBUG": colorBlue,
		"INFO":  colorGreen,
		"WARN":  colorYellow,
		"ERROR": colorRed,
		"FATAL": colorRed,
		"PANIC": colorRed,
	}
)

// PatternEncoder encodes logging event with pattern.
type PatternEncoder struct {
	mutex     sync.Mutex
	buf       *bytes.Buffer
	converter Converter
}

// NewPatternEncoder creates a new instance of pattern encoder.
func NewPatternEncoder(layouts ...string) *PatternEncoder {
	var layout string
	if len(layouts) == 0 || len(layouts[0]) == 0 {
		layout = DefaultLayout
	} else {
		layout = layouts[0]
	}

	patternParser := NewPatternParser(layout)
	node, err := patternParser.Parse()
	if err != nil {
		ReportfExit("parse pattern error, %v", err)
	}

	converters := map[string]NewConverter{
		"color":   newColorConverter,
		"level":   newLevelConverter,
		"date":    newLogDateConverter,
		"message": newMessageConverter,
		"fields":  newFieldsConverter,
	}
	converter, err := NewPatternCompiler(node, converters).Compile()
	if err != nil {
		ReportfExit("compile pattern error, %v", err)
	}

	return &PatternEncoder{
		buf:       &bytes.Buffer{},
		converter: converter,
	}
}

func (pe *PatternEncoder) Encode(p []byte) (data []byte, err error) {
	var eventMap map[string]interface{}
	err = json.Unmarshal(p, &eventMap)
	if err != nil {
		return nil, err
	}

	lvl := getAndRemove(LevelFieldKey, eventMap)
	ts := getAndRemove(TimestampFieldKey, eventMap)
	msg := getAndRemove(MessageFieldKey, eventMap)
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, err
	}
	level := ParseLevel(lvl)

	event := logEvent{
		level:     level,
		timestamp: t,
		message:   msg,
		fields:    eventMap,
	}
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	for c := pe.converter; c != nil; c = c.Next() {
		pe.buf.WriteString(c.Convert(event))
	}
	pe.buf.WriteByte('\n')
	data = pe.buf.Bytes()
	pe.buf.Reset()

	return data, err
}

type logEvent struct {
	level     Level
	timestamp time.Time
	message   string
	fields    map[string]interface{}
}

type colorConverter struct {
	next  Converter
	child Converter
	opts  []string
	buf   *bytes.Buffer
}

func newColorConverter() Converter {
	return &colorConverter{
		buf: new(bytes.Buffer),
	}
}

func (cc *colorConverter) AttatchNext(next Converter) {
	cc.next = next
}

func (cc *colorConverter) Next() Converter {
	return cc.next
}

func (cc *colorConverter) AttachChild(child Converter) {
	cc.child = child
}

func (cc *colorConverter) AttachOptions(opts []string) {
	cc.opts = opts
}

func (cc *colorConverter) Convert(event interface{}) string {
	var level string
	if logEvent, ok := event.(logEvent); ok {
		level = logEvent.level.String()
	}

	if len(cc.opts) != 0 {
		color, ok := colorMap[cc.opts[0]]
		if !ok {
			color = colorWhite
		}
		cc.writeColor(color)
	}

	for c := cc.child; c != nil; c = c.Next() {
		result := c.Convert(event)
		if _, ok := c.(*levelConverter); ok {
			color, ok := levelColorMap[level]
			if !ok {
				color = colorWhite
			}

			cc.writeColor(color)
		}
		cc.buf.WriteString(result)
		cc.writeColorEnd()
	}

	cc.writeColorEnd()

	data := cc.buf.String()
	cc.buf.Reset()

	return data
}

func (cc *colorConverter) writeColor(color int) {
	cc.buf.WriteString("\x1b[")
	cc.buf.WriteString(strconv.Itoa(color))
	cc.buf.WriteByte('m')
}

func (cc *colorConverter) writeColorEnd() {
	cc.buf.WriteString("\x1b[0m")
}

type levelConverter struct {
	next Converter
}

func newLevelConverter() Converter {
	return &levelConverter{}
}

func (lc *levelConverter) AttatchNext(next Converter) {
	lc.next = next
}

func (lc *levelConverter) Next() Converter {
	return lc.next
}

func (lc *levelConverter) AttachChild(child Converter) {
}

func (lc *levelConverter) AttachOptions(opts []string) {
}

func (lc *levelConverter) Convert(event interface{}) string {
	logEvent, ok := event.(logEvent)
	if !ok {
		return ""
	}

	return logEvent.level.String()
}

type logDateConverter struct {
	next  Converter
	child Converter
	opts  []string
}

func newLogDateConverter() Converter {
	return &logDateConverter{
		opts: []string{"2006-01-02"},
	}
}

func (c *logDateConverter) AttatchNext(next Converter) {
	c.next = next
}

func (c *logDateConverter) Next() Converter {
	return c.next
}

func (c *logDateConverter) AttachChild(child Converter) {
	c.child = child
}

func (c *logDateConverter) AttachOptions(opts []string) {
	if len(opts) != 0 && len(opts[0]) != 0 {
		c.opts = opts
	}
}

func (c *logDateConverter) Convert(event interface{}) string {
	logEvent, ok := event.(logEvent)
	if !ok {
		return ""
	}

	return logEvent.timestamp.Format(c.opts[0])
}

type messageConverter struct {
	next Converter
}

func newMessageConverter() Converter {
	return &messageConverter{}
}

func (mc *messageConverter) AttatchNext(next Converter) {
	mc.next = next
}

func (mc *messageConverter) Next() Converter {
	return mc.next
}

func (mc *messageConverter) AttachChild(child Converter) {
}

func (mc *messageConverter) AttachOptions(opts []string) {
}

func (mc *messageConverter) Convert(event interface{}) string {
	logEvent, ok := event.(logEvent)
	if !ok {
		return "-"
	}

	message := logEvent.message
	if len(message) == 0 {
		return "-"
	}

	return message
}

type fieldsConverter struct {
	next Converter
	buf  *bytes.Buffer
}

func newFieldsConverter() Converter {
	return &fieldsConverter{
		buf: new(bytes.Buffer),
	}
}

func (fc *fieldsConverter) AttatchNext(next Converter) {
	fc.next = next
}

func (fc *fieldsConverter) Next() Converter {
	return fc.next
}

func (fc *fieldsConverter) AttachChild(child Converter) {
}

func (fc *fieldsConverter) AttachOptions(opts []string) {
}

func (fc *fieldsConverter) Convert(event interface{}) string {
	logEvent, ok := event.(logEvent)
	if !ok {
		return ""
	}

	fields := logEvent.fields
	for k, v := range fields {
		fc.buf.WriteString(fmt.Sprintf("%s=%v ", k, v))
	}
	data := fc.buf.String()
	fc.buf.Reset()

	return data
}
