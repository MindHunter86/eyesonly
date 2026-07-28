package stats

import (
	"sync"
)

type IncrementMetric uint16
type AverageMetric uint16
type PersistentMetric uint16

const (
	// increment metrics
	IMHTTPServerRequest IncrementMetric = iota
	IMHTTPServerInvalidRequest
	IMHTTPServerPanic
	IMHTTPServerCode100
	IMHTTPServerCode200
	IMHTTPServerCode300
	IMHTTPServerCode400
	IMHTTPServerCode500
	IMHTTPServerNoCode

	IMFiberReqPathStatic
	IMFiberReqPathAssets
	IMFiberReqPathAuthMod
	IMFiberReqPathVerify
	IMFiberReqPathRoot
	IMFiberReqPathSettings
	IMFiberReqPathConfigPatch

	IMStatsLoopCount
	IMStatsLoopLockedError

	IMSvcKernSignDebug

	// average metrics
	AMHTTPServerLatency AverageMetric = iota

	// persistent metrics
	// PMSvcModulesRunning PersistentMetric = iota

	_metricsEOL
)

var MetricNames = map[any]string{
	IMHTTPServerRequest:        "fiber.requests",
	IMHTTPServerInvalidRequest: "fiber.requests_invalid",
	IMHTTPServerPanic:          "fiber.requests_panics",
	IMHTTPServerCode100:        "fiber.requests_100",
	IMHTTPServerCode200:        "fiber.requests_200",
	IMHTTPServerCode300:        "fiber.requests_300",
	IMHTTPServerCode400:        "fiber.requests_400",
	IMHTTPServerCode500:        "fiber.requests_500",
	IMHTTPServerNoCode:         "fiber.requests_000",
	IMFiberReqPathStatic:       "fiber.path.static",
	IMFiberReqPathAssets:       "fiber.path.assets",
	IMFiberReqPathAuthMod:      "fiber.path.authmod",
	IMFiberReqPathVerify:       "fiber.path.verify",
	IMFiberReqPathRoot:         "fiber.path.root",
	IMFiberReqPathSettings:     "fiber.path.settings",
	IMFiberReqPathConfigPatch:  "fiber.path.config_patch",
	IMStatsLoopCount:           "stats.loop_count",
	IMStatsLoopLockedError:     "stats.loop_lockederr",
	IMSvcKernSignDebug:         "service.kernsign_debug_calls",

	AMHTTPServerLatency: "fiber.latency_ns",
}

var queueTaskPool = sync.Pool{New: func() any {
	return &queueTask{}
}}

func acquireQueueTask() *queueTask  { return queueTaskPool.Get().(*queueTask) }
func releaseQueueTask(v *queueTask) { queueTaskPool.Put(v) }
