/*
 * Copyright 2018 The NATS Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nats-io/nsc/cmd/store"
	"github.com/spf13/cobra"
	"github.com/xlab/tablewriter"
)

func envSet(varName string) string {
	return yn(os.Getenv(varName) != "")
}

func yn(v bool) string {
	if v {
		return "Yes"
	}
	return "No"

}

func createEnvCmd() *cobra.Command {
	var params SetContextParams
	cmd := &cobra.Command{
		Use:           "env",
		Short:         fmt.Sprintf("Prints and manage the %s environment", GetToolName()),
		Args:          MaxArgs(0),
		SilenceErrors: false,
		SilenceUsage:  false,
		Example:       "env",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := params.Run(cmd); err != nil {
				return err
			}
			switch OutputFormatFlag {
			case "table":
				params.PrintEnv(cmd)
			case "json":
				params.PrintEnvJSON(cmd)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.StoreRoot, "store", "s", "", "set store directory")
	cmd.Flags().StringVarP(&params.Operator, "operator", "o", "", "set operator name")
	cmd.Flags().StringVarP(&params.Account, "account", "a", "", "set account name")

	return cmd
}

func init() {
	GetRootCmd().AddCommand(createEnvCmd())
}

type SetContextParams struct {
	StoreRoot string
	Operator  string
	Account   string
}

func (p *SetContextParams) Run(cmd *cobra.Command) error {
	if *p == (SetContextParams{}) {
		// no edits
		return nil
	}
	current := GetConfig()
	if err := current.ContextConfig.Update(p.StoreRoot, p.Operator, p.Account); err != nil {
		return err
	}
	return current.Save()
}

func (p *SetContextParams) PrintEnv(cmd *cobra.Command) {
	conf := GetConfig()
	table := tablewriter.CreateTable()
	table.UTF8Box()
	table.AddTitle("NSC Environment")
	table.AddHeaders("Setting", "Set", "Effective Value")
	table.AddRow("$"+store.NKeysPathEnv, envSet(store.NKeysPathEnv), AbbrevHomePaths(store.GetKeysDir()))
	table.AddRow("$"+homeEnv, envSet(homeEnv), AbbrevHomePaths(toolHome))
	table.AddRow("Config", "", AbbrevHomePaths(conf.configFile()))
	table.AddSeparator()
	r := conf.StoreRoot
	if r == "" {
		r = "Not Set"
	}
	table.AddRow("From CWD", "", yn(GetCwdCtx() != nil))
	table.AddRow("Stores Dir", "", AbbrevHomePaths(r))
	table.AddRow("Default Operator", "", conf.Operator)
	table.AddRow("Default Account", "", conf.Account)
	cmd.Println(table.Render())
}

type Environment struct {
	NKeysPath       string
	NSCHome         string
	Config          string
	FromCWD         bool
	StoresDir       string
	DefaultOperator string
	DefaultAccount  string
}

func (p *SetContextParams) PrintEnvJSON(cmd *cobra.Command) {
	conf := GetConfig()

	environment := &Environment{
		NKeysPath:       store.GetKeysDir(),
		NSCHome:         toolHome,
		Config:          conf.configFile(),
		FromCWD:         GetCwdCtx() != nil,
		StoresDir:       conf.StoreRoot,
		DefaultOperator: conf.Operator,
		DefaultAccount:  conf.Account,
	}
	data, _ := json.Marshal(environment)
	cmd.Println(string(data))
}
