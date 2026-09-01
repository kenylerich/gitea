// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repostats

import (
	"testing"

	activities_model "gitea.dev/backend/models/activities"
	"gitea.dev/backend/models/organization"
	repo_model "gitea.dev/backend/models/repo"
	"gitea.dev/backend/models/unittest"
	user_model "gitea.dev/backend/models/user"

	_ "gitea.dev/backend/models/actions"
	_ "gitea.dev/backend/models/system"

	"github.com/stretchr/testify/assert"
)

// TestFixturesAreConsistent assert that test fixtures are consistent
func TestFixturesAreConsistent(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	unittest.CheckConsistencyFor(t,
		&user_model.User{},
		&repo_model.Repository{},
		&organization.Team{},
		&activities_model.Action{})
}

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
