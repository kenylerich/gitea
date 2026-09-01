// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package system_test

import (
	"testing"

	"gitea.dev/backend/models/unittest"

	_ "gitea.dev/backend/models" // register models
	_ "gitea.dev/backend/models/actions"
	_ "gitea.dev/backend/models/activities"
	_ "gitea.dev/backend/models/system" // register models of system
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
