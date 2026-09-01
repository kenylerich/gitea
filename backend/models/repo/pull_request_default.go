// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"context"

	"gitea.dev/backend/models/unit"
	"gitea.dev/backend/modules/util"
)

func (repo *Repository) GetPullRequestTargetBranch(ctx context.Context) string {
	unitPRConfig := repo.MustGetUnit(ctx, unit.TypePullRequests).PullRequestsConfig()
	return util.IfZero(unitPRConfig.DefaultTargetBranch, repo.DefaultBranch)
}
