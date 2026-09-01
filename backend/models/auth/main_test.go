// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package auth_test

import (
	"testing"

	"gitea.dev/backend/models/unittest"

	_ "gitea.dev/backend/models"
	_ "gitea.dev/backend/models/actions"
	_ "gitea.dev/backend/models/activities"
	_ "gitea.dev/backend/models/auth"
	_ "gitea.dev/backend/models/perm/access"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
