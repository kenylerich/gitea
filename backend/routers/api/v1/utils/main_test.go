// Copyright 2018 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package utils

import (
	"testing"

	"gitea.dev/backend/models/unittest"
	"gitea.dev/backend/modules/setting"
	webhook_service "gitea.dev/backend/services/webhook"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		SetUp: func() error {
			setting.LoadQueueSettings()
			return webhook_service.Init()
		},
	})
}
