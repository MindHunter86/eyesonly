package dynamic

var whitelistedFlags = map[string]string{
	"AnalyticsAnalyzeInterval":  "analytics-analyze-interval",
	"AnalyticsAnalyzeWindow":    "analytics-analyze-window",
	"AnalyticsLockdownDuration": "analytics-lockdown-time",
	"AnalyticsLockdown":         "analytics-lockdown-force",
	"AnalyticsTrigBootHold":     "analytics-trigger-boothold",
	"AnalyticsTrigRatioWrn":     "analytics-trigger-ratio-warn",
	"AnalyticsTrigRatioErr":     "analytics-trigger-ratio-err",

	// RO example
	// "LogLevel":                  "_log-level",
}

var availableFlags = make([]string, 0, len(whitelistedFlags))
