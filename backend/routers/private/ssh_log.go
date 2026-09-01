// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package private

import (
	"net/http"

	"gitea.dev/backend/modules/log"
	"gitea.dev/backend/modules/private"
	"gitea.dev/backend/modules/setting"
	"gitea.dev/backend/modules/web"
	"gitea.dev/backend/services/context"
)

// SSHLog hook to response ssh log
func SSHLog(ctx *context.PrivateContext) {
	if !setting.Log.EnableSSHLog {
		ctx.Status(http.StatusOK)
		return
	}

	opts := web.GetForm[*private.SSHLogOption](ctx)

	if opts.IsError {
		log.Error("ssh: %v", opts.Message)
		ctx.Status(http.StatusOK)
		return
	}

	log.Debug("ssh: %v", opts.Message)
	ctx.Status(http.StatusOK)
}
