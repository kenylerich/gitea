// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_21

import (
	"context"

	"gitea.dev/backend/modelmigration/base"
)

func AddTriggerEventToActionRun(_ context.Context, x base.EngineMigration) error {
	type ActionRun struct {
		TriggerEvent string
	}

	return x.Sync(new(ActionRun))
}
