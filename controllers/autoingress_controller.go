/*
Copyright 2022.

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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1 "github.com/bryant-rh/auto-ingress-operator/api/v1"
	"github.com/bryant-rh/auto-ingress-operator/controllers/helper"
	"github.com/bryant-rh/auto-ingress-operator/controllers/util"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// AutoIngressReconciler reconciles a AutoIngress object
type AutoIngressReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.ingress.com,resources=autoingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.ingress.com,resources=autoingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.ingress.com,resources=autoingresses/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources="services",verbs=get;list;watch
//+kubebuilder:rbac:groups="networking.k8s.io",resources="ingresses",verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AutoIngress object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *AutoIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.WithValues("func", "Reconcile")
	logger.Info("进入 Reconcile")
	defer logger.Info("退出 Reconcile")

	// TODO(user): your logic here

	instance := &appsv1.AutoIngress{}

	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 删除
	if !instance.DeletionTimestamp.IsZero() {
		autoIngressSet.Remove(*instance)

		return ctrl.Result{}, nil
	}

	// 保存
	// if len(instance.Spec.ServicePrefixes) == 0 {
	// 	instance.Spec.ServicePrefixes = []string{"web-", "srv-"}
	// }
	autoIngressSet.Add(*instance)

	r.ReconcileServices(ctx, instance)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.AutoIngress{}).
		Watches(
			&source.Kind{
				Type: &corev1.Service{},
			},
			handler.Funcs{
				CreateFunc: r.onCreateService,
			},
		).
		Watches(
			&source.Kind{
				Type: &netv1.Ingress{},
			},
			handler.Funcs{
				DeleteFunc: r.onIngressDelete,
			},
		).
		Complete(r)
}

//onCreateService
func (r *AutoIngressReconciler) onCreateService(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	logger := log.FromContext(context.TODO())
	logger.WithValues("func", "onCreateService")
	logger.Info("新 Service 被创建")

	svc := r.getService(e.Object)
	if svc == nil {
		return
	}

	ctx := context.TODO()
	for _, op := range autoIngressSet.List() {
		r.HandleIngress(ctx, op, svc)
	}
}

//onIngressDelete
func (r *AutoIngressReconciler) onIngressDelete(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
	ctx := context.TODO()
	logger := log.FromContext(ctx)
	logger.WithValues("func", "onIngressDelete")
	logger.Info(fmt.Sprintf("Ingress 删除: %s.%s", e.Object.GetName(), e.Object.GetNamespace()))

	svc := &corev1.Service{}
	for _, owner := range e.Object.GetOwnerReferences() {
		if owner.Kind == "Service" {
			key := types.NamespacedName{
				Namespace: e.Object.GetNamespace(),
				Name:      owner.Name,
			}
			err := r.Client.Get(ctx, key, svc)
			// 如果是删除 svc 早成的 删除 ingress 就略过
			if apierrors.IsNotFound(err) {
				return
			}
		}
	}

	for _, op := range autoIngressSet.List() {
		r.HandleIngress(ctx, op, svc)
	}

}

//getService
func (r *AutoIngressReconciler) getService(e client.Object) *corev1.Service {

	key := r.objectKey(e)
	svc := &corev1.Service{}

	err := r.Client.Get(context.TODO(), key, svc)
	if err != nil {
		return nil
	}

	return svc
}

//ReconcileServices
func (r *AutoIngressReconciler) ReconcileServices(ctx context.Context, op *appsv1.AutoIngress) {
	logger := log.FromContext(ctx)
	logger.WithValues("func", "ReconcileServices")

	svcs := &corev1.ServiceList{}
	err := r.Client.List(ctx, svcs)
	if err != nil {
		logger.Error(err, "list services failed")
		return
	}

	for _, svc := range svcs.Items {
		r.HandleIngress(ctx, *op, &svc)
	}
}

//HandleIngress
func (r *AutoIngressReconciler) HandleIngress(ctx context.Context, op appsv1.AutoIngress, svc *corev1.Service) {
	logger := log.FromContext(ctx)
	logger.WithValues("func", "HandleIngress")

	if len(op.Spec.ServicePrefixes) != 0 {
		if !util.IsValidServcieName(svc.Name, op.Spec.ServicePrefixes) {
			return
		}
	}
	if len(op.Spec.NameSpaces) != 0 {
		if !util.IsValidNsName(svc.Namespace, op.Spec.NameSpaces) {
			return
		}
	}
	if len(op.Spec.NameSpaceBlackList) != 0 {
		if util.IsValidNsName(svc.Namespace, op.Spec.NameSpaceBlackList) {
			return
		}
	}

	ing := helper.NewIngress(op, svc)
	_ing := helper.NewIngress(op, svc)
	action := "create"

	if r.isExistInK8s(_ing) {
		action = "update"
		ing.SetResourceVersion(_ing.ResourceVersion)
	}

	_ = controllerutil.SetOwnerReference(svc, ing, r.Scheme)
	// _ = controllerutil.SetOwnerReference(&op, ing, r.Scheme)

	err := r.HandleObject(ctx, ing, action)
	if err != nil {
		logger.Error(err, fmt.Sprintf("handle(%s) ingress (%s.%s) failed", action, ing.Name, ing.Namespace))

		return
	}

	logger.Info(fmt.Sprintf("handle(%s) ingress (%s.%s) success", action, ing.Name, ing.Namespace))
}

//HandleObject
func (r *AutoIngressReconciler) HandleObject(ctx context.Context, obj client.Object, action string) error {

	switch action {
	case "update":
		return r.Client.Update(ctx, obj)
	case "create":
		return r.Client.Create(ctx, obj)
	}

	return nil
}

// isExistInK8s 检查对象是否在 k8s 中存在
func (r *AutoIngressReconciler) isExistInK8s(obj client.Object) bool {

	key := r.objectKey(obj)
	err := r.Client.Get(context.TODO(), key, obj)

	return err == nil
}

//objectKey
func (r *AutoIngressReconciler) objectKey(e client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: e.GetNamespace(),
		Name:      e.GetName(),
	}
}
