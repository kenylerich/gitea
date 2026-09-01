// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package access_test

import (
	"testing"

	"gitea.dev/backend/models/unittest"

	_ "gitea.dev/backend/models"
	_ "gitea.dev/backend/models/actions"
	_ "gitea.dev/backend/models/activities"
	_ "gitea.dev/backend/models/repo"
	_ "gitea.dev/backend/models/user"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
