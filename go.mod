module github.com/MindHunter86/eyesonly

// I know, it's too legacy for 2026, sorry
go 1.22.12

// godebug (
// 	gctrace=1

// 	http2client=0
// 	http2server=0

// 	tlsmaxrsasize=4096
// 	netedns0=0
// )

exclude (
	github.com/valyala/fasthttp v1.60.0
	github.com/valyala/fasthttp v1.61.0
	github.com/valyala/fasthttp v1.62.0
	github.com/valyala/fasthttp v1.63.0
	github.com/valyala/fasthttp v1.64.0
	github.com/valyala/fasthttp v1.65.0
	github.com/valyala/fasthttp v1.66.0
	github.com/valyala/fasthttp v1.67.0
	github.com/valyala/fasthttp v1.68.0
	github.com/valyala/fasthttp v1.69.0
	github.com/valyala/fasthttp v1.70.0
	github.com/valyala/fasthttp v1.71.0
)

exclude (
	golang.org/x/sys v0.31.0
	golang.org/x/sys v0.32.0
	golang.org/x/sys v0.33.0
	golang.org/x/sys v0.34.0
	golang.org/x/sys v0.35.0
	golang.org/x/sys v0.36.0
	golang.org/x/sys v0.37.0
)

exclude (
	go.etcd.io/bbolt v1.3.12
	go.etcd.io/bbolt v1.4.0
	go.etcd.io/bbolt v1.4.1
	go.etcd.io/bbolt v1.4.2
	go.etcd.io/bbolt v1.4.3
)

exclude (
	golang.org/x/tools v0.30.0
	golang.org/x/tools v0.31.0
	golang.org/x/tools v0.32.0
	golang.org/x/tools v0.33.0
	golang.org/x/tools v0.34.0
	golang.org/x/tools v0.35.0
	golang.org/x/tools v0.36.0
	golang.org/x/tools v0.37.0
	golang.org/x/tools v0.38.0
)

exclude (
	github.com/apache/arrow-go/v18 v18.2.0
	github.com/apache/arrow-go/v18 v18.3.0
	github.com/apache/arrow-go/v18 v18.3.1
	github.com/apache/arrow-go/v18 v18.4.0
	github.com/apache/arrow-go/v18 v18.4.1
)

exclude (
	golang.org/x/mod v0.23.0
	golang.org/x/mod v0.24.0
	golang.org/x/mod v0.25.0
	golang.org/x/mod v0.26.0
	golang.org/x/mod v0.27.0
	golang.org/x/mod v0.28.0
	golang.org/x/mod v0.29.0
)

exclude (
	golang.org/x/sync v0.11.0
	golang.org/x/sync v0.12.0
	golang.org/x/sync v0.13.0
	golang.org/x/sync v0.14.0
	golang.org/x/sync v0.15.0
	golang.org/x/sync v0.16.0
	golang.org/x/sync v0.17.0
)

exclude (
	golang.org/x/exp v0.0.0-20251009144603-d2f985daa21b
	golang.org/x/exp v0.0.0-20251125195548-87e1e737ad39
)

exclude github.com/klauspost/compress v1.18.1

exclude (
	github.com/rs/zerolog v1.35.0
	github.com/rs/zerolog v1.35.1
)

require (
	github.com/go-kit/kit v0.13.0
	github.com/go-sql-driver/mysql v1.8.1
	github.com/gofiber/fiber/v2 v2.52.13
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.34.0
	github.com/urfave/cli/v2 v2.27.7
	github.com/valyala/bytebufferpool v1.0.0
	github.com/valyala/tcplisten v1.0.0
	go.uber.org/atomic v1.11.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/stretchr/testify v1.11.0 // indirect
	github.com/tinylib/msgp v1.2.5 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/sys v0.30.0 // indirect
)
