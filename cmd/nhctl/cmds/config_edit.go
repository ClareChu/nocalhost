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
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"nocalhost/internal/nhctl/profile"

	"nocalhost/pkg/nhctl/log"
)

type ConfigEditFlags struct {
	CommonFlags
	Content   string
	AppConfig bool
}

var configEditFlags = ConfigEditFlags{}

func init() {
	configEditCmd.Flags().StringVarP(&configEditFlags.SvcName, "deployment", "d", "", "k8s deployment which your developing service exists")
	configEditCmd.Flags().StringVarP(&configEditFlags.Content, "content", "c", "", "base64 encode json content")
	configEditCmd.Flags().BoolVar(&configEditFlags.AppConfig, "app-config", false, "edit application config")
	configCmd.AddCommand(configEditCmd)
}

var configEditCmd = &cobra.Command{
	Use:   "edit [Name]",
	Short: "edit service config",
	Long:  "edit service config",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.Errorf("%q requires at least 1 argument\n", cmd.CommandPath())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		configEditFlags.AppName = args[0]

		initApp(configEditFlags.AppName)

		if len(configEditFlags.Content) == 0 {
			log.Fatal("--content required")
		}

		bys, err := base64.StdEncoding.DecodeString(configEditFlags.Content)
		if err != nil {
			log.Fatalf("--content must be a valid base64 string: %s", err.Error())
		}

		// set application config, plugin do not provide services struct, update application config only
		if configEditFlags.AppConfig {
			applicationConfig := &profile.ApplicationConfig{}
			err = json.Unmarshal(bys, applicationConfig)
			if err != nil {
				log.Fatalf("fail to unmarshal content: %s", err.Error())
			}
			// update config
			// update profile
			if err := nocalhostApp.SaveAppProfileV2(applicationConfig); err != nil {
				log.FatalE(err, "fail to save app config")
			}
			return
		}
		// Deprecated: V1
		//svcConfig := &app.ServiceDevOptions{}
		//err = json.Unmarshal(bys, svcConfig)
		//if err != nil {
		//	log.Fatalf("fail to unmarshal content: %s", err.Error())
		//}
		//err = nocalhostApp.SaveSvcProfile(configEditFlags.SvcName, svcConfig)
		//if err != nil {
		//	log.FatalE(err, "fail to save svc config")
		//}

		svcConfig := &profile.ServiceConfigV2{}

		CheckIfSvcExist(configEditFlags.SvcName)

		err = json.Unmarshal(bys, svcConfig)
		if err != nil {
			log.Fatalf("fail to unmarshal content: %s", err.Error())
		}
		err = nocalhostApp.SaveSvcProfileV2(configEditFlags.SvcName, svcConfig)
		if err != nil {
			log.FatalE(err, "fail to save svc config")
		}
	},
}
