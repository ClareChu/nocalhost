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
	"nocalhost/internal/nhctl/envsubst"
	"nocalhost/internal/nhctl/fp"
	"nocalhost/pkg/nhctl/log"
)

var renderOps = &RenderOps{}

type RenderOps struct {
	envPath string
}

func init() {
	renderCmd.Flags().StringVarP(&renderOps.envPath, "env path", "e", "", "the env file for render injection")
	rootCmd.AddCommand(renderCmd)
}

var renderCmd = &cobra.Command{
	Use:   "render [NAME]",
	Short: "Render the file for debugging",
	Long:  `Render the file for debugging`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.Errorf("%q requires at least 1 argument\n", cmd.CommandPath())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		render, err := envsubst.Render(fp.NewFilePath(args[0]), fp.NewFilePath(renderOps.envPath))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(render)
	},
}
