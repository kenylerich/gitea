// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markdown

import (
	"os"
	"testing"

	"gitea.dev/backend/modules/markup"
	"gitea.dev/backend/modules/setting"
)

func TestMain(m *testing.M) {
	setting.IsInTesting = true
	markup.RenderBehaviorForTesting.DisableAdditionalAttributes = true
	os.Exit(m.Run())
}
