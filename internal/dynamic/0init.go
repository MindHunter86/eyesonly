package dynamic

import (
	"context"
	"sort"
	"strings"

	"github.com/MindHunter86/eyesonly/internal/utils"
)

func init() {
	utils.RegisterSubservice(utils.CtxDynamic,
		func(context.Context) (any, error) {
			for k := range whitelistedFlags {
				availableFlags = append(availableFlags, k)
			}

			sort.Slice(availableFlags, func(i, j int) bool {
				return strings.ToLower(availableFlags[i]) < strings.ToLower(availableFlags[j])
			})

			return nil, nil
		})
}
