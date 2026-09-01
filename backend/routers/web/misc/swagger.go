// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package misc

import (
	"net/http"

	"gitea.dev/backend/services/context"
)

func Swagger(ctx *context.Context) {
	ctx.HTML(http.StatusOK, "swagger/openapi-viewer")
}
