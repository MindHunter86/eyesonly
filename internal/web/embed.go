package web

import (
	"embed"
	"io/fs"
)

//go:generate go install github.com/valyala/quicktemplate/qtc
//go:generate qtc -skipLineComments -dir=.

//go:generate npm ci --include dev --no-audit --no-fund --progress=false
//go:generate node node_modules/eslint/bin/eslint.js ./src/
//go:generate npm run build
//go:generate ls -la dist/

////go:generate gzip -f -9 -k ./dist/js/*js
////go:generate zstd -f -k --ultra -22 static/js/main.mjs
////go:generate brotli -fZk static/js/main.mjs

var (
	//go:embed dist/*
	Static        embed.FS
	StripedStatic fs.FS
)
