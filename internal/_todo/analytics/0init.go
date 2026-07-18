package analytics

import (
	"context"

	"github.com/MindHunter86/eyesonly/internal/stats"
	"github.com/MindHunter86/eyesonly/internal/utils"
)

func init() {
	utils.RegisterSubservice(utils.CtxAnalytics,
		func(c context.Context) (any, error) {
			a, e := NewAnalytics(c)
			if e != nil {
				return a, e
			}

			// callbacks
			utils.RegisterCallback(utils.OnServiceBootstrap, func(ctx context.Context) error {
				a.sts = utils.ContextValueExtract[*stats.Stats](ctx, utils.CtxStats)
				a.sts.RegisterWriter(a)
				a.rdy.Store(true)
				return nil
			})
			utils.RegisterCallback(utils.OnServiceTicker1sec, a.onServiceTicker1sec)
			utils.RegisterCallback(utils.OnServiceDestruct, func(context.Context) error {
				return a.db.Close()
			})

			return a, nil
		})
}
