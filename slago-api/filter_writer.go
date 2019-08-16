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
	jsoniter "github.com/json-iterator/go"
	"io"
)

type FilterWriter struct {
	originWirters []io.Writer
}

func NewFilterWriter(w ...io.Writer) *FilterWriter {
	return &FilterWriter{
		originWirters: w,
	}
}

func (w *FilterWriter) Write(p []byte) (int, error) {
	var event map[string]interface{}
	jsoniter.Unmarshal(p, &event)

	return len(p), nil
}
