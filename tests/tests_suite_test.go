package tests_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"kubevirt.io/kubevirt-ansible/tests"
	"kubevirt.io/qe-tools/pkg/ginkgo-reporters"
	ktests "kubevirt.io/kubevirt/tests"
)

func TestTests(t *testing.T) {
	RegisterFailHandler(ktests.KubevirtFailHandler)
	reporters := make([]Reporter, 0)
	if ginkgo_reporters.Polarion.Run {
		reporters = append(reporters, &ginkgo_reporters.Polarion)
	}
	if ginkgo_reporters.JunitOutput != "" {
		reporters = append(reporters, ginkgo_reporters.NewJunitReporter())
	}
	RunSpecsWithDefaultAndCustomReporters(t, "Tests Suite", reporters)
}

var _ = BeforeSuite(func() {
	tests.CreateNamespaces()
})

var _ = AfterSuite(func() {
	tests.RemoveNamespaces()
})
