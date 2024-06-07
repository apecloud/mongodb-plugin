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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubefakeclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	restfakeclient "k8s.io/client-go/rest/fake"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/apecloud/kubeblocks/apis/apps/v1alpha1"

	"github.com/apecloud/mongodb_plugin/constant"
)

const (
	ClusterName        = "fake-cluster-name"
	Namespace          = "fake-namespace"
	ComponentName      = "fake-component-name"
	ClusterCompName    = "fake-cluster-component-name"
	ClusterVersionName = "fake-cluster-version"
	ClusterDefName     = "fake-cluster-definition"
	ComponentDefName   = "fake-component-type"
	PodName            = "fake-pod-name"
)

func mockCluster(name, namespace string, conditions ...metav1.Condition) *v1alpha1.Cluster {
	var replicas int32 = 1

	return &v1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       v1alpha1.ClusterKind,
			APIVersion: v1alpha1.APIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       "b262b889-a27f-42d8-b066-2978561c8167",
		},
		Status: v1alpha1.ClusterStatus{
			Phase:      v1alpha1.RunningClusterPhase,
			Components: map[string]v1alpha1.ClusterComponentStatus{},
			Conditions: conditions,
		},
		Spec: v1alpha1.ClusterSpec{
			ClusterDefRef:     ClusterDefName,
			ClusterVersionRef: ClusterVersionName,
			TerminationPolicy: v1alpha1.WipeOut,
			ComponentSpecs: []v1alpha1.ClusterComponentSpec{
				{
					Name:            ComponentName,
					ComponentDefRef: ComponentDefName,
					Replicas:        replicas,
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("100Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("200m"),
							corev1.ResourceMemory: resource.MustParse("2Gi"),
						},
					},
					VolumeClaimTemplates: []v1alpha1.ClusterComponentVolumeClaimTemplate{
						{
							Name: "data",
							Spec: v1alpha1.PersistentVolumeClaimSpec{
								AccessModes: []corev1.PersistentVolumeAccessMode{
									corev1.ReadWriteOnce,
								},
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceStorage: resource.MustParse("1Gi"),
									},
								},
							},
						},
					},
				},
				{
					Name:            ComponentName + "-1",
					ComponentDefRef: ComponentDefName,
					Replicas:        replicas,
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("100Mi"),
						},
					},
					VolumeClaimTemplates: []v1alpha1.ClusterComponentVolumeClaimTemplate{
						{
							Name: "data",
							Spec: v1alpha1.PersistentVolumeClaimSpec{
								AccessModes: []corev1.PersistentVolumeAccessMode{
									corev1.ReadWriteOnce,
								},
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceStorage: resource.MustParse("1Gi"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func mockClusterRestClient(cluster *v1alpha1.Cluster, store *KubernetesStore) {
	_ = v1alpha1.AddToScheme(scheme.Scheme)
	store.client = &restfakeclient.RESTClient{
		GroupVersion:         v1alpha1.GroupVersion,
		NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		Client: restfakeclient.CreateHTTPClient(func(request *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader([]byte(runtime.EncodeOrDie(
					scheme.Codecs.LegacyCodec(scheme.Scheme.PrioritizedVersionsAllGroups()...),
					cluster),
				))),
			}
			return resp, nil
		}),
	}
}

func mockPods(replicas int, namespace string, cluster string) *corev1.PodList {
	pods := &corev1.PodList{}
	for i := 0; i < replicas; i++ {
		role := "follower"
		pod := corev1.Pod{}
		pod.Name = fmt.Sprintf("%s-pod-%d", cluster, i)
		pod.Namespace = namespace

		if i == 0 {
			role = "leader"
		}

		pod.Labels = map[string]string{
			constant.AppInstanceLabelKey:    cluster,
			constant.RoleLabelKey:           role,
			constant.KBAppComponentLabelKey: ComponentName,
			constant.AppNameLabelKey:        "mysql-apecloud-mysql",
			constant.AppManagedByLabelKey:   constant.AppName,
		}
		pod.Spec.NodeName = PodName
		pod.Spec.Containers = []corev1.Container{
			{
				Name:  "fake-container",
				Image: "fake-container-image",
				Ports: []corev1.ContainerPort{
					{
						Name:     "fake-port",
						HostPort: 1111,
					},
				},
			},
		}
		pod.Status.Phase = corev1.PodRunning
		pods.Items = append(pods.Items, pod)
	}
	return pods
}

func mockConfigMap(cmName string, namespace string, data map[string]string) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: Namespace,
		},
		Data: data,
	}
	if namespace != "" {
		cm.Namespace = namespace
	}
	return cm
}

