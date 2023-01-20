package controllers

import (
	"github.com/malike/mock-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *SampleKindReconciler) newSampleAppDeployment(name string, m *v1alpha1.SampleKind) *appsv1.Deployment {
	depLabel := getObjectLabels(m.Name)
	depLabel["appliedHash"] = generateHash(m.Spec.String())
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: m.Namespace,
			Labels:    depLabel,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: getSelectorLabels(m.Name),
			},
			Replicas: &m.Spec.Nodes,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getSelectorLabels(m.Name),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            name,
							Image:           m.Spec.Image.Repository + ":" + m.Spec.Image.Tag,
							ImagePullPolicy: m.Spec.Image.ImagePullPolicy,
							Ports: []corev1.ContainerPort{{
								Name:          "tcp",
								ContainerPort: m.Spec.ContainerPort,
							}},
						},
					},
					ImagePullSecrets: m.Spec.Image.ImagePullSecrets,
				},
			},
		},
	}

	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}
func (r *SampleKindReconciler) updateSampleAppDeployment(dep appsv1.Deployment, m *v1alpha1.SampleKind) appsv1.Deployment {
	depLabel := getObjectLabels(m.Name)
	depLabel["appliedHash"] = generateHash(m.Spec.String())
	dep.Labels = depLabel
	dep.Spec.Replicas = &m.Spec.Nodes
	return dep
}

func (r *SampleKindReconciler) newSampleAppService(name string, m *v1alpha1.SampleKind) *corev1.Service {
	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getServiceName(name),
			Namespace: m.Namespace,
			Labels:    getSelectorLabels(name),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp",
					Protocol:   corev1.ProtocolTCP,
					Port:       m.Spec.ServicePort,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: m.Spec.ContainerPort},
				},
			},
			Selector: getSelectorLabels(name),
		},
	}
	// Set instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}
