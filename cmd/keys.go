/*
 * Copyright 2018-2019 The NATS Authors
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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nats-io/nsc/cmd/store"
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage keys for operators, accounts, and users",
}

func init() {
	GetRootCmd().AddCommand(keysCmd)
	keysCmd.AddCommand(createMigrateKeysCmd())
	keysCmd.AddCommand(createGetNKeysCmd())
}

func createMigrateKeysCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrates keystore to new layout, original keystore is preserved",
		RunE: func(cmd *cobra.Command, args []string) error {
			migration, err := store.KeysNeedMigration()
			if err != nil {
				return err
			}
			if !migration {
				cmd.Printf("keystore %q does not need migration\n", AbbrevHomePaths(store.GetKeysDir()))
				return nil
			}

			old, err := store.Migrate()
			if err != nil {
				return err
			}
			cmd.Printf("keystore %q was migrated - old store was renamed to %q - remove at your convenience\n",
				AbbrevHomePaths(store.GetKeysDir()),
				AbbrevHomePaths(old))

			return nil
		},
	}

	return cmd
}

func createGetNKeysCmd() *cobra.Command {
	var params GetKeysParams
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "prints the contents of the credentialsfile (.creds) for a specific user",
		RunE: func(cmd *cobra.Command, args []string) error {
			// cmd.Print("here:", params.Account)
			// cmd.Print("here:", store.GetKeysDir())
			filename := store.GetKeysDir() + "/creds/" + params.Operator + "/" + params.Account + "/" + params.User + ".creds"
			file, err := os.Open(filename)
			if err != nil {
				return errors.New("credsfile not found")
			}
			defer file.Close()

			content, err := ioutil.ReadFile(filepath.FromSlash(filename))
			if err != nil {
				return err
			}
			cmd.Print(string(content))
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if params.Operator == "" {
				return errors.New("operator flag may not be empty")
			}
			if params.Account == "" {
				return errors.New("account flag may not be empty")
			}
			if params.User == "" {
				return errors.New("user flag may not be empty")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.Operator, "operator", "", "", "export specified operator key")
	cmd.Flags().StringVarP(&params.Account, "account", "", "", "change account context to the named account")
	cmd.Flags().StringVarP(&params.User, "user", "", "", "export specified user key")

	return cmd
}

type GetKeysParams struct {
	Account  string
	User     string
	Operator string
}
