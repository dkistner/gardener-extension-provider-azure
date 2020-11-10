// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package internal

import (
	azureapi "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure"

	"github.com/Azure/go-autorest/autorest"
)

// AzureAPIErrorNotFound tries to determine if an error is a resource not found error.
func AzureAPIErrorNotFound(err error) bool {
	switch err.(type) {
	case autorest.DetailedError:
		detailedErr := autorest.DetailedError(err.(autorest.DetailedError))
		if detailedErr.Response != nil && detailedErr.Response.StatusCode == 404 {
			return true
		}
	}
	return false
}

// IsVMORequired determines if VMO is required.
func IsVMORequired(infrastructureStatus *azureapi.InfrastructureStatus) bool {
	if infrastructureStatus.Zoned || len(infrastructureStatus.AvailabilitySets) > 0 {
		return false
	}
	return true
}
