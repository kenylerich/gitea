// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package queue

import (
	"testing"

	"gitea.dev/backend/modules/setting"
	"gitea.dev/backend/modules/test"
)

func TestBaseRedis(t *testing.T) {
	redisConn := test.PrepareTestRedis(t)
	queueSetting := setting.QueueSettings{Length: 10, ConnStr: redisConn}
	optsSimple := testQueueBasicOptions{NotifiableQueue: true}
	optsUnique := testQueueBasicOptions{UniqueQueue: true, NotifiableQueue: true}
	testQueueBasic(t, newBaseRedisSimple, toBaseConfig("baseRedis", queueSetting), optsSimple)
	testQueueBasic(t, newBaseRedisUnique, toBaseConfig("baseRedisUnique", queueSetting), optsUnique)
}
