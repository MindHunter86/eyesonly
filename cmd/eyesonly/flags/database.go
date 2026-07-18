package flags

import (
	"time"

	"github.com/urfave/cli/v2"
)

func databaseFlags(expertMode bool) []cli.Flag {
	return []cli.Flag{

		// database settings
		&cli.StringFlag{
			Name:     "database-path",
			Category: "Database",
			Value:    "data/eyesonly.db",
		},
		&cli.DurationFlag{
			Name:     "database-open-timeout",
			Category: "Database",
			Hidden:   expertMode,
			Value:    2 * time.Second,
		},
		&cli.BoolFlag{
			Name:               "database-no-freelist-sync",
			Category:           "Database",
			Usage:              "This improves the database write performance under normal operation, but requires a full database re-sync during recovery.",
			Hidden:             expertMode,
			DisableDefaultText: true,
		},
		&cli.BoolFlag{
			Name:               "database-no-sync",
			Category:           "Database",
			Hidden:             expertMode,
			DisableDefaultText: true,
		},
	}
}
