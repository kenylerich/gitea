// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repository

import (
	"testing"

	"gitea.dev/backend/models/unittest"

	_ "gitea.dev/backend/models"
	_ "gitea.dev/backend/models/actions"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
