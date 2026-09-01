// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package db_test

import (
	"testing"

	"gitea.dev/backend/models/unittest"

	_ "gitea.dev/backend/models"
	_ "gitea.dev/backend/models/repo"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