func mockKubernetesStore() *KubernetesStore {
	ctx := context.TODO()
	logger := ctrl.Log.WithName("DCS-K8S-TEST")

	return &KubernetesStore{
		ctx:               ctx,
		clusterName:       ClusterName,
		componentName:     ComponentName,
		clusterCompName:   ClusterCompName,
		currentMemberName: PodName,
		namespace:         Namespace,
		client:            nil,
		clientset:         nil,
		logger:            logger,
	}
}

func TestInitialize(t *testing.T) {
	store := mockKubernetesStore()

	t.Run("get cluster failed", func(t *testing.T) {
		_ = v1alpha1.AddToScheme(scheme.Scheme)
		store.client = &restfakeclient.RESTClient{
			GroupVersion: v1alpha1.GroupVersion,
			Err:          fmt.Errorf("some error"),
		}

		err := store.Initialize()
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("initialize successfully", func(t *testing.T) {
		mockClusterRestClient(mockCluster("test", "test-ns"), store)
		store.clientset = kubefakeclient.NewSimpleClientset()

		err := store.Initialize()
		assert.Nil(t, err)
	})
}

func TestGetCluster(t *testing.T) {
	store := mockKubernetesStore()
	objs := []runtime.Object{
		mockPods(3, Namespace, ClusterName),
		mockConfigMap(store.getLeaderName(), Namespace, map[string]string{}),
	}

	t.Run("k8s get cluster error", func(t *testing.T) {
		_ = v1alpha1.AddToScheme(scheme.Scheme)
		store.client = &restfakeclient.RESTClient{
			GroupVersion: v1alpha1.GroupVersion,
			Err:          fmt.Errorf("some error"),
		}

		cluster, err := store.GetCluster()
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "some error")
		assert.Nil(t, cluster)
	})

	t.Run("get cluster success", func(t *testing.T) {
		mockClusterRestClient(mockCluster("test", "test-ns"), store)
		store.clientset = kubefakeclient.NewSimpleClientset(objs...)

		cluster, err := store.GetCluster()
		assert.Nil(t, err)
		assert.NotNil(t, cluster)
		assert.Equal(t, ClusterCompName, cluster.ClusterCompName)
		assert.Equal(t, Namespace, cluster.Namespace)
		assert.Equal(t, int32(1), cluster.Replicas)
		assert.Equal(t, 3, len(cluster.Members))
	})
}

func TestGetMembers(t *testing.T) {
	store := mockKubernetesStore()
	pods := mockPods(2, Namespace, ClusterName)

	store.clientset = kubefakeclient.NewSimpleClientset(pods)

	members, err := store.GetMembers()
	assert.NotNil(t, members)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(members))
}

func TestGetLeaderConfigMap(t *testing.T) {
	store := mockKubernetesStore()
	leaderConfigmap := mockConfigMap(store.getLeaderName(), Namespace, map[string]string{})

	t.Run("Leader configmap is not found", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset()
		leader, err := store.GetLeaderConfigMap()
		assert.Nil(t, leader)
		assert.Nil(t, err)
	})

	t.Run("get leader configmap success", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset(leaderConfigmap)
		leader, err := store.GetLeaderConfigMap()
		assert.Nil(t, err)
		assert.NotNil(t, leader)
		assert.Equal(t, store.getLeaderName(), leader.Name)
	})
}

func TestIsLeaseExist(t *testing.T) {
	store := mockKubernetesStore()

	configmap := mockConfigMap(store.getLeaderName(), Namespace, map[string]string{})
	cluster := mockCluster("test", "test-ns")

	t.Run("previous leader configmap resource exists", func(t *testing.T) {
		configmap.CreationTimestamp = metav1.NewTime(time.Unix(-1, 0))
		cluster.CreationTimestamp = metav1.NewTime(time.Unix(1, 0))
		store.clientset = kubefakeclient.NewSimpleClientset(configmap)
		mockClusterRestClient(cluster, store)

		_, err := store.GetCluster()
		assert.Nil(t, err)
		isExist, err := store.IsLeaseExist()
		assert.False(t, isExist)
		assert.Nil(t, err)
	})

	t.Run("Lease exist", func(t *testing.T) {
		configmap.CreationTimestamp = metav1.NewTime(time.Unix(1, 0))
		cluster.CreationTimestamp = metav1.NewTime(time.Unix(-1, 0))
		store.clientset = kubefakeclient.NewSimpleClientset(configmap)
		mockClusterRestClient(cluster, store)

		_, err := store.GetCluster()
		assert.Nil(t, err)
		isExist, err := store.IsLeaseExist()
		assert.True(t, isExist)
		assert.Nil(t, err)
	})
}

