package dynamic

import (
	"bytes"

	"github.com/MindHunter86/eyesonly/internal/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	patchRequestKey   = "key"
	patchRequestValue = "value"
)

func ConfigPatchHandler(c *fiber.Ctx) (e error) {
	k := bytes.TrimSpace(c.Context().QueryArgs().Peek(patchRequestKey))
	v := bytes.TrimSpace(c.Context().QueryArgs().Peek(patchRequestValue))

	// !! todo : ADD CSRF!!

	// !! todo : 2DELETE
	c.Set("Access-Control-Allow-Origin", "*")

	if len(k) == 0 || len(v) == 0 {
		return utils.AcquireFiberError(fiber.StatusBadRequest, "key or value couldn't be empty")
	}

	if wlv, ok := whitelistedFlags[utils.UnsafeString(k)]; !ok {
		return utils.AcquireFiberError(fiber.StatusBadRequest, "given key is not found in runtime")
	} else if wlv[0] == '_' {
		return utils.AcquireFiberError(fiber.StatusBadRequest, "given key is readonly and couldn't be changed")
	}

	val := utils.CopyBytes(v)

	ctx := c.UserContext()
	if e = PatchValue(ctx, utils.UnsafeString(k), utils.UnsafeString(val)); e != nil {
		return utils.AcquireFiberError(fiber.StatusInternalServerError, e.Error())
	}

	return c.SendStatus(201)
}
