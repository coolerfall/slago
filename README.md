slago
=====
Simple Logging Abstraction for Go. 

Usage
====
* Add logger you want to bind to:
```go
slago.Bind(salzero.NewZeroLogger())
```

* Install the bridges for other logger :
```go
slago.Install(bridge.NewLogBridge())
slago.Install(bridge.NewLogrusBridge())
slago.Install(bridge.NewZapBrige())
```

* Configure the output writer:
```go
cw := slago.NewConsoleWriter(func(o *slago.ConsoleWriterOption) {
		o.Encoder = slago.NewPatternEncoder(
			"#color(#date{2006-01-02T15:04:05.000Z07:00}){cyan} #color(" +
				"#level) #message #fields")
	})
slago.Logger().AddWriter(cw)
fw := slago.NewFileWriter(func(o *slago.FileWriterOption) {
		o.Encoder = slago.NewJsonEncoder()
		o.Filter = slago.NewLevelFilter(slago.TraceLevel)
		o.Filename = "slago-test.log"
		o.RollingPolicy = slago.NewSizeAndTimeBasedRollingPolicy(
			func(o *slago.SizeAndTimeBasedRPOption) {
				o.FilenamePattern = "slago-archive.#date{2006-01-02}.#index.log"
				o.MaxFileSize = "10MB"
			})
	})
slago.Logger().AddWriter(fw)
```

* Add logging...:
```go
slago.Logger().Trace().Msg("slago")
slago.Logger().Info().Int("int", 88).Interface("slago", "val").Msg("")
```

* If you log with other logger, it will send to the bound logger:
```go
zap.L().With().Warn("this is zap")
log.Printf("this is builtin logger")
```

License
=======

    Copyright (C) 2019 Anbillon Team
    
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
         http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
