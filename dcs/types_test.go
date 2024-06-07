/*
Copyright (C) 2022-2023 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package dcs

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/apecloud/mongodb_plugin/constant"
)

func TestCluster(t *testing.T) {
	cluster := Cluster{
		Namespace:       Namespace,
		ClusterCompName: ClusterCompName,
	}

	t.Run("cluster has member", func(t *testing.T) {
		cluster.Members = []Member{
			{
				Name:  "pod-0",
				PodIP: "pod-0",
			},
		}

		hasMember := cluster.HasMember("pod-0")
		assert.True(t, hasMember)
		hasMember = cluster.HasMember("pod-1")
		assert.False(t, hasMember)
	})

	t.Run("get leader member", func(t *testing.T) {
		leader := cluster.GetLeaderMember()
		assert.Nil(t, leader)

		cluster.Leader = &Leader{
			Name: "pod-1",
		}

		leader = cluster.GetLeaderMember()
		assert.Nil(t, leader)

		cluster.Members = []Member{
			{
				Name:  "pod-1",
				PodIP: "pod-1",
			},
		}

		leader = cluster.GetLeaderMember()
		assert.Equal(t, "pod-1", leader.Name)
	})

	t.Run("get member with host", func(t *testing.T) {
		viper.Set(constant.KubernetesClusterDomainEnv, "cluster.local")
		defer viper.Reset()

		member := cluster.GetMemberWithHost(fmt.Sprintf("%s.%s-headless.%s.svc.%s", "pod-2", cluster.ClusterCompName, cluster.Namespace, "cluster.local"))
		assert.Nil(t, member)

		cluster.Members = []Member{
			{
				Name:  "pod-2",
				PodIP: "pod-2",
			},
		}

		member = cluster.GetMemberWithHost(fmt.Sprintf("%s.%s-headless.%s.svc.%s", "pod-2", cluster.ClusterCompName, cluster.Namespace, "cluster.local"))
		assert.NotNil(t, member)
		assert.Equal(t, "pod-2", member.Name)
	})

	t.Run("get member name and addrs", func(t *testing.T) {
		viper.Set(constant.KubernetesClusterDomainEnv, "cluster.local")
		defer viper.Reset()

		cluster.Members = []Member{
			{
				Name:   "pod-3",
				PodIP:  "pod-3",
				DBPort: "1",
			},
			{
				Name:   "pod-4",
				PodIP:  "pod-3",
				DBPort: "1",
			},
		}

		memberNames := cluster.GetMemberName()
		assert.Equal(t, []string{"pod-3", "pod-4"}, memberNames)

		addrs := cluster.GetMemberAddrs()
		assert.Equal(t, []string{"pod-3.pod-headless.fake-namespace.svc.cluster.local:1",
			"pod-4.pod-headless.fake-namespace.svc.cluster.local:1"}, addrs)
	})

	t.Run("is locked", func(t *testing.T) {
		cluster.Leader = nil
		isLocked := cluster.IsLocked()
		assert.False(t, isLocked)

		cluster.Leader = &Leader{
			Name: "pod-5",
		}
		isLocked = cluster.IsLocked()
		assert.True(t, isLocked)
	})
}

func TestHAConfig(t *testing.T) {
	haConfig := &HaConfig{
		ttl:                5,
		maxLagOnSwitchover: 100,
		enable:             true,
	}

	ttl := haConfig.GetTTL()
	maxLag := haConfig.GetMaxLagOnSwitchover()
	enable := haConfig.IsEnable()
	assert.Equal(t, 5, ttl)
	assert.Equal(t, int64(100), maxLag)
	assert.True(t, enable)

	t.Run("test delete", func(t *testing.T) {
		haConfig.DeleteMembers = map[string]MemberToDelete{
			"pod-0": {
				UID:        "test-0",
				IsFinished: false,
			},
			"pod-1": {
				UID:        "test-1",
				IsFinished: true,
			},
		}
		member0 := &Member{
			Name: "pod-0",
			UID:  "test-0",
		}
		member1 := &Member{
			Name: "pod-1",
			UID:  "test-1",
		}
		member2 := &Member{
			Name: "pod-2",
			UID:  "test-2",
		}

		isDeleting := haConfig.IsDeleting(member0)
		assert.True(t, isDeleting)

		isDeleted := haConfig.IsDeleted(member0)
		assert.False(t, isDeleted)
		isDeleted = haConfig.IsDeleted(member1)
		assert.True(t, isDeleted)
		member1.UID = "test"
		isDeleted = haConfig.IsDeleted(member1)
		assert.False(t, isDeleted)
		isDeleted = haConfig.IsDeleted(member2)
		assert.False(t, isDeleted)

		haConfig.FinishDeleted(member0)
		assert.True(t, haConfig.DeleteMembers["pod-0"].IsFinished)

		haConfig.AddMemberToDelete(member2)
		isDeleted = haConfig.IsDeleting(member2)
		assert.True(t, isDeleted)
	})
}
