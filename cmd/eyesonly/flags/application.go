package flags

import (
	"time"

	"github.com/urfave/cli/v2"
)

func applicationFlags(expertMode bool) []cli.Flag {
	return []cli.Flag{
		// web render settings
		&cli.DurationFlag{
			Name:     "webrender-static-client-cache",
			Category: "Web Render",
			Value:    86400 * time.Second,
		},
		&cli.StringFlag{
			Name:     "webrender-static-cdn-host",
			Category: "Web Render",
			Value:    "cdn.anilibria.top",
		},

		// stats settings
		&cli.StringFlag{
			Name:     "stats-graphite-prefix",
			Category: "Stats Collection",
			Value:    "mh00net.golang.eyesonly",
		},
		&cli.StringFlag{
			Name:     "stats-graphite-proto",
			Category: "Stats Collection",
			Hidden:   expertMode,
			Value:    "udp4",
		},
		&cli.StringFlag{
			Name:     "stats-graphite-server",
			Category: "Stats Collection",
			Usage:    "127.0.0.1:2003",
			EnvVars:  []string{"GRAPHITE_SERVER"},
			Value:    "",
		},
		&cli.DurationFlag{
			Name:     "stats-metrics-interval",
			Category: "Stats Collection",
			Usage:    "must be a multiple of 5 seconds",
			Hidden:   expertMode,
			Value:    5 * time.Second,
		},
		&cli.DurationFlag{
			Name:     "stats-metrics-loop-warn",
			Category: "Stats Collection",
			Usage:    "warning for too long loop metrics collection",
			Hidden:   expertMode,
			Value:    100 * time.Millisecond,
		},
		&cli.IntFlag{
			Name:     "stats-queue-buffer-size",
			Category: "Stats Collection",
			Hidden:   expertMode,
			Value:    1 << 8, // 256
		},

		// analytics
		&cli.StringFlag{
			Name:     "analytics-database-path",
			Category: "Analytics",
			Value:    "data/analytics.db",
		},
		&cli.StringFlag{
			Name:     "analytics-database-table-stats",
			Category: "Analytics",
			Hidden:   expertMode,
			Value:    "stats",
		},
		&cli.BoolFlag{
			Name:     "analytics-database-reinit",
			Category: "Analytics",
			Usage:    "WARNING! It destroys all your data! Add analytics-database-i-am-sure to complete reinit",
			Hidden:   expertMode,
		},
		&cli.BoolFlag{
			Name:     "analytics-database-i-am-sure",
			Category: "Analytics",
			Hidden:   expertMode,
		},
		&cli.StringFlag{
			Name:     "analytics-metric-filter",
			Category: "Analytics",
			Usage:    "must be a valid regexp",
			Hidden:   expertMode,
			Value:    "",
		},
		&cli.DurationFlag{
			Name:     "analytics-analyze-interval",
			Category: "Analytics",
			Usage:    "must be a multiple of 5 seconds",
			Hidden:   expertMode,
			Value:    10 * time.Second,
		},
		&cli.DurationFlag{
			Name:     "analytics-analyze-window",
			Category: "Analytics",
			Usage:    "a period of time from which analyzer will be calculate differences ratio (from NOW()-`DURATION` TO NOW())",
			Hidden:   expertMode,
			Value:    2 * time.Minute,
		},
		&cli.BoolFlag{
			Name:     "analytics-lockdown-force",
			Category: "Analytics",
			Usage:    "lockdown mode force ignoring all unlockdown triggers",
			Hidden:   expertMode,
		},
		&cli.DurationFlag{
			Name:     "analytics-lockdown-time",
			Category: "Analytics",
			Usage:    "a period of time when JS challenge will be forced for all requests after triggering",
			Hidden:   expertMode,
			Value:    10 * time.Minute,
		},
		&cli.DurationFlag{
			Name:     "analytics-trigger-boothold",
			Category: "Analytics",
			Usage:    "wait `DURATION` after service initialization before starting the analyzer",
			Hidden:   expertMode,
			Value:    5 * time.Minute,
		},
		&cli.Float64Flag{
			Name:     "analytics-trigger-ratio-warn",
			Category: "Analytics",
			Usage:    "set warn in range of `warn value` and `error value` (trigger-ratio-err)",
			Hidden:   expertMode,
			Value:    2,
		},
		&cli.Float64Flag{
			Name:     "analytics-trigger-ratio-err",
			Category: "Analytics",
			Usage:    "set err when ratio is above `error value`",
			Hidden:   expertMode,
			Value:    3,
		},

		// todo : query timeouts
		// &cli.DurationFlag{
		// 	Name:     "analytics-database-query-timeout",
		// 	Category: "Analytics",
		// 	Hidden:   expertMode,
		// 	Value:    500 * time.Millisecond,
		// },

		// todo challenger settings
		// &cli.StringFlag{
		// 	Name:     "challenger-ja3service-url",
		// 	Category: "Challenger",
		// 	Value:    "https://anilibria.top/internal/justanimesan",
		// },
	}
}
