module github.com/sh0rez/klint

go 1.15

require (
	github.com/containous/yaegi v0.8.15-0.20200821085603-358a57b4b9bc
	github.com/fatih/color v1.9.0
	github.com/go-clix/cli v0.1.2
	github.com/grafana/tanka v0.11.1
	github.com/pkg/errors v0.8.1
	gopkg.in/yaml.v3 v3.0.0-20191010095647-fc94e3f71652
)

replace github.com/grafana/tanka => ../tanka