func TestCreateLease(t *testing.T) {
	store := mockKubernetesStore()
	configmap := mockConfigMap(store.getLeaderName(), Namespace, map[string]string{})
	cluster := mockCluster("test", "test-ns")

	t.Run("Lease exist", func(t *testing.T) {
		configmap.CreationTimestamp = metav1.NewTime(time.Unix(1, 0))
		cluster.CreationTimestamp = metav1.NewTime(time.Unix(-1, 0))
		store.clientset = kubefakeclient.NewSimpleClientset(configmap)
		mockClusterRestClient(cluster, store)

		_, err := store.GetCluster()
		assert.Nil(t, err)
		err = store.CreateLease()
		assert.Nil(t, err)
	})

	t.Run("create Lease success", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset()

		err := store.CreateLease()
		assert.Nil(t, err)
	})
}

func TestGetLeader(t *testing.T) {
	store := mockKubernetesStore()
	configmap := mockConfigMap(store.getLeaderName(), Namespace, map[string]string{})

	t.Run("get configmap nil", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset()

		leader, err := store.GetLeader()
		assert.Nil(t, leader)
		assert.Nil(t, err)
	})

	t.Run("get leader success", func(t *testing.T) {
		dbState := &DBState{
			OpTimestamp: 1000,
			Extra: map[string]string{
				"timeline": "1",
			},
		}
		dbStateStr, _ := json.Marshal(dbState)

		configmap.Annotations = map[string]string{
			"acquire-time": "100",
			"renew-time":   "101",
			"ttl":          "0",
			"leader":       "test-pod-0",
			"dbstate":      string(dbStateStr),
		}

		store.clientset = kubefakeclient.NewSimpleClientset(configmap)
		leader, err := store.GetLeader()
		assert.NotNil(t, leader)
		assert.Nil(t, err)
		assert.Equal(t, "test-pod-0", leader.Name)
		assert.Equal(t, int64(100), leader.AcquireTime)
		assert.Equal(t, int64(101), leader.RenewTime)
		assert.Equal(t, 0, leader.TTL)
		assert.Equal(t, dbState, leader.DBState)
	})
}

func TestAttemptAcquireLease(t *testing.T) {
	store := mockKubernetesStore()
	store.cluster = &Cluster{
		Leader: &Leader{
			Resource: mockConfigMap("test", Namespace, map[string]string{}),
			DBState: &DBState{
				OpTimestamp: 1,
			},
		},
		HaConfig: &HaConfig{
			ttl: 5,
		},
	}

	t.Run("Acquire Lease failed", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset()
		err := store.AttemptAcquireLease()
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, `configmaps "test" not found`)
	})

	t.Run("Acquire Lease success", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset(mockConfigMap("test", Namespace, map[string]string{}))

		err := store.AttemptAcquireLease()
		assert.Nil(t, err)
		leaderConfigmap := store.cluster.Leader.Resource.(*corev1.ConfigMap)
		assert.Equal(t, PodName, leaderConfigmap.Annotations["leader"])
		assert.Equal(t, "5", leaderConfigmap.Annotations["ttl"])
	})
}

func TestHasLease(t *testing.T) {
	store := mockKubernetesStore()

	t.Run("cluster nil", func(t *testing.T) {
		hasLease := store.HasLease()
		assert.False(t, hasLease)
	})

	store.cluster = &Cluster{}
	t.Run("leader nil", func(t *testing.T) {
		hasLease := store.HasLease()
		assert.False(t, hasLease)
	})

	store.cluster.Leader = &Leader{
		Name: store.currentMemberName,
	}
	t.Run("has Lease", func(t *testing.T) {
		hasLease := store.HasLease()
		assert.True(t, hasLease)
	})
}

func TestUpdateLease(t *testing.T) {
	store := mockKubernetesStore()
	configMap := mockConfigMap(store.getLeaderName(), Namespace, map[string]string{})
	store.cluster = &Cluster{
		Leader: &Leader{
			Resource: configMap,
			DBState: &DBState{
				OpTimestamp: 100,
			},
		},
	}

	t.Run("lost Lease", func(t *testing.T) {
		err := store.UpdateLease()
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "lost lease")
	})

	t.Run("update Lease successfully", func(t *testing.T) {
		configMap.Annotations = map[string]string{
			"leader": store.currentMemberName,
		}
		store.cluster.HaConfig = &HaConfig{
			ttl: 5,
		}
		store.clientset = kubefakeclient.NewSimpleClientset(configMap)

		err := store.UpdateLease()
		assert.Nil(t, err)
		newConfigMap, err := store.GetLeaderConfigMap()
		assert.Nil(t, err)
		assert.Equal(t, "5", newConfigMap.Annotations["ttl"])
	})
}

