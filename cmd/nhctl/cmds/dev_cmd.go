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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"nocalhost/pkg/nhctl/log"
	"os"
)

type DevCommandType string

const (
	buildCommand          DevCommandType = "build"
	runCommand            DevCommandType = "run"
	debugCommand          DevCommandType = "debug"
	hotReloadRunCommand   DevCommandType = "hotReloadRun"
	hotReloadDebugCommand DevCommandType = "hotReloadDebug"
)

var commandType string
var container string

func init() {
	devCmdCmd.Flags().StringVarP(&deployment, "deployment", "d", "", "K8s deployment which your developing service exists")
	devCmdCmd.Flags().StringVarP(&container, "container", "c", "", "which container of pod to run command")
	devCmdCmd.Flags().StringVar(&commandType, "dev-command-type", "", fmt.Sprintf("Dev command type can be: %s, %s, %s, %s, %s", buildCommand, runCommand, debugCommand, hotReloadRunCommand, hotReloadDebugCommand))
	debugCmd.AddCommand(devCmdCmd)
}

var devCmdCmd = &cobra.Command{
	Use:   "cmd [NAME]",
	Short: "Run cmd in dev container",
	Long:  `Run cmd in dev container`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.Errorf("%q requires at least 1 argument\n", cmd.CommandPath())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if commandType == "" {
			log.Fatal("--dev-command-type mush be specified")
		}
		applicationName := args[0]
		initAppAndCheckIfSvcExist(applicationName, deployment, nil)
		if !nocalhostApp.CheckIfSvcIsDeveloping(deployment) {
			log.Fatalf("%s is not in DevMode", deployment)
		}
		profile := nocalhostApp.GetSvcProfileV2(deployment)
		if profile == nil {
			log.Fatal("Failed to get service profile")
			os.Exit(1)
		}

		if profile.GetContainerDevConfigOrDefault(container) == nil || profile.GetContainerDevConfigOrDefault(container).Command == nil {
			log.Fatalf("%s command not defined", commandType)
		}
		var targetCommand []string
		switch commandType {
		case string(buildCommand):
			targetCommand = profile.GetContainerDevConfigOrDefault(container).Command.Build
		case string(runCommand):
			targetCommand = profile.GetContainerDevConfigOrDefault(container).Command.Run
		case string(debugCommand):
			targetCommand = profile.GetContainerDevConfigOrDefault(container).Command.Debug
		case string(hotReloadDebugCommand):
			targetCommand = profile.GetContainerDevConfigOrDefault(container).Command.HotReloadDebug
		case string(hotReloadRunCommand):
			targetCommand = profile.GetContainerDevConfigOrDefault(container).Command.HotReloadRun
		default:
			log.Fatalf("%s is not supported", commandType)

		}
		if len(targetCommand) == 0 {
			log.Fatalf("%s command not defined", commandType)
		}

		err := nocalhostApp.Exec(deployment, container, targetCommand)
		if err != nil {
			log.Fatalf("Failed to exec : %s", err.Error())
		}
	},
}
