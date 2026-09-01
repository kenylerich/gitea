// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"testing"

	"gitea.dev/backend/models/unittest"
	"gitea.dev/backend/modules/hostmatcher"
	"gitea.dev/backend/modules/setting"

	_ "gitea.dev/backend/models"
	_ "gitea.dev/backend/models/actions"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		SetUp: func() error {
			// for tests, allow only loopback IPs. This must run after the test config is loaded (which
			// resets the shared Security.AllowedHostList) and before Init() builds the delivery client.
			setting.Security.AllowedHostList = hostmatcher.MatchBuiltinLoopback
			setting.LoadQueueSettings()
			return Init()
		},
	})
}
