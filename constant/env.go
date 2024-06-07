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

package constant

const (
	KBEnvNamespace            = "KB_NAMESPACE"
	KBEnvHostIP               = "KB_HOST_IP"
	KBEnvNodeName             = "KB_NODENAME"
	KBEnvPodName              = "KB_POD_NAME"
	KBEnvPodIP                = "KB_POD_IP"
	KBEnvPodUID               = "KB_POD_UID"
	KBEnvVolumeProtectionSpec = "KB_VOLUME_PROTECTION_SPEC"
)

// syncer env names
const (
	KBEnvClusterName     = "KB_CLUSTER_NAME"
	KBEnvComponentName   = "KB_COMP_NAME"
	KBEnvClusterCompName = "KB_CLUSTER_COMP_NAME"
	KBEnvWorkloadType    = "KB_WORKLOAD_TYPE"
	KBEnvCharacterType   = "KB_SERVICE_CHARACTER_TYPE"
	KBEnvServiceUser     = "KB_SERVICE_USER"
	KBEnvServicePassword = "KB_SERVICE_PASSWORD"
	KBEnvServicePort     = "KB_SERVICE_PORT"
	KBEnvBuiltinHandler  = "KB_BUILTIN_HANDLER"
	KBEnvTTL             = "KB_TTL"
	KBEnvMaxLag          = "KB_MAX_LAG"
	KBEnvEnableHA        = "KB_ENABLE_HA"
	KBEnvScriptsPath     = "KB_SCRIPTS_PATH"
)
