/*
Copyright 2020 The Nocalhost Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"encoding/base64"
	"nocalhost/internal/nocalhost-api/service"
	"nocalhost/pkg/nocalhost-api/app/api"
	"nocalhost/pkg/nocalhost-api/pkg/clientgo"
	"nocalhost/pkg/nocalhost-api/pkg/errno"
	"nocalhost/pkg/nocalhost-api/pkg/log"
	"nocalhost/pkg/nocalhost-api/pkg/setupcluster"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
)

// Create Add cluster
// @Summary Add cluster
// @Description Add cluster
// @Tags Cluster
// @Accept  json
// @Produce  json
// @param Authorization header string true "Authorization"
// @Param createCluster body cluster.CreateClusterRequest true "The cluster info"
// @Success 200 {object} model.ClusterModel
// @Router /v1/cluster [post]
func Create(c *gin.Context) {
	var req CreateClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("createCluster bind params err: %v", err)
		api.SendResponse(c, errno.ErrBind, nil)
		return
	}
	// decode kubeconfig
	DecKubeconfig, err := base64.StdEncoding.DecodeString(req.KubeConfig)
	if err != nil {
		api.SendResponse(c, errno.ErrClusterKubeErr, nil)
		return
	}
	// check kubeconfig server already exist
	t := KubeConfig{}
	err = yaml.Unmarshal(DecKubeconfig, &t)
	if err != nil {
		api.SendResponse(c, errno.ErrClusterKubeErr, nil)
		return
	}
	// only allow one cluster kubeconfig
	if len(t.Clusters) != 1 {
		api.SendResponse(c, errno.ErrClusterKubeCreate, nil)
		return
	}
	where := make(map[string]interface{}, 0)
	where["server"] = t.Clusters[0].Cluster.Server
	_, err = service.Svc.ClusterSvc().GetAny(c, where)
	if err == nil {
		api.SendResponse(c, errno.ErrClusterExistCreate, nil)
		return
	}

	goClient, err := clientgo.NewAdminGoClient(DecKubeconfig)

	// get client go and check if is admin Kubeconfig
	if err != nil {
		switch err.(type) {
		case *errno.Errno:
			api.SendResponse(c, err, nil)
		default:
			api.SendResponse(c, errno.ErrClusterKubeErr, nil)
		}
		return
	}

	// 1. check if Namespace nocalhost-reserved already exist, ignore cause by nocalhost-dep-job installer.sh has checkout this condition and will exit
	// 2. use admin Kubeconfig create configmap for nocalhost-dep-job to create admission webhook cert
	// 3. deploy nocalhost-dep-job and pull on nocalhost-dep
	// see https://codingcorp.coding.net/p/nocalhost/wiki/115
	clusterSetUp := setupcluster.NewSetUpCluster(goClient)

	clusterInfo, err, errRes := clusterSetUp.InitCluster("")
	if err != nil {
		api.SendResponse(c, errRes, nil)
		return
	}

	// Pre pull images DaemonSet
	prePullImages, _ := service.Svc.PrePull().GetAll(c)
	_, err = goClient.DeployPrePullImages(prePullImages, "")
	if err != nil {
		log.Warnf("deploy pre pull images err: %v", err)
	}

	userId, _ := c.Get("userId")
	cluster, err := service.Svc.ClusterSvc().Create(c, req.Name, string(DecKubeconfig), req.StorageClass, t.Clusters[0].Cluster.Server, clusterInfo, userId.(uint64))
	if err != nil {
		log.Warnf("create cluster err: %v", err)
		api.SendResponse(c, errno.ErrClusterCreate, nil)
		return
	}

	api.SendResponse(c, nil, cluster)
}