func TestReleaseLease(t *testing.T) {
	store := mockKubernetesStore()
	configMap := mockConfigMap(store.getLeaderName(), Namespace, map[string]string{})
	configMap.Annotations = map[string]string{}
	store.cluster = &Cluster{
		Leader: &Leader{
			Resource: configMap,
			DBState: &DBState{
				OpTimestamp: 100,
			},
		},
	}

	t.Run("configmap not found", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset()

		err := store.ReleaseLease()
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, `configmaps "fake-cluster-component-name-leader" not found`)
	})

	t.Run("release Lease successfully", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset(configMap)

		err := store.ReleaseLease()
		assert.Nil(t, err)
		newConfigMap, err := store.GetLeaderConfigMap()
		assert.Nil(t, err)
		assert.Equal(t, "", newConfigMap.Annotations["leader"])
	})
}

func TestHaConfig(t *testing.T) {
	store := mockKubernetesStore()
	configMap := mockConfigMap(store.getHAConfigName(), Namespace, map[string]string{})

	deleteMember := &MemberToDelete{
		UID:        "test",
		IsFinished: false,
	}
	deleteMemberStr, err := json.Marshal(deleteMember)
	assert.Nil(t, err)
	configMap.Annotations = map[string]string{
		"enable":         "true",
		"delete-members": string(deleteMemberStr),
	}
	store.cluster = &Cluster{
		Resource: mockCluster(ClusterName, Namespace),
	}

	t.Run("has previous ha config", func(t *testing.T) {
		cluster, ok := store.cluster.Resource.(*v1alpha1.Cluster)
		assert.True(t, ok)
		cluster.SetCreationTimestamp(metav1.NewTime(configMap.CreationTimestamp.Add(time.Second)))
		store.cluster.Resource = cluster
		store.clientset = kubefakeclient.NewSimpleClientset(configMap)

		err = store.CreateHaConfig("")
		assert.Nil(t, err)
	})

	t.Run("ha config exists", func(t *testing.T) {
		configMap.SetCreationTimestamp(metav1.NewTime(configMap.CreationTimestamp.Add(time.Second * 2)))
		store.clientset = kubefakeclient.NewSimpleClientset(configMap)

		err = store.CreateHaConfig("")
		assert.NotNil(t, err)
		haConfig, err := store.GetHaConfig()
		assert.Nil(t, err)
		assert.NotNil(t, haConfig.resource)
		assert.True(t, haConfig.enable)
		assert.Equal(t, int64(1048576), haConfig.maxLagOnSwitchover)
	})

	t.Run("create ha config successfully", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset()
		store.cluster = &Cluster{
			Resource: mockCluster(ClusterName, Namespace),
		}

		err = store.CreateHaConfig("")
		assert.Nil(t, err)
	})

	t.Run("no ha config", func(t *testing.T) {
		store.cluster = &Cluster{
			HaConfig: &HaConfig{},
		}

		err = store.UpdateHaConfig()
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "No HA configmap")
	})

	t.Run("update ha config successfully", func(t *testing.T) {
		store.cluster = &Cluster{
			HaConfig: &HaConfig{
				ttl:                10,
				maxLagOnSwitchover: 100,
				resource:           configMap,
			},
		}
		store.clientset = kubefakeclient.NewSimpleClientset(configMap)

		err = store.UpdateHaConfig()
		assert.Nil(t, err)
		haConfig, err := store.GetHaConfig()
		assert.Nil(t, err)
		assert.Equal(t, 10, haConfig.ttl)
		assert.Equal(t, int64(100), haConfig.maxLagOnSwitchover)
	})
}

func TestSwitchoverConfig(t *testing.T) {
	store := mockKubernetesStore()
	configMap := mockConfigMap(store.getSwitchoverName(), Namespace, map[string]string{})
	configMap.Annotations = map[string]string{
		"scheduled-at": "100",
		"leader":       "pod-0",
		"candidate":    "pod-1",
	}

	t.Run("there is another switchover unfinished", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset(configMap)

		err := store.CreateSwitchover("pod-0", "pod-1")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "there is another switchover fake-cluster-component-name-switchover unfinished")
		switchover, err := store.GetSwitchover()
		assert.Nil(t, err)
		assert.Equal(t, int64(100), switchover.ScheduledAt)
		assert.Equal(t, "pod-0", switchover.Leader)
		assert.Equal(t, "pod-1", switchover.Candidate)
	})

	t.Run("create switchover successfully", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset()
		store.cluster = &Cluster{
			Resource: mockCluster(ClusterName, Namespace),
		}

		err := store.CreateSwitchover("pod-0", "pod-1")
		assert.Nil(t, err)
		switchover, err := store.GetSwitchover()
		assert.Nil(t, err)
		assert.Equal(t, "pod-0", switchover.Leader)
		assert.Equal(t, "pod-1", switchover.Candidate)
	})

	t.Run("delete switchover failed", func(t *testing.T) {
		store.clientset = kubefakeclient.NewSimpleClientset()

		err := store.DeleteSwitchover()
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, `configmaps "fake-cluster-component-name-switchover" not found`)
	})
}
