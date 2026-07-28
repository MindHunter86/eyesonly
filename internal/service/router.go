package service

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/MindHunter86/eyesonly/internal/stats"
	"github.com/MindHunter86/eyesonly/internal/utils"
	"github.com/MindHunter86/eyesonly/internal/web"
	"github.com/MindHunter86/eyesonly/internal/webapp/platform/crypto"
	"github.com/MindHunter86/eyesonly/internal/webapp/secrets"
	"github.com/MindHunter86/eyesonly/internal/webapp/sessions"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/valyala/tcplisten"

	db "github.com/MindHunter86/eyesonly/internal/database"
)

func (m *Service) fiberMiddlewareInitialization() {
	// pprof profiler
	// manual:
	// 	curl -o profile.out https://host/debug/pprof -H 'X-Authorization: $TOKEN'
	// 	go tool pprof profile.out
	ppEnabled, ppSecret :=
		gCli.Bool("http-pprof-enable"),
		utils.UnsafeBytes(gCli.String("http-pprof-secret"))

	m.fb.Use(pprof.New(pprof.Config{
		Next: func(c *fiber.Ctx) bool {
			if !ppEnabled {
				return true
			}

			return !bytes.Equal(ppSecret, c.Request().Header.Peek("x-pprof-secret"))
		},
		Prefix: gCli.String("http-pprof-prefix"),
	}))

	// Small middleware for optimizing working with IP
	// and reducing allocs by using net.IP package
	m.fb.Use(func(c *fiber.Ctx) (e error) {
		_, e = utils.IPFromFiberRequest(c), c.Next()

		c.Response().Header.Del(utils.IP_RESPONSE_HEADER)
		return e
	})

	// simple limiter for all requests
	limitMaxRps := gCli.Int("limit-request-maxrps")
	m.fb.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return limitMaxRps == 0 || utils.IPFromFiberRequest(c) == "127.0.0.1"
		},
		KeyGenerator: utils.IPFromFiberRequest,

		Max:        gCli.Int("limit-request-maxrps"),
		Expiration: gCli.Duration("limit-request-expiration"),
	}))

	// push global app context for request
	m.fb.Use(func(c *fiber.Ctx) (e error) {
		c.SetUserContext(gCtx)
		return c.Next()
	})

	// panic recover for all handlers
	m.fb.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			utils.Rlog(c, zerolog.ErrorLevel).Str("request", c.Request().String()).Bytes("stack", debug.Stack()).
				Msg("panic has been caught")
			_, _ = fmt.Fprintf(os.Stderr, "panic: %v\n%s\n", e, debug.Stack()) //nolint:errcheck // This will never fail

			c.Status(fiber.StatusInternalServerError)

			// stat paniced request
			ctx := c.UserContext()
			sts := utils.ContextValueExtract[*stats.Stats](ctx, utils.CtxStats)

			sts.WriteIncMetric(stats.IMHTTPServerPanic)
		},
	}))

	// time collector + logger + stats
	m.fb.Use(func(c *fiber.Ctx) (e error) {
		started, e := time.Now(), c.Next()

		status, lvl := c.Response().StatusCode(), utils.HTTPAccessLogLevel

		// defaults for errored requests
		var cause string
		if e != nil {
			cause, lvl = e.Error(), zerolog.ErrorLevel
			status = fiber.StatusInternalServerError
		}

		// redefine variables for Fiber errors
		switch e := e.(type) { // skipcq: CRT-A0014 type should be inside switch
		case *fiber.Error:
			if e.Code > fiber.StatusBadRequest && e.Code < fiber.StatusInternalServerError {
				cause, lvl = e.Error(), zerolog.WarnLevel
			} else if e.Code > fiber.StatusInternalServerError {
				cause, lvl = e.Error(), zerolog.ErrorLevel
			}

			status = e.Code
		}

		// dump request data for debugging "error" cases
		if lvl >= zerolog.ErrorLevel && cause != "" {
			utils.Rlog(c, lvl).Msg(c.Request().String())
		}

		// utils.Rlog(c, lvl).
		// 	Int("status", status).
		// 	Str("method", c.Method()).
		// 	Str("path", c.Path()).
		// 	Str("ip", utils.IPFromFiberRequest(c)).
		// 	Dur("latency", elapsed).
		// 	Str("user-agent", c.Get(fiber.HeaderUserAgent)).Msg(cause)

		// todo : need some tests in production, revert if causes errs
		elapsed := time.Since(started)
		utils.RlogFast(c, lvl, status, elapsed, cause)

		// stats record
		ctx := c.UserContext()
		sts := utils.ContextValueExtract[*stats.Stats](ctx, utils.CtxStats)

		if sts != nil {
			sts.WriteIncMetric(stats.IMHTTPServerRequest)
			sts.WriteAvgMetric(stats.AMHTTPServerLatency, uint64(elapsed.Nanoseconds()))

			// sts.QueueWriteMetric(stats.IMHTTPServerRequest)
			if status >= 100 && status <= 199 {
				sts.WriteIncMetric(stats.IMHTTPServerCode100)
			} else if status >= 200 && status <= 299 {
				// sts.QueueWriteMetric(stats.IMHTTPServerCode200)
				sts.WriteIncMetric(stats.IMHTTPServerCode200)
			} else if status >= 300 && status <= 399 {
				// sts.QueueWriteMetric(stats.IMHTTPServerCode300)
				sts.WriteIncMetric(stats.IMHTTPServerCode300)
			} else if status >= 400 && status <= 499 {
				// sts.QueueWriteMetric(stats.IMHTTPServerCode400)
				sts.WriteIncMetric(stats.IMHTTPServerCode400)
			} else if status >= 500 && status <= 599 {
				// sts.QueueWriteMetric(stats.IMHTTPServerCode500)
				sts.WriteIncMetric(stats.IMHTTPServerCode500)
			} else {
				sts.WriteIncMetric(stats.IMHTTPServerNoCode)
			}
		}

		return
	})
}

