// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package private

import (
	"errors"
	"net/http"

	issues_model "gitea.dev/backend/models/issues"
	user_model "gitea.dev/backend/models/user"
	"gitea.dev/backend/modules/git"
	"gitea.dev/backend/modules/private"
	"gitea.dev/backend/modules/web"
	"gitea.dev/backend/services/agit"
	gitea_context "gitea.dev/backend/services/context"
)

// HookProcReceive proc-receive hook - only handles agit Proc-Receive requests at present
func HookProcReceive(ctx *gitea_context.PrivateContext) {
	opts := web.GetForm[*private.HookOptions](ctx)
	if !git.DefaultFeatures().SupportProcReceive {
		ctx.Status(http.StatusNotFound)
		return
	}
	if !loadContextDoerPermission(ctx, opts.UserID, opts.UserExtDoerData) {
		return
	}

	results, err := agit.ProcReceive(ctx, ctx.Repo.Repository, ctx.Repo.GitRepo, &agit.ProcReceiveOptions{
		OldCommitIDs:   opts.OldCommitIDs,
		NewCommitIDs:   opts.NewCommitIDs,
		RefFullNames:   opts.RefFullNames,
		GitPushOptions: opts.GitPushOptions,
		Doer:           ctx.Doer,
	})
	if err != nil {
		if errors.Is(err, issues_model.ErrMustCollaborator) {
			ctx.PrivateUserErrorf(http.StatusUnauthorized, "You must be a collaborator to create pull request.")
		} else if errors.Is(err, user_model.ErrBlockedUser) {
			ctx.PrivateUserErrorf(http.StatusUnauthorized, "Cannot create pull request because you are blocked by the repository owner.")
		} else {
			ctx.PrivateInternalErrorf("agit.ProcReceive failed: %v", err)
		}

		return
	}

	ctx.JSON(http.StatusOK, private.HookProcReceiveResult{
		Results: results,
	})
}
