// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package issues_test

import (
	"testing"

	issues_model "gitea.dev/backend/models/issues"
	"gitea.dev/backend/models/unittest"

	_ "gitea.dev/backend/models"
	_ "gitea.dev/backend/models/actions"
	_ "gitea.dev/backend/models/activities"
	_ "gitea.dev/backend/models/repo"
	_ "gitea.dev/backend/models/user"

	"github.com/stretchr/testify/assert"
)

func TestFixturesAreConsistent(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	unittest.CheckConsistencyFor(t,
		&issues_model.Issue{},
		&issues_model.PullRequest{},
		&issues_model.Milestone{},
		&issues_model.Label{},
	)
}

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
