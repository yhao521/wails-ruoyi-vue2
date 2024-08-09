package main

import (
	"embed"
	"mySparkler/backend/run"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed config/config.yaml.example
var ymlDefault string

//go:embed sql/sqlite_data.sql
var sqlfiles embed.FS

//go:embed view/template
var templates embed.FS

func main() {
	run.WailsRun(assets, ymlDefault, sqlfiles, templates)
}
