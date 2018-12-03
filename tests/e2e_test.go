package tests_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	"kubevirt.io/kubevirt-ansible/tests"
)

// template parameters
const (
	pvcEPHTTPNOAUTHURL = "https://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
	invalidPVCURL      = "https://noneexist.com"
	pvcName            = "golden-pvc"
	pvcName1           = "golden-pvc1"
	vmName             = "test-vm"
	vmAPIVersion       = "kubevirt.io/v1alpha2"
	rawPVCFilePath     = "tests/manifests/golden-pvc.yml"
	rawVMFilePath      = "tests/manifests/test-vm.yml"
)

var _ = Describe("Importing and starting a VM using CDI", func() {
	var dstPVCFilePath, dstVMFilePath, newPVCName, url string

	BeforeEach(func() {
		var ok bool
		dstPVCFilePath = "/tmp/test-pvc.json"
		dstVMFilePath = "/tmp/test-vm.json"
		newPVCName = pvcName
		url, ok = os.LookupEnv("STREAM_IMAGE_URL")
		if !ok {
			url = pvcEPHTTPNOAUTHURL
		}
	})

	JustBeforeEach(func() {
		tests.ProcessTemplateWithParameters(rawPVCFilePath, dstPVCFilePath, "PVC_NAME="+newPVCName, "EP_URL="+url)
		tests.CreateResourceWithFilePathTestNamespace(dstPVCFilePath)
	})

	Context("PVC with valid image url", func() {

		It("will succeed", func() {
			tests.WaitUntilResourceReadyByNameTestNamespace("pvc", pvcName, "-o=jsonpath='{.metadata.annotations}'", "pv.kubernetes.io/bind-completed:yes")
			tests.ProcessTemplateWithParameters(rawVMFilePath, dstVMFilePath, "VM_NAME="+vmName, "PVC_NAME="+pvcName, "VM_APIVERSION="+vmAPIVersion)
			tests.CreateResourceWithFilePathTestNamespace(dstVMFilePath)
			tests.WaitUntilResourceReadyByNameTestNamespace("vmi", vmName, "-o=jsonpath='{.status.phase}'", "Running")
		})
	})

	Context("PVC with invalid image url", func() {
		BeforeEach(func() {
			newPVCName = pvcName1
			url = invalidPVCURL
		})

		It("will be failed because the PVC should become failed", func() {
			tests.WaitUntilResourceReadyByLabelTestNamespace("pod", tests.CDI_LABEL_SELECTOR, "", "CrashLoopBackOff")
		})
	})

})
