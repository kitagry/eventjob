/*
Copyright 2021.

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
	"math/rand"
	"sort"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	batchv1 "github.com/kitagry/eventjob/api/v1"
	k8sbatchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	eventTriggerField = ".spec.trigger"
)

// EventJobReconciler reconciles a EventJob object
type EventJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=batch.kitagry.github.io,resources=eventjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.kitagry.github.io,resources=eventjobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.kitagry.github.io,resources=eventjobs/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EventJob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *EventJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logging := log.FromContext(ctx)

	// your logic here
	var eventJob batchv1.EventJob
	err := r.Get(ctx, req.NamespacedName, &eventJob)
	if err != nil && errors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		logging.Error(err, "Failed to get EventJob")
		return ctrl.Result{}, err
	}

	var events corev1.EventList
	err = r.List(ctx, &events, &client.ListOptions{
		Namespace: eventJob.Namespace,
	})
	if err != nil {
		logging.Error(err, "Failed to get EventJob")
		return ctrl.Result{}, err
	}

	newStatus := batchv1.EventJobStatus{}
	ok, eventTime := shouldCreate(events, eventJob)
	if !ok {
		return ctrl.Result{}, nil
	}
	newStatus.LastScheduledTime = eventTime

	randstr, err := makeRandomStr(5)
	if err != nil {
		logging.Error(err, "Failed to make random string")
		return ctrl.Result{}, err
	}
	job := k8sbatchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: eventJob.GetNamespace(),
			Name:      fmt.Sprintf("%s-%s", eventJob.GetName(), randstr),
		},
		Spec: eventJob.Spec.JobTemplate.Spec,
	}

	err = ctrl.SetControllerReference(&eventJob, &job, r.Scheme)
	if err != nil {
		logging.Error(err, "failed to SetControllerReference")
		return ctrl.Result{}, err
	}

	err = r.Create(ctx, &job)
	if err != nil {
		logging.Error(err, "Failed to create job")
		return ctrl.Result{}, err
	}

	eventJob.Status = newStatus
	err = r.Status().Update(ctx, &eventJob)
	if err != nil {
		logging.Error(err, "Failed to create job")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func shouldCreate(events corev1.EventList, eventJob batchv1.EventJob) (bool, metav1.Time) {
	sort.Slice(events.Items, func(i, j int) bool {
		ei := events.Items[i]
		ej := events.Items[j]
		return ei.CreationTimestamp.After(ej.CreationTimestamp.Time)
	})

	for _, e := range events.Items {
		if e.CreationTimestamp.Before(&eventJob.Status.LastScheduledTime) {
			return false, metav1.Time{}
		}

		if eventJob.Spec.Trigger.Match(e) && e.CreationTimestamp.After(eventJob.Status.LastScheduledTime.Time) {
			return true, e.CreationTimestamp
		}
	}
	return false, metav1.Time{}
}

func makeRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("can't rand.Read: %w", err)
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EventJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &batchv1.EventJob{}, eventTriggerField, func(rawObj client.Object) []string {
		eventJob := rawObj.(*batchv1.EventJob)
		if eventJob.Spec.Trigger.APIVersion == "" || eventJob.Spec.Trigger.Kind == "" {
			return nil
		}
		return []string{formatTrigger(eventJob.Spec.Trigger.APIVersion, eventJob.Spec.Trigger.Kind)}
	}); err != nil {
		return err
	}

	notReconcileOption := builder.WithPredicates(predicate.ResourceVersionChangedPredicate{
		Funcs: predicate.Funcs{
			CreateFunc:  func(event.CreateEvent) bool { return false },
			UpdateFunc:  func(event.UpdateEvent) bool { return false },
			DeleteFunc:  func(event.DeleteEvent) bool { return false },
			GenericFunc: func(event.GenericEvent) bool { return false },
		},
	})
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.EventJob{}, notReconcileOption).
		Owns(&k8sbatchv1.Job{}, notReconcileOption).
		Watches(&source.Kind{Type: &corev1.Event{}}, handler.EnqueueRequestsFromMapFunc(r.findObjectsForEvent), builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		Complete(r)
}

func (r *EventJobReconciler) findObjectsForEvent(object client.Object) []reconcile.Request {
	event := object.(*corev1.Event)
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(eventTriggerField, formatTrigger(event.InvolvedObject.APIVersion, event.InvolvedObject.Kind)),
		Namespace:     event.Namespace,
	}

	var attachedEventJobs batchv1.EventJobList
	err := r.List(context.Background(), &attachedEventJobs, listOps)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, 0, len(attachedEventJobs.Items))
	for _, eventJob := range attachedEventJobs.Items {
		if !eventJob.Spec.Trigger.Match(*event) {
			continue
		}
		requests = append(requests, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: eventJob.Namespace,
				Name:      eventJob.Name,
			},
		})
	}

	return requests
}

func formatTrigger(apiVersion, kind string) string {
	return fmt.Sprintf("%s/%s", apiVersion, kind)
}
