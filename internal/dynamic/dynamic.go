package dynamic

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/MindHunter86/eyesonly/internal/utils"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"go.uber.org/atomic"
)

var dynamicFlagsMu sync.RWMutex
var dynamicFlags = map[string]*atomic.String{}

func String(c context.Context, k string) string {
	dynamicFlagsMu.RLock()

	if fl, ok := dynamicFlags[k]; ok {
		defer dynamicFlagsMu.RUnlock()
		return fl.Load()
	}

	dynamicFlagsMu.RUnlock()
	dynamicFlagsMu.Lock()
	defer dynamicFlagsMu.Unlock()

	cl := utils.ContextValueExtract[*cli.Context](c, utils.CtxCliContext)
	dynamicFlags[k] = atomic.NewString(cl.String(k))
	return cl.String(k)
}

func Duration(c context.Context, k string) (dur time.Duration) {
	dynamicFlagsMu.RLock()

	if fl, ok := dynamicFlags[k]; ok {
		defer dynamicFlagsMu.RUnlock()

		var e error
		if dur, e = time.ParseDuration(fl.Load()); e != nil {
			log := utils.ContextValueExtract[*zerolog.Logger](c, utils.CtxZeroLogger)
			log.Info().Msg(utils.ExtraErrorWrapper(e, "trying parse dur %s with value %s", k, fl.Load()).Error())
			return 0
		}

		return dur
	}

	dynamicFlagsMu.RUnlock()
	dynamicFlagsMu.Lock()
	defer dynamicFlagsMu.Unlock()

	cl := utils.ContextValueExtract[*cli.Context](c, utils.CtxCliContext)
	dynamicFlags[k] = atomic.NewString(cl.Duration(k).String())
	return cl.Duration(k)
}

func Bool(c context.Context, k string) bool {
	dynamicFlagsMu.RLock()

	if fl, ok := dynamicFlags[k]; ok {
		defer dynamicFlagsMu.RUnlock()
		return fl.Load() == "true"
	}

	dynamicFlagsMu.RUnlock()
	dynamicFlagsMu.Lock()
	defer dynamicFlagsMu.Unlock()

	cl := utils.ContextValueExtract[*cli.Context](c, utils.CtxCliContext)
	dynamicFlags[k] = atomic.NewString("false")
	if cl.Bool(k) {
		dynamicFlags[k].Store("true")
	}

	return cl.Bool(k)
}

func Int(c context.Context, k string) (num int) {
	dynamicFlagsMu.RLock()

	if fl, ok := dynamicFlags[k]; ok {
		defer dynamicFlagsMu.RUnlock()

		var e error
		if num, e = strconv.Atoi(fl.Load()); e != nil {
			log := utils.ContextValueExtract[*zerolog.Logger](c, utils.CtxZeroLogger)
			log.Info().Msg(utils.ExtraErrorWrapper(e, "trying parse int %s with value %s", k, fl.Load()).Error())
			return 0
		}

		return num
	}

	dynamicFlagsMu.RUnlock()
	dynamicFlagsMu.Lock()
	defer dynamicFlagsMu.Unlock()

	cl := utils.ContextValueExtract[*cli.Context](c, utils.CtxCliContext)
	dynamicFlags[k] = atomic.NewString(strconv.Itoa(cl.Int(k)))
	return cl.Int(k)
}

func Float64(c context.Context, k string) (num float64) {
	dynamicFlagsMu.RLock()

	if fl, ok := dynamicFlags[k]; ok {
		defer dynamicFlagsMu.RUnlock()

		var e error
		if num, e = strconv.ParseFloat(fl.Load(), 64); e != nil {
			log := utils.ContextValueExtract[*zerolog.Logger](c, utils.CtxZeroLogger)
			log.Info().Msg(utils.ExtraErrorWrapper(e, "trying parse fl64 %s with value %s", k, fl.Load()).Error())
			return 0
		}

		return num
	}

	dynamicFlagsMu.RUnlock()
	dynamicFlagsMu.Lock()
	defer dynamicFlagsMu.Unlock()

	cl := utils.ContextValueExtract[*cli.Context](c, utils.CtxCliContext)
	dynamicFlags[k] = atomic.NewString(strconv.FormatFloat(cl.Float64(k), 'f', 6, 64))
	return cl.Float64(k)
}

func PatchValue(c context.Context, k, v string) error {
	dynamicFlagsMu.RLock()
	defer dynamicFlagsMu.RUnlock()

	var ok bool
	var fl *atomic.String
	if fl, ok = dynamicFlags[GetCliName(k)]; !ok {
		return utils.ExtraErrorWrapper(
			errors.New("requested key wasn't defined; new config values creation isn't possible"),
			"lookup for typped value of %s in cli.Context", k)
	}

	ov := fl.Swap(v)
	log := utils.ContextValueExtract[*zerolog.Logger](c, utils.CtxZeroLogger)
	log.Info().Msgf("configuration value has been changed via API; %s was %s become %s", k, ov, v)

	forEachCallback(c)
	return nil
}

func FetchValues() (vals []string) {
	return availableFlags
}

func GetCliName(k string) string {
	return whitelistedFlags[k]
}

func IsFlagRO(k string) bool {
	return whitelistedFlags[k][0] == '_'
}
