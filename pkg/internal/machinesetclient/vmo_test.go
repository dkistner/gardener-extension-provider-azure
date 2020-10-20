// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package machinesetclient_test

import (
	azureapi "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gardener/gardener-extension-provider-azure/pkg/internal/machinesetclient"
)

var _ = Describe("VMO", func() {
	var (
		infrastructureStatus        *azureapi.InfrastructureStatus
		defaultInfrastructureStatus = azureapi.InfrastructureStatus{
			Zoned: false,
		}
	)

	BeforeEach(func() {
		infrastructureStatus = defaultInfrastructureStatus.DeepCopy()
	})

	Describe("#IsVMORequired", func() {
		It("should require an VMO", func() {
			Expect(IsVMORequired(infrastructureStatus)).To(BeTrue())
		})

		It("should not require VMO for zoned cluster", func() {
			infrastructureStatus.Zoned = true
			Expect(IsVMORequired(infrastructureStatus)).To(BeFalse())
		})

		It("should not require VMO for a cluster with primary availabilityset (non zoned)", func() {
			infrastructureStatus.AvailabilitySets = []azureapi.AvailabilitySet{
				{
					ID:      "/my/azure/availabilityset/id",
					Name:    "my-availabilityset",
					Purpose: "nodes",
				},
			}
			Expect(IsVMORequired(infrastructureStatus)).To(BeFalse())
		})
	})
})
