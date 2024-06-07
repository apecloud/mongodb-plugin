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
	KBPrefix = "KB"
)

const (
	// ReasonNotFoundCR referenced custom resource not found
	ReasonNotFoundCR = "NotFound"
	// ReasonRefCRUnavailable  referenced custom resource is unavailable
	ReasonRefCRUnavailable = "Unavailable"
	// ReasonDeletedCR deleted custom resource
	ReasonDeletedCR = "DeletedCR"
	// ReasonDeletingCR deleting custom resource
	ReasonDeletingCR = "DeletingCR"
	// ReasonCreatedCR created custom resource
	ReasonCreatedCR = "CreatedCR"
	// ReasonRunTaskFailed run task failed
	ReasonRunTaskFailed = "RunTaskFailed"
	// ReasonDeleteFailed delete failed
	ReasonDeleteFailed = "DeleteFailed"
)

const (
	// Container port name
	LorryHTTPPortName                  = "lorry-http-port"
	LorryGRPCPortName                  = "lorry-grpc-port"
	SyncerHTTPPortName                 = "ha"
	ProbeInitContainerName             = "kb-initprobe"
	StatusProbeContainerName           = "kb-checkstatus"
	RunningProbeContainerName          = "kb-checkrunning"
	VolumeProtectionProbeContainerName = "kb-volume-protection"

	// the filedpath name used in event.InvolvedObject.FieldPath
	ProbeCheckStatusPath  = "spec.containers{" + StatusProbeContainerName + "}"
	ProbeCheckRunningPath = "spec.containers{" + RunningProbeContainerName + "}"
)

const (
	Primary   = "primary"
	Secondary = "secondary"

	Leader    = "leader"
	Follower  = "follower"
	Learner   = "learner"
	Candidate = "candidate"
)

// switchover constants

// username and password are keys in created secrets for others to refer to.
const (
	AccountNameForSecret   = "username"
	AccountPasswdForSecret = "password"
)

const (
	KubernetesClusterDomainEnv = "KUBERNETES_CLUSTER_DOMAIN"
	DefaultDNSDomain           = "cluster.local"
)

const (
	Replication = "Replication"
	Consensus   = "Consensus"
)

