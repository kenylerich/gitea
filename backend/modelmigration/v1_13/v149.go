// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_13

import (
	"context"
	"fmt"

	"gitea.dev/backend/modelmigration/base"
	"gitea.dev/backend/modules/timeutil"
)

func AddCreatedAndUpdatedToMilestones(_ context.Context, x base.EngineMigration) error {
	type Milestone struct {
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
	}

	if err := x.Sync(new(Milestone)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}
	return nil
}
