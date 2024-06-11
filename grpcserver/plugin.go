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

package grpcserver

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/apecloud/kubeblocks/pkg/kb_agent/plugin"
	"github.com/apecloud/mongodb_plugin/constant"
	"github.com/apecloud/mongodb_plugin/dcs"
	"github.com/apecloud/mongodb_plugin/mongodb"
)

type DBPlugin struct {
	plugin.UnimplementedEnginePluginServer
	dbManager *mongodb.Manager
	store     dcs.DCS
}

func NewDBPlugin() *DBPlugin {
	dbManager, _ := mongodb.NewManager(nil)
	return &DBPlugin{
		dbManager: dbManager,
		store:     dcs.GetStore(),
	}
}

func (p *DBPlugin) GetPluginInfo(ctx context.Context, in *plugin.GetPluginInfoRequest) (*plugin.GetPluginInfoResponse, error) {
	resp := &plugin.GetPluginInfoResponse{
		Name:       "DBPlugin",
		EngineType: viper.GetString(constant.KBEnvCharacterType),
		Version:    "0.1.0",
	}
	return resp, nil
}

func (p *DBPlugin) IsEngineReady(ctx context.Context, in *plugin.IsEngineReadyRequest) (*plugin.IsEngineReadyResponse, error) {
	isReady := p.dbManager.IsDBStartupReady()

	resp := &plugin.IsEngineReadyResponse{
		Ready: isReady,
	}
	return resp, nil
}

func (p *DBPlugin) GetRole(ctx context.Context, in *plugin.GetRoleRequest) (*plugin.GetRoleResponse, error) {
	role, err := p.dbManager.GetReplicaRole(ctx, nil)

	if err != nil {
		return nil, err
	}

	resp := &plugin.GetRoleResponse{
		Role: role,
	}
	return resp, nil
}

func (p *DBPlugin) JoinMember(ctx context.Context, in *plugin.JoinMemberRequest) (*plugin.JoinMemberResponse, error) {
	cluster, err := p.store.GetCluster()
	if err != nil {
		return nil, err
	}

	memberName := in.NewMember
	err = p.dbManager.JoinMemberToCluster(ctx, cluster, memberName)
	if err != nil {
		return nil, err
	}

	resp := &plugin.JoinMemberResponse{}
	return resp, nil
}

func (p *DBPlugin) LeaveMember(ctx context.Context, in *plugin.LeaveMemberRequest) (*plugin.LeaveMemberResponse, error) {
	memberName := in.LeaveMember
	cluster, err := p.store.GetCluster()
	if err != nil {
		return nil, err
	}

	err = p.dbManager.LeaveMemberFromCluster(ctx, cluster, memberName)
	if err != nil {
		return nil, err
	}

	resp := &plugin.LeaveMemberResponse{}
	return resp, nil
}

func (p *DBPlugin) ReadOnly(ctx context.Context, in *plugin.ReadOnlyRequest) (*plugin.ReadOnlyResponse, error) {
	err := p.dbManager.Lock(ctx, in.Reason)
	return &plugin.ReadOnlyResponse{}, err
}

func (p *DBPlugin) ReadWrite(ctx context.Context, in *plugin.ReadWriteRequest) (*plugin.ReadWriteResponse, error) {
	err := p.dbManager.Unlock(ctx)
	return &plugin.ReadWriteResponse{}, err
}

// func (p *DBPlugin) AccountProvision(ctx context.Context, in *plugin.AccountProvisionRequest) (*plugin.AccountProvisionResponse, error) {
// 	userInfo := models.UserInfo{
// 		UserName: in.UserName,
// 		Password: in.Password,
// 		RoleName: in.Role,
// 	}
//
// 	err := p.dbManager.CreateUser(ctx, userInfo.UserName, userInfo.Password)
// 	if err != nil {
// 		return &plugin.AccountProvisionResponse{}, err
// 	}
//
// 	if userInfo.RoleName != "" {
// 		err := p.dbManager.GrantUserRole(ctx, userInfo.UserName, userInfo.RoleName)
// 		if err != nil {
// 			return &plugin.AccountProvisionResponse{}, err
// 		}
// 	}
// 	return &plugin.AccountProvisionResponse{}, err
// }

func (p *DBPlugin) Switchover(ctx context.Context, in *plugin.SwitchoverRequest) (*plugin.SwitchoverResponse, error) {
	resp := &plugin.SwitchoverResponse{}
	primary := in.Primary
	candidate := in.Candidate
	if primary == "" && candidate == "" {
		return resp, errors.New("primary or candidate must be set")
	}

	cluster, err := p.store.GetCluster()
	if cluster == nil {
		return resp, errors.Wrap(err, "get cluster failed")
	}

	if cluster.HaConfig == nil || !cluster.HaConfig.IsEnable() {
		return resp, errors.New("cluster's ha is disabled")
	}
	if primary != "" {
		leaderMember := cluster.GetMemberWithName(primary)
		if leaderMember == nil {
			message := fmt.Sprintf("primary %s not exists", primary)
			return resp, errors.New(message)
		}

		ok, err := p.dbManager.IsLeaderMember(ctx, cluster, leaderMember)
		if err != nil {
			return resp, errors.Wrap(err, "check leader member failed")
		}
		if !ok {
			message := fmt.Sprintf("%s is not the primary", primary)
			return resp, errors.New(message)
		}
	}

	if candidate != "" {
		candidateMember := cluster.GetMemberWithName(candidate)
		if candidateMember == nil {
			message := fmt.Sprintf("candidate %s not exists", candidate)
			return resp, errors.New(message)
		}

		if !p.dbManager.IsMemberHealthy(ctx, cluster, candidateMember) {
			message := fmt.Sprintf("candidate %s is unhealthy", candidate)
			return resp, errors.New(message)
		}
	} else if len(p.dbManager.HasOtherHealthyMembers(ctx, cluster, primary)) == 0 {
		return resp, errors.New("candidate is not set and has no other healthy members")
	}

	err = p.store.CreateSwitchover(primary, candidate)
	if err != nil {
		message := fmt.Sprintf("Create switchover failed: %v", err)
		return resp, errors.New(message)
	}

	return resp, nil
}
