// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package project

import (
	"testing"

	"gitea.dev/backend/models/unittest"

	_ "gitea.dev/backend/models/repo"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		FixtureFiles: []string{
			"project.yml",
			"project_board.yml",
			"project_issue.yml",
			"repository.yml",
		},
	})
}