func (m *Service) fiberRouterInitialization() (e error) {
	//
	//	Router pre-initialization
	//

	// dynamic settings and helpers:
	// statsToken := utils.UnsafeBytes(gCli.String("http-stats-secret"))

	//
	//	Router handlers configuration
	//

	// routing base
	// root := m.fb.Group("/.within.website/x/cmd/" + m.fb.Config().AppName)

	// expvars stats page
	// internal := root.Group("/internal")
	// internal.Get("/stats", func(c *fiber.Ctx) error {
	// 	if !bytes.Equal(c.Request().Header.Peek(fasthttp.HeaderAuthorization), statsToken) {
	// 		c.Status(fiber.StatusNotFound)
	// 		return fiber404ErrorHandler(c)
	// 	}

	// 	expvarhandler.ExpvarHandler(c.Context())
	// 	return nil
	// })

	//
	// APIv1 prepare
	//

	var database *sql.DB
	if database, e = db.Open(gCtx); e != nil {
		return
	}
	// defer database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if e = db.Migrate(ctx, database, gCli.String("database-driver")); e != nil {
		return
	}

	var cryptosvc *crypto.Service
	if cryptosvc, e = crypto.NewService(gCli.String("secrets-encryption-key")); e != nil {
		return
	}

	jwtsecret := []byte(gCli.String("sessions-jwt-secret"))
	sessionsvc := sessions.NewService(jwtsecret, gCli.Duration("sessions-jwt-ttl"))

	secretrepo := secrets.NewRepository(database)
	secretsvc := secrets.NewService(secretrepo, cryptosvc,
		strings.TrimRight(gCli.String("app-baseurl"), "/"))

	sessionhand := sessions.NewHandler(sessionsvc)
	secrethand := secrets.NewHandler(gCtx, sessionsvc, secretsvc)

	//
	// APIv1 initialization
	//

	apiv1 := m.fb.Group("/v1")
	sessionhand.Register(apiv1)
	secrethand.Register(apiv1)

	// add healthz and readyz
	// app.Get("/healthz", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })
	// app.Get("/readyz", func(c *fiber.Ctx) error {
	// 	if e := database.Ping(); e != nil {
	// 		return err
	// 	}
	// 	return c.JSON(fiber.Map{"status": "ready"})
	// })

	// index page
	if e = web.RegisterStatic(m.fb); e != nil {
		return
	}

	// ! should be at the end of handlers list !
	// custom 404 handler for fiber.Error allocs reduce
	m.fb.Use(fiber404ErrorHandler)
	return
}

func fiberErrorHandler(c *fiber.Ctx, err error) (_ error) {
	// reject invalid requests
	if strings.TrimSpace(c.Hostname()) == "" {
		gLog.Warn().Msgf("invalid request from %s: %+v ; error - %+v",
			utils.IPFromFiberRequest(c), c, err)

		sts := utils.ContextValueExtract[*stats.Stats](gCtx, utils.CtxStats)
		if sts != nil {
			sts.WriteIncMetric(stats.IMHTTPServerRequest)
			sts.WriteIncMetric(stats.IMHTTPServerInvalidRequest)
		}

		return c.Context().Conn().Close()
	}

	// disable caching for error responder
	c.Set(fiber.HeaderCacheControl, "no-cache")

	// JSON error content-type:
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	// * using errors.As is more readable way, yes
	// * but we building anti-ddos solution, so we need allocs minimization
	// * for responding on invalid requests
	switch err := err.(type) {
	case *fiber.Error:
		writeJsonErrorFastTo(c, err.Code, err.Message)
		c.Status(err.Code)

		utils.ReleaseFiberError(err)
	default:
		writeJsonErrorFastTo(c, fiber.StatusInternalServerError, err.Error())
		c.Status(fiber.StatusInternalServerError)
	}

	// TODO : fixme; I think it can be dropped
	if zerolog.GlobalLevel() <= zerolog.DebugLevel {
		utils.Rlog(c, zerolog.DebugLevel).Msgf("%+v", err)
	}

	return
}

func (m *Service) fhttpListenerInitialization() func() error {
	reuseport, deferaccept, fastopen, backlog :=
		gCli.Bool("http-adv-reuseport"),
		gCli.Bool("http-adv-deferaccept"),
		gCli.Bool("http-adv-tcpfastopen"),
		gCli.Int("http-adv-backlog")

	if !reuseport && !deferaccept && !fastopen && backlog == 0 {
		return nil
	}

	tcpopts := &tcplisten.Config{
		ReusePort:   reuseport,
		DeferAccept: deferaccept,
		FastOpen:    fastopen,
		Backlog:     backlog,
	}

	ln, e := tcpopts.NewListener("tcp4", gCli.String("http-listen-addr"))
	if e != nil {
		gLog.Error().Msg("could not initialize custom net.Listener, due to - " + e.Error())
		return nil
	}

	return func() error {
		return m.fb.Listener(ln)
	}
}
