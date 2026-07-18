package main

import "github.com/MindHunter86/eyesonly/internal/utils"

var (
	buildtime = "never"

	name    = "eyesonly"
	version = utils.DevelVersionIdent // -ldflags="-X main.version=X.X.X"
	usage   = "One Time service for short-lived secrets"

	copyright = "(c) 2026 mindhunter86\nQualification work for CBC. 2026 Spring"
)
