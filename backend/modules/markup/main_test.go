// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markup_test

import (
	"os"
	"testing"

	"gitea.dev/backend/modules/markup"
	"gitea.dev/backend/modules/setting"
)

func TestMain(m *testing.M) {
	setting.IsInTesting = true
	markup.RenderBehaviorForTesting.DisableAdditionalAttributes = true
	markup.RefreshFileNamePatterns()
	os.Exit(m.Run())
}
