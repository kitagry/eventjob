package v1_test

import (
	"testing"

	v1 "github.com/kitagry/eventjob/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEventJobTrigger_Match(t *testing.T) {
	tests := map[string]struct {
		trigger v1.EventJobTrigger
		event   corev1.Event
		expect  bool
	}{
		"trigger match CronJob complete": {
			trigger: v1.EventJobTrigger{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "batch/v1",
					Kind:       "CronJob",
				},
				Type: v1.CompleteTrigger,
			},
			event: corev1.Event{
				InvolvedObject: corev1.ObjectReference{
					APIVersion: "batch/v1",
					Kind:       "CronJob",
				},
				Message: "aw completed job: cronjob-27210675, status: Complete",
			},
			expect: true,
		},
		"trigger don't match CronJob failed": {
			trigger: v1.EventJobTrigger{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "batch/v1",
					Kind:       "CronJob",
				},
				Type: v1.CompleteTrigger,
			},
			event: corev1.Event{
				InvolvedObject: corev1.ObjectReference{
					APIVersion: "batch/v1",
					Kind:       "CronJob",
				},
				Message: "aw completed job: cronjob-27210675, status: Failed",
			},
			expect: false,
		},
		"trigger match CronJob failed": {
			trigger: v1.EventJobTrigger{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "batch/v1",
					Kind:       "CronJob",
				},
				Type: v1.FailedTrigger,
			},
			event: corev1.Event{
				InvolvedObject: corev1.ObjectReference{
					APIVersion: "batch/v1",
					Kind:       "CronJob",
				},
				Message: "aw completed job: cronjob-27210675, status: Failed",
			},
			expect: true,
		},
	}

	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			matched := tt.trigger.Match(tt.event)
			if matched != tt.expect {
				t.Errorf("EventJobTrigger expect %t, got %t", tt.expect, matched)
			}
		})
	}
}
