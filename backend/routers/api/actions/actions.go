// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"net/http"

	"gitea.dev/backend/modules/web"
	"gitea.dev/backend/routers/api/actions/ping"
	"gitea.dev/backend/routers/api/actions/runner"
)

func Routes(prefix string) *web.Router {
	m := web.NewRouter()

	path, handler := ping.NewPingServiceHandler()
	m.Post(path+"*", http.StripPrefix(prefix, handler).ServeHTTP)

	path, handler = runner.NewRunnerServiceHandler()
	m.Post(path+"*", http.StripPrefix(prefix, handler).ServeHTTP)

	return m
}