const (
	APIGroup = "kubeblocks.io"

	AppName = "kubeblocks"

	// K8s recommonded well-known labels and annotation keys
	AppInstanceLabelKey  = "app.kubernetes.io/instance"
	AppNameLabelKey      = "app.kubernetes.io/name"
	AppComponentLabelKey = "app.kubernetes.io/component"
	AppVersionLabelKey   = "app.kubernetes.io/version"
	AppManagedByLabelKey = "app.kubernetes.io/managed-by"
	RegionLabelKey       = "topology.kubernetes.io/region"
	ZoneLabelKey         = "topology.kubernetes.io/zone"

	// kubeblocks.io labels
	BackupProtectionLabelKey                 = "kubeblocks.io/backup-protection" // BackupProtectionLabelKey Backup delete protection policy label
	BackupToolTypeLabelKey                   = "kubeblocks.io/backup-tool-type"
	AddonProviderLabelKey                    = "kubeblocks.io/provider" // AddonProviderLabelKey marks the addon provider
	RoleLabelKey                             = "kubeblocks.io/role"     // RoleLabelKey consensusSet and replicationSet role label key
	ModeKey                                  = "kubeblocks.io/mode"     // ModeKey is in enum of standalone/replication/raftGroup
	VolumeTypeLabelKey                       = "kubeblocks.io/volume-type"
	ClusterAccountLabelKey                   = "account.kubeblocks.io/name"
	KBAppComponentLabelKey                   = "apps.kubeblocks.io/component-name"
	KBAppComponentDefRefLabelKey             = "apps.kubeblocks.io/component-def-ref"
	AppConfigTypeLabelKey                    = "apps.kubeblocks.io/config-type"
	KBManagedByKey                           = "apps.kubeblocks.io/managed-by" // KBManagedByKey marks resources that auto created
	PVCNameLabelKey                          = "apps.kubeblocks.io/pvc-name"
	VolumeClaimTemplateNameLabelKey          = "apps.kubeblocks.io/vct-name"
	VolumeClaimTemplateNameLabelKeyForLegacy = "vct.kubeblocks.io/name" // Deprecated: only compatible with version 0.5, will be removed in 0.7
	WorkloadTypeLabelKey                     = "apps.kubeblocks.io/workload-type"
	ClassProviderLabelKey                    = "class.kubeblocks.io/provider"
	ClusterDefLabelKey                       = "clusterdefinition.kubeblocks.io/name"
	ClusterVerLabelKey                       = "clusterversion.kubeblocks.io/name"
	CMConfigurationSpecProviderLabelKey      = "config.kubeblocks.io/config-spec"    // CMConfigurationSpecProviderLabelKey is ComponentConfigSpec name
	CMConfigurationCMKeysLabelKey            = "config.kubeblocks.io/configmap-keys" // CMConfigurationCMKeysLabelKey Specify configmap keys
	CMConfigurationTemplateNameLabelKey      = "config.kubeblocks.io/config-template-name"
	CMTemplateNameLabelKey                   = "config.kubeblocks.io/template-name"
	CMConfigurationTypeLabelKey              = "config.kubeblocks.io/config-type"
	CMInsConfigurationHashLabelKey           = "config.kubeblocks.io/config-hash"
	CMInsCurrentConfigurationHashLabelKey    = "config.kubeblocks.io/update-config-hash"
	CMConfigurationConstraintsNameLabelKey   = "config.kubeblocks.io/config-constraints-name"
	CMConfigurationTemplateVersion           = "config.kubeblocks.io/config-template-version"
	ConsensusSetAccessModeLabelKey           = "cs.apps.kubeblocks.io/access-mode"
	BackupTypeLabelKeyKey                    = "dataprotection.kubeblocks.io/backup-type"
	DataProtectionLabelBackupNameKey         = "dataprotection.kubeblocks.io/backup-name"
	AddonNameLabelKey                        = "extensions.kubeblocks.io/addon-name"
	OpsRequestTypeLabelKey                   = "ops.kubeblocks.io/ops-type"
	OpsRequestNameLabelKey                   = "ops.kubeblocks.io/ops-name"
	ServiceDescriptorNameLabelKey            = "servicedescriptor.kubeblocks.io/name"

	// kubeblocks.io annotations
	ClusterSnapshotAnnotationKey                = "kubeblocks.io/cluster-snapshot"            // ClusterSnapshotAnnotationKey saves the snapshot of cluster.
	DefaultClusterVersionAnnotationKey          = "kubeblocks.io/is-default-cluster-version"  // DefaultClusterVersionAnnotationKey specifies the default cluster version.
	OpsRequestAnnotationKey                     = "kubeblocks.io/ops-request"                 // OpsRequestAnnotationKey OpsRequest annotation key in Cluster
	ReconcileAnnotationKey                      = "kubeblocks.io/reconcile"                   // ReconcileAnnotationKey Notify k8s object to reconcile
	RestartAnnotationKey                        = "kubeblocks.io/restart"                     // RestartAnnotationKey the annotation which notices the StatefulSet/DeploySet to restart
	RestoreFromTimeAnnotationKey                = "kubeblocks.io/restore-from-time"           // RestoreFromTimeAnnotationKey specifies the time to recover from the backup.
	RestoreFromSrcClusterAnnotationKey          = "kubeblocks.io/restore-from-source-cluster" // RestoreFromSrcClusterAnnotationKey specifies the source cluster to recover from the backup.
	RestoreFromBackUpAnnotationKey              = "kubeblocks.io/restore-from-backup"         // RestoreFromBackUpAnnotationKey specifies the component to recover from the backup.
	SnapShotForStartAnnotationKey               = "kubeblocks.io/snapshot-for-start"
	ComponentReplicasAnnotationKey              = "apps.kubeblocks.io/component-replicas" // ComponentReplicasAnnotationKey specifies the number of pods in replicas
	BackupPolicyTemplateAnnotationKey           = "apps.kubeblocks.io/backup-policy-template"
	LastAppliedClusterAnnotationKey             = "apps.kubeblocks.io/last-applied-cluster"
	PVLastClaimPolicyAnnotationKey              = "apps.kubeblocks.io/pv-last-claim-policy"
	HaltRecoveryAllowInconsistentCVAnnotKey     = "clusters.apps.kubeblocks.io/allow-inconsistent-cv"
	HaltRecoveryAllowInconsistentResAnnotKey    = "clusters.apps.kubeblocks.io/allow-inconsistent-resource"
	LeaderAnnotationKey                         = "cs.apps.kubeblocks.io/leader"
	PrimaryAnnotationKey                        = "rs.apps.kubeblocks.io/primary"
	DefaultBackupPolicyAnnotationKey            = "dataprotection.kubeblocks.io/is-default-policy"          // DefaultBackupPolicyAnnotationKey specifies the default backup policy.
	DefaultBackupPolicyTemplateAnnotationKey    = "dataprotection.kubeblocks.io/is-default-policy-template" // DefaultBackupPolicyTemplateAnnotationKey specifies the default backup policy template.
	DefaultBackupRepoAnnotationKey              = "dataprotection.kubeblocks.io/is-default-repo"            // DefaultBackupRepoAnnotationKey specifies the default backup repo.
	BackupDataPathPrefixAnnotationKey           = "dataprotection.kubeblocks.io/path-prefix"                // BackupDataPathPrefixAnnotationKey specifies the backup data path prefix.
	ReconfigureRefAnnotationKey                 = "dataprotection.kubeblocks.io/reconfigure-ref"
	DataProtectionLabelClusterUIDKey            = "dataprotection.kubeblocks.io/cluster-uid"
	DisableUpgradeInsConfigurationAnnotationKey = "config.kubeblocks.io/disable-reconfigure"
	LastAppliedConfigAnnotationKey              = "config.kubeblocks.io/last-applied-configuration"
	LastAppliedOpsCRAnnotationKey               = "config.kubeblocks.io/last-applied-ops-name"
	UpgradePolicyAnnotationKey                  = "config.kubeblocks.io/reconfigure-policy"
	KBParameterUpdateSourceAnnotationKey        = "config.kubeblocks.io/reconfigure-source"
	UpgradeRestartAnnotationKey                 = "config.kubeblocks.io/restart"
	ConfigAppliedVersionAnnotationKey           = "config.kubeblocks.io/config-applied-version"
	KubeBlocksGenerationKey                     = "kubeblocks.io/generation"
	ExtraEnvAnnotationKey                       = "kubeblocks.io/extra-env"
	LastRoleSnapshotVersionAnnotationKey        = "apps.kubeblocks.io/last-role-snapshot-version"

	// kubeblocks.io well-known finalizers
	DBClusterFinalizerName             = "cluster.kubeblocks.io/finalizer"
	ConfigurationTemplateFinalizerName = "config.kubeblocks.io/finalizer"
	ServiceDescriptorFinalizerName     = "servicedescriptor.kubeblocks.io/finalizer"

	// ConfigurationTplLabelPrefixKey clusterVersion or clusterdefinition using tpl
	ConfigurationTplLabelPrefixKey         = "config.kubeblocks.io/tpl"
	ConfigurationConstraintsLabelPrefixKey = "config.kubeblocks.io/constraints"

	// CMInsLastReconfigurePhaseKey defines the current phase
	CMInsLastReconfigurePhaseKey = "config.kubeblocks.io/last-applied-reconfigure-phase"

	// ConfigurationRevision defines the current revision
	// TODO support multi version
	ConfigurationRevision          = "config.kubeblocks.io/configuration-revision"
	LastConfigurationRevisionPhase = "config.kubeblocks.io/revision-reconcile-phase"

	// Deprecated: only compatible with version 0.6, will be removed in 0.8
	// CMInsEnableRerenderTemplateKey is used to enable rerender template
	CMInsEnableRerenderTemplateKey = "config.kubeblocks.io/enable-rerender"

	// IgnoreResourceConstraint is used to specify whether to ignore the resource constraint
	IgnoreResourceConstraint = "resource.kubeblocks.io/ignore-constraint"

	RBACRoleName        = "kubeblocks-cluster-pod-role"
	RBACClusterRoleName = "kubeblocks-volume-protection-pod-role"
)
