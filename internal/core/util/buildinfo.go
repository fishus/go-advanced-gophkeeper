package util

import (
	"os"
	"text/template"
)

func PrintBuildInfo(date, ver string) {
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

	err := t.Execute(os.Stdout, info)
	if err != nil {
		panic(err)
	}
}

func defaultValue(v, defaultValue string) string {
	if v == "" {
		return defaultValue
	}
	return v
}
