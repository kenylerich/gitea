// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package misc

import (
	api "gitea.dev/backend/modules/structs"
	"gitea.dev/backend/modules/util"
	"gitea.dev/backend/modules/web"
	"gitea.dev/backend/routers/common"
	"gitea.dev/backend/services/context"
)

// Markup render markup document to HTML
func Markup(ctx *context.Context) {
	form := web.GetForm[*api.MarkupOption](ctx)
	mode := util.Iif(form.Wiki, "wiki", form.Mode) //nolint:staticcheck // form.Wiki is deprecated
	common.RenderMarkup(ctx.Base, ctx.Repo, mode, form.Text, form.Context, form.FilePath)
}
