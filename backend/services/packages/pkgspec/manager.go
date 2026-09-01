// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package pkgspec

import (
	packages_model "gitea.dev/backend/models/packages"
	packages_service "gitea.dev/backend/services/packages"
	"gitea.dev/backend/services/packages/debian"
	"gitea.dev/backend/services/packages/terraform"
)

func InitManager() error {
	mgr := packages_service.GetSpecManager()
	mgr.Add(packages_model.TypeDebian, &debian.Specialization{})
	mgr.Add(packages_model.TypeTerraformState, &terraform.Specialization{})
	// TODO: add more in the future, refactor the existing code to use this approach
	return nil
}
