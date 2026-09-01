// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package explore

import (
	"net/http"

	"gitea.dev/backend/models/db"
	repo_model "gitea.dev/backend/models/repo"
	api "gitea.dev/backend/modules/structs"
	"gitea.dev/backend/services/context"
	"gitea.dev/backend/services/convert"
)

// TopicSearch search for creating topic
func TopicSearch(ctx *context.Context) {
	opts := &repo_model.FindTopicOptions{
		Keyword: ctx.FormString("q"),
		ListOptions: db.ListOptions{
			Page:     ctx.FormInt("page"),
			PageSize: convert.ToCorrectPageSize(ctx.FormInt("limit")),
		},
	}

	topics, total, err := db.FindAndCount[repo_model.Topic](ctx, opts)
	if err != nil {
		ctx.HTTPError(http.StatusInternalServerError)
		return
	}

	topicResponses := make([]*api.TopicResponse, len(topics))
	for i, topic := range topics {
		topicResponses[i] = convert.ToTopicResponse(topic)
	}

	ctx.SetTotalCountHeader(total)
	ctx.JSON(http.StatusOK, map[string]any{
		"topics": topicResponses,
	})
}
