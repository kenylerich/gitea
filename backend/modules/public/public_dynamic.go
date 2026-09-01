// Copyright 2016 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !bindata

package public

import (
	"gitea.dev/backend/modules/assetfs"
	"gitea.dev/backend/modules/setting"
)

func BuiltinAssets() *assetfs.Layer {
	return assetfs.Local("builtin(static)", setting.StaticRootPath, "public")
}
