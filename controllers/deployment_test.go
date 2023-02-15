package controllers

import (
	samplekind "github.com/malike/mock-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

var _ = Describe("Deployment test", func() {

	const (
		name      = "deployment-test"
		namespace = "default"
	)

	Context("When SampleKind is created, Deployment is created", func() {

		It("allows deployment to be created and deleted", func() {

			By("set up deployment", func() {

				skMockOperator := &samplekind.SampleKind{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Spec: samplekind.SampleKindSpec{
						Nodes: 1,
					},
				}
				Expect(k8sClient.Create(ctx, skMockOperator)).Should(Succeed())

				EventuallyWithOffset(10, func() bool {
					smDeployment := &v1.Deployment{}
					err := k8sClient.Get(ctx, types.NamespacedName{Name: skMockOperator.Name, Namespace: skMockOperator.Namespace}, smDeployment)
					return err == nil
				}).WithTimeout(20 * time.Second).Should(BeTrue())

				//delete samplekind delete deployment
				Expect(k8sClient.Delete(ctx, skMockOperator)).To(Succeed())

			})

		})
	})

})
