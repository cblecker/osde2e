package operators

import (
	"context"
	"fmt"
	"net/http"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/osde2e/pkg/common/alert"
	"github.com/openshift/osde2e/pkg/common/config"
	"github.com/openshift/osde2e/pkg/common/helper"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	osdMetricsExporterTestPrefix   = "[Suite: operators] [OSD] OSD Metrics Exporter"
	osdMetricsExporterBasicTest    = osdMetricsExporterTestPrefix + " Basic Test"
)

func init() {
	alert.RegisterGinkgoAlert(osdMetricsExporterBasicTest, "SD_SREP", "Arjun Naik", "sd-cicd-alerts", "sd-cicd@redhat.com", 4)
}

var _ = ginkgo.Describe(osdMetricsExporterBasicTest, func() {
	var (
		operatorNamespace = "openshift-osd-metrics"
		operatorName      = "osd-metrics-exporter"
		clusterRoles      = []string{
			"osd-metrics-exporter",
		}
		clusterRoleBindings = []string{
			"osd-metrics-exporter",
		}
		servicePort = 8383
	)
	h := helper.New()
	checkClusterServiceVersion(h, operatorNamespace, operatorName)
	checkDeployment(h, operatorNamespace, operatorName, 1)
	checkClusterRoles(h, clusterRoles, true)
	checkClusterRoleBindings(h, clusterRoleBindings, true)
	checkService(h, operatorNamespace, operatorName, servicePort)
	checkUpgrade(helper.New(), operatorNamespace, operatorName, operatorName, "osd-metrics-exporter-registry")
})

func checkService(h *helper.H, namespace string, name string, port int) {
	pollTimeout := viper.GetFloat64(config.Tests.PollingTimeout)
	serviceEndpoint := fmt.Sprintf("http://%s.%s:%d/metrics", name, namespace, port)
	ginkgo.Context("service", func() {
		ginkgo.It(
			"should exist",
			func() {
				Eventually(func() bool {
					_, err := h.Kube().CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
					if err != nil {
						return false
					}
					return true
				}, "30m", "1m").Should(BeTrue())
			},
			pollTimeout,
		)
		ginkgo.It(
			"should return response",
			func() {
				Eventually(func() (*http.Response, error) {
					response, err := http.Get(serviceEndpoint)
					return response, err
				}, "30m", "1m").Should(HaveHTTPStatus(http.StatusOK))
			},
			pollTimeout,
		)
	})
}