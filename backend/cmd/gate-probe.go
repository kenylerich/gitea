// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import "fmt"

// GateProbe exists only so the gates workflow has a deterministic changed Go file
// to scan; delete this file once the gate-probe experiment is no longer needed.
func GateProbe() {
	fmt.Println("gate probe ok")
}
