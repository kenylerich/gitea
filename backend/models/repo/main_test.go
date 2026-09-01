// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo_test

import (
	"testing"

	"gitea.dev/backend/models/unittest"

	_ "gitea.dev/backend/models" // register table model
	_ "gitea.dev/backend/models/actions"
	_ "gitea.dev/backend/models/activities"
	_ "gitea.dev/backend/models/perm/access" // register table model
	_ "gitea.dev/backend/models/repo"        // register table model
	_ "gitea.dev/backend/models/user"        // register table model
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
