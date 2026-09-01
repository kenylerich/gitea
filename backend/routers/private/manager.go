// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package private

import (
	"net/http"

	"gitea.dev/backend/models/db"
	"gitea.dev/backend/modules/graceful"
	"gitea.dev/backend/modules/graceful/releasereopen"
	"gitea.dev/backend/modules/log"
	"gitea.dev/backend/modules/private"
	"gitea.dev/backend/modules/queue"
	"gitea.dev/backend/modules/templates"
	"gitea.dev/backend/modules/web"
	"gitea.dev/backend/services/context"
)

// ReloadTemplates reloads all the templates
func ReloadTemplates(ctx *context.PrivateContext) {
	err := templates.ReloadAllTemplates()
	if err != nil {
		ctx.PrivateInternalErrorf("Template error: %v", err)
		return
	}
	ctx.PlainText(http.StatusOK, "success")
}

// FlushQueues flushes all the Queues
func FlushQueues(ctx *context.PrivateContext) {
	opts := web.GetForm[*private.FlushOptions](ctx)
	if opts.NonBlocking {
		// Save the hammer ctx here - as a new one is created each time you call this.
		baseCtx := graceful.GetManager().HammerContext()
		go func() {
			err := queue.GetManager().FlushAll(baseCtx, opts.Timeout)
			if err != nil {
				log.Error("Flushing request timed-out with error: %v", err)
			}
		}()
		ctx.JSON(http.StatusAccepted, private.Response{
			UserMsg: "Flushing",
		})
		return
	}
	err := queue.GetManager().FlushAll(ctx, opts.Timeout)
	if err != nil {
		ctx.PrivateUserErrorf(http.StatusRequestTimeout, "%v", err)
		return
	}
	ctx.PlainText(http.StatusOK, "success")
}

// PauseLogging pauses logging
func PauseLogging(ctx *context.PrivateContext) {
	log.GetManager().PauseAll()
	ctx.PlainText(http.StatusOK, "success")
}

// ResumeLogging resumes logging
func ResumeLogging(ctx *context.PrivateContext) {
	log.GetManager().ResumeAll()
	ctx.PlainText(http.StatusOK, "success")
}

// ReleaseReopenLogging releases and reopens logging files
func ReleaseReopenLogging(ctx *context.PrivateContext) {
	if err := releasereopen.GetManager().ReleaseReopen(); err != nil {
		ctx.PrivateInternalErrorf("Error during release and reopen: %v", err)
		return
	}
	ctx.PlainText(http.StatusOK, "success")
}

// SetLogSQL re-sets database SQL logging
func SetLogSQL(ctx *context.PrivateContext) {
	db.SetLogSQL(ctx, ctx.FormBool("on"))
	ctx.PlainText(http.StatusOK, "success")
}
