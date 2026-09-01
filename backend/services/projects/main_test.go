// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package project

import (
	"testing"

	"gitea.dev/backend/models/unittest"

	_ "gitea.dev/backend/models/actions"
	_ "gitea.dev/backend/models/activities"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
