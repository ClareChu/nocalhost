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

package cmds

import (
	"fmt"
	"io/ioutil"
	"nocalhost/internal/nhctl/app"
	"nocalhost/internal/nhctl/fp"
	"nocalhost/internal/nhctl/nocalhost"
	"nocalhost/internal/nhctl/utils"
	"nocalhost/pkg/nhctl/clientgoutils"
	"nocalhost/pkg/nhctl/log"
	"os"
	"strings"
)

func initApp(appName string) {
	var err error

	if nameSpace == "" {
		nameSpace, err = clientgoutils.GetNamespaceFromKubeConfig(kubeConfig)
		if err != nil {
			log.FatalE(err, "Failed to get namespace")
		}
		if nameSpace == "" {
			log.Fatal("--namespace or --kubeconfig mush be provided")
		}
	}

	if !nocalhost.CheckIfApplicationExist(appName, nameSpace) {
		// init nocalhost default application and wait
		if appName == app.DefaultNocalhostApplication {
			err := initNocalhostDefaultApplicationAndWait()
			if err != nil {
				log.FatalE(err, "Failed to init default application in ns "+nameSpace)
			}
		} else {
			log.FatalE(err, fmt.Sprintf("Application \"%s\" not found", appName))
		}
	}
	nocalhostApp, err = app.NewApplication(appName, nameSpace, kubeConfig, true)
	if err != nil {
		log.FatalE(err, "Failed to get application info")
	}
	log.AddField("APP", nocalhostApp.Name)
}

func CheckIfSvcExist(svcName string, svcType ...string) {
	serviceType := app.Deployment
	if len(svcType) > 0 {
		svcTypeLower := strings.ToLower(svcType[0])
		switch svcTypeLower {
		case strings.ToLower(string(app.StatefulSet)):
			serviceType = app.StatefulSet
		case strings.ToLower(string(app.DaemonSet)):
			serviceType = app.DaemonSet
		case strings.ToLower(string(app.Job)):
			serviceType = app.Job
		case strings.ToLower(string(app.CronJob)):
			serviceType = app.CronJob
		default:
			serviceType = app.Deployment
		}
	}
	if svcName == "" {
		log.Fatal("please use -d to specify a k8s workload")
	}
	exist, err := nocalhostApp.CheckIfSvcExist(svcName, serviceType)
	if err != nil {
		log.FatalE(err, fmt.Sprintf("failed to check if svc exists: %s", err.Error()))
	} else if !exist {
		log.Fatalf("\"%s\" not found", svcName)
	}

	log.AddField("SVC", svcName)
}

func initAppAndCheckIfSvcExist(appName string, svcName string, svcAttr []string) {
	serviceType := "deployment"
	if len(svcAttr) > 0 {
		serviceType = svcAttr[0]
	}
	initApp(appName)
	CheckIfSvcExist(svcName, serviceType)
}

func initNocalhostDefaultApplicationAndWait() error {
	nsLock := nocalhost.NsLock(nameSpace)

	err := nsLock.Lock()
	defer nsLock.Unlock()
	if err != nil {
		return err
	}

	// double check
	if !nocalhost.CheckIfApplicationExist(app.DefaultNocalhostApplication, nameSpace) {
		log.Logf("Default virtual application in ns %s haven't init yet...", nameSpace)

		switch nocalhost.EstimateApplicationCounts(nameSpace) {
		case 1: // means we can move user's configurations to default app
			log.Logf("Init virtual application and copy unique app's configuration in ns %s ...", nameSpace)
			err := InitDefaultApplicationByFirstValid()
			if err != nil {
				return err
			}
		default:
			log.Logf("Init virtual application in ns %s ...", nameSpace)
			err := InitDefaultApplicationInCurrentNs()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func InitDefaultApplicationByFirstValid() error {
	appNeedToCopy := nocalhost.GetFirstApplication(nameSpace)
	if appNeedToCopy == app.DefaultNocalhostApplication {
		log.Error("Error while init %s, %s need to be created but has been created. ", appNeedToCopy, appNeedToCopy)
	}
	var err error

	defer func() {
		if err != nil {
			os.RemoveAll(nocalhost.GetAppDirUnderNs(app.DefaultNocalhostApplication, nameSpace))
		}
	}()

	err = utils.CopyDir(
		nocalhost.GetAppDirUnderNs(appNeedToCopy, nameSpace),
		nocalhost.GetAppDirUnderNs(app.DefaultNocalhostApplication, nameSpace),
	)
	if err != nil {
		return err
	}
	return nil
}

func InitDefaultApplicationInCurrentNs() error {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	baseDir := fp.NewFilePath(tmpDir)
	nocalhostDir := baseDir.RelOrAbs(app.DefaultGitNocalhostDir)
	err = nocalhostDir.Mkdir()
	if err != nil {
		return err
	}

	var cfg = ".default_config"

	err = nocalhostDir.RelOrAbs(cfg).WriteFile("name: nocalhost.default\nmanifestType: rawManifestLocal")
	if err != nil {
		return err
	}

	installFlags.Config = cfg
	installFlags.AppType = string(app.ManifestLocal)
	installFlags.LocalPath = baseDir.Abs()
	err = InstallApplication(app.DefaultNocalhostApplication)
	if err != nil {
		return err
	}

	return nil
}
