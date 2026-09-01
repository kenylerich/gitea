// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package git

// RefCommit represents a commit resolved from a ref.
// It is used to avoid a services -> routers dependency.
type RefCommit struct {
	InputRef string
	RefName  RefName
	Commit   *Commit
	CommitID string
}

// NewRefCommit creates a RefCommit from a ref name and a resolved commit.
func NewRefCommit(refName RefName, commit *Commit) *RefCommit {
	return &RefCommit{InputRef: refName.ShortName(), RefName: refName, Commit: commit, CommitID: commit.ID.String()}
}
