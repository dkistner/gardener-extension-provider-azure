package worker_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gardener/gardener-extension-provider-azure/pkg/controller/worker"

	"github.com/gardener/gardener/extensions/pkg/controller/common"
	"github.com/gardener/gardener/extensions/pkg/controller/worker/genericactuator"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	mockclient "github.com/gardener/gardener/pkg/mock/controller-runtime/client"
)

var _ = Describe("MachinesDependencies", func() {
	var (
		ctrl *gomock.Controller
		c    *mockclient.MockClient

		scheme, decoder = getSchemaAndDecoder()

		workerDelegate genericactuator.WorkerDelegate
		w              *extensionsv1alpha1.Worker
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		c = mockclient.NewMockClient(ctrl)

		w = getDefaultWorker("shoot--foobar--azure", nil)
		expectGetSecretCallToWork(c, w)
		workerDelegate, _ = NewWorkerDelegate(common.NewClientContext(c, scheme, decoder), nil, "", w, nil)
	})

	Context("#DeployMachineDependencies", func() {
		It("should return no error", func() {
			err := workerDelegate.DeployMachineDependencies(context.TODO())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("#CleanupMachineDependencies", func() {
		It("should return no error", func() {
			err := workerDelegate.CleanupMachineDependencies(context.TODO())
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
