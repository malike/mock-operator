/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	appv1alpha1 "github.com/malike/mock-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// SampleKindReconciler reconciles a SampleKind object
type SampleKindReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.malike.kendeh.com,resources=samplekinds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.malike.kendeh.com,resources=samplekinds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.malike.kendeh.com,resources=samplekinds/finalizers,verbs=update

//+kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;delete;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete


func (r *SampleKindReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the SampleKind instance
	sampleApp := &appv1alpha1.SampleKind{}
	err := r.Get(ctx, req.NamespacedName, sampleApp)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("SampleKind resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object
		return ctrl.Result{}, err
	} else {
		log.V(1).Info("Detected existing SampleKind", " sampleApp.Name", sampleApp.Name)
	}

	// Check if the Deployment already exists, if not create a new one
	deployment := &appsv1.Deployment{}
	deploymentName := sampleApp.Name
	err = r.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: sampleApp.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new configmap
		deployment := r.newSampleAppDeployment(deploymentName, sampleApp)
		log.Info("Creating a new SampleApp", "SampleKind.Namespace", sampleApp.Namespace, "SampleKind.Name", sampleApp.Name)
		err = r.Create(ctx, deployment)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		//check if hash is outdated
		oldHashLabel := getHashLabel(deployment.Labels)
		if oldHashLabel != generateHash(sampleApp.Spec) {
			//outdated update deployment
			deployment := r.updateSampleAppDeployment(*deployment, sampleApp)
			err := r.Update(ctx, &deployment)
			if err != nil {
				//log erro
			}
		}
		log.V(1).Info("Detected existing SampleApp", " deployment.Name", sampleApp.Name)
	}

	service := &corev1.Service{}
	serviceName := getServiceName(deploymentName)
	err = r.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: sampleApp.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		// New service
		service = r.newSampleAppService(deploymentName, sampleApp)
		log.Info("Creating a new Service for SampleApp ", "Service.Namespace", service.Namespace, "Service.Name", service.Name)

		err = r.Create(ctx, service)
		if err != nil {
			//log failed to create
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		//log failed to create
		return ctrl.Result{}, err
	} else {
		log.V(1).Info("Detected existing Service", " Service.Name", service.Name)
	}

	// When finished still reconcile periodically
	return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SampleKindReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.SampleKind{}).
		Complete(r)
}
