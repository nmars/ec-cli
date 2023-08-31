#!/usr/bin/env bash
# Copyright The Enterprise Contract Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0

# The EC golden container, see https://github.com/enterprise-contract/golden-container/
IMAGE=${IMAGE:-"ghcr.io/enterprise-contract/golden-container:latest"}
IDENTITY_REGEXP=${IDENTITY_REGEXP:-"https:\/\/github\.com\/(slsa-framework\/slsa-github-generator|enterprise-contract\/golden-container)\/"}
IDENTITY_ISSUER=${IDENTITY_ISSUER:-"https://token.actions.githubusercontent.com"}

# Festoji, see https://github.com/lcarva/festoji
#IMAGE=${IMAGE:-"quay.io/lucarval/festoji:latest"}
#IDENTITY_REGEXP=${IDENTITY_REGEXP:-"https:\/\/github\.com\/(slsa-framework\/slsa-github-generator|lcarva\/festoji)\/"}
#IDENTITY_ISSUER=${IDENTITY_ISSUER:-"https://token.actions.githubusercontent.com"}

# Todo: Use a useful policy here
POLICY=""

OPTS=${1:-}
MAIN_GO=$(git rev-parse --show-toplevel)/main.go
go run $MAIN_GO validate image --image "${IMAGE}" \
  --policy "${POLICY}" \
  --certificate-identity-regexp ${IDENTITY_REGEXP} \
  --certificate-oidc-issuer ${IDENTITY_ISSUER} \
  --info $OPTS | yq -P
