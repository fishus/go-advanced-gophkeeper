package util

import (
	"bytes"
	"text/template"
)

func GetBuildInfo(date, ver string) string {
	type buildInfo struct {
		Date    string
		Version string
	}

	date = defaultValue(date, "N/A")
	ver = defaultValue(ver, "N/A")

	info := buildInfo{
		Date:    date,
		Version: ver,
	}

	const tpl = `Build date: {{.Date}}
Build version: {{.Version}}
`

	t := template.Must(template.New("list").Parse(tpl))

	buf := new(bytes.Buffer)
	err := t.Execute(buf, info)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func defaultValue(v, defaultValue string) string {
	if v == "" {
		return defaultValue
	}
	return v
}
