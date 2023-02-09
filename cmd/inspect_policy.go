// Copyright 2022 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

// Define the `ec inspect policy` command
package cmd

import (
	"encoding/json"

	hd "github.com/MakeNowJust/heredoc"
	"github.com/open-policy-agent/opa/ast"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/hacbs-contract/ec-cli/internal/opa"
	"github.com/hacbs-contract/ec-cli/internal/policy/source"
	"github.com/hacbs-contract/ec-cli/internal/utils"
)

func inspectPolicyCmd() *cobra.Command {
	var (
		sourceUrls   []string
		destDir      string
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "policy --source <source-url>",
		Short: "Read policies from source urls and show information about the rules inside them",

		Long: hd.Doc(`
			Read policies from a source url and show information about the rules inside them.

			This fetches policy sources similar to the 'ec fetch policy' command, but once
			the policy is fetched the equivalent of 'opa inspect' is run against the
			downloaded policies.

			This can be used to extract information about each rule in the policy source,
			including the rule annotations which include the rule's title and description
			and custom fields used by ec to filter the results produced by conftest.

			Note that this command is not typically required to verify the Enterprise
			Contract. It has been made available for troubleshooting and debugging purposes.
		`),

		Example: hd.Doc(`
			Print a list of rules and their descriptions from the latest Stonesoup release policy:

			  ec inspect policy --source quay.io/hacbs-contract/ec-release-policy

			Display details about the latest Stonesoup release policy in json format:

			  ec inspect policy --source quay.io/hacbs-contract/ec-release-policy -o json | jq
		`),

		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if destDir == "" {
				workDir, err := utils.CreateWorkDir(afero.NewOsFs())
				if err != nil {
					log.Debug("Failed to create work dir!")
					return err
				}
				destDir = workDir
			}

			allResults := make(map[string][]*ast.AnnotationsRef)
			for _, url := range sourceUrls {
				s := &source.PolicyUrl{Url: url, Kind: source.PolicyKind}

				// Download
				policyDir, err := s.GetPolicy(cmd.Context(), destDir, false)
				if err != nil {
					return err
				}

				// Inspect
				result, err := opa.InspectDir(afero.NewOsFs(), policyDir)
				if err != nil {
					return err
				}

				// Collect results
				allResults[s.PolicyUrl()] = result
			}

			out := cmd.OutOrStdout()
			if outputFormat == "json" {
				return json.NewEncoder(out).Encode(allResults)
			}
			return opa.OutputText(out, allResults)
		},
	}

	cmd.Flags().StringArrayVarP(&sourceUrls, "source", "s", []string{}, "policy source url. multiple values are allowed")
	cmd.Flags().StringVarP(&destDir, "dest", "d", "", "use the specified destination directory to download the policy. if not set, a temporary directory will be used")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "output format, either json or text.")

	if err := cmd.MarkFlagRequired("source"); err != nil {
		panic(err)
	}

	return cmd
}