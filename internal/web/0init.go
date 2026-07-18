package web

import (
	"context"
	"io/fs"
	"math"
	"strconv"
	"sync"

	"github.com/MindHunter86/eyesonly/internal/utils"
	"github.com/urfave/cli/v2"
)

var (
	DataDivPostfix string

	dynamicMu        sync.RWMutex
	DynamicTemplates Templates
)

func init() {
	utils.RegisterSubservice(utils.CtxWebRender, func(c context.Context) (cfg any, e error) {
		cli := utils.ContextValueExtract[*cli.Context](c, utils.CtxCliContext)
		cfg = &Config{
			AppName:    cli.App.Name,
			AppVersion: cli.App.Version,

			CDNDomain: cli.String("webrender-static-cdn-host"),

			ClientCacheDur: cli.Duration("webrender-static-client-cache"),
			ClientCacheSeconds: strconv.Itoa(int(math.Floor(
				cli.Duration("webrender-static-client-cache").Seconds()))),
		}

		// new templates subsystem initialization
		DynamicTemplates = make(Templates)

		// strip embed fs for fasthttp.FS
		if StripedStatic, e = fs.Sub(Static, "dist"); e != nil {
			panic(e)
		}

		// load randomized div token for payload in html
		var buf []byte
		if buf, e = openFile("dist/payload-token.txt"); e != nil {
			panic(utils.ExtraErrorWrapper(e, "try to open payload-token.txt file"))
		}

		if len(buf) != 0 {
			DataDivPostfix = utils.UnsafeString(buf)
			return
		}

		// todo : research its importance
		// l := utils.ContextValueExtract[*zerolog.Logger](c, utils.CtxZeroLogger)
		// l.Error().Msg("found empty payload-token file, token wasn't configured")
		// l.Info().Msg("replace random div token with 00000000")
		// DataDivPostfix = "00000000"
		// return

		panic("no random token found in payload-token.txt")
	})
}
