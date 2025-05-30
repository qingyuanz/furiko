/*
 * Copyright 2022 The Furiko Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package validation_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	fakeclock "k8s.io/utils/clock/testing"
	"k8s.io/utils/pointer"

	configv1alpha1 "github.com/furiko-io/furiko/apis/config/v1alpha1"
	"github.com/furiko-io/furiko/apis/execution/v1alpha1"
	"github.com/furiko-io/furiko/pkg/execution/util/jobconfig"
	"github.com/furiko-io/furiko/pkg/execution/validation"
	"github.com/furiko-io/furiko/pkg/runtime/controllercontext"
	"github.com/furiko-io/furiko/pkg/runtime/controllercontext/mock"
	"github.com/furiko-io/furiko/pkg/utils/testutils"
)

const (
	mockStartTime = "2021-02-09T04:06:09Z"
)

var (
	startTime = testutils.Mkmtimep(mockStartTime)

	objectMetaJobConfig = metav1.ObjectMeta{
		Name:      "jobconfig-sample",
		Namespace: metav1.NamespaceDefault,
		UID:       "7172ef3a-7754-4e73-a99a-d7e8a8e9fe5b",
	}

	ownerReferences = []metav1.OwnerReference{
		{
			APIVersion:         v1alpha1.GroupVersion.String(),
			Kind:               v1alpha1.KindJobConfig,
			Name:               objectMetaJobConfig.Name,
			UID:                objectMetaJobConfig.UID,
			Controller:         pointer.Bool(true),
			BlockOwnerDeletion: pointer.Bool(true),
		},
	}

	objectMetaJob = metav1.ObjectMeta{
		Name:      "job",
		Namespace: metav1.NamespaceDefault,
	}

	objectMetaJobWithNoLabelUID = metav1.ObjectMeta{
		Name:            "job",
		Namespace:       metav1.NamespaceDefault,
		OwnerReferences: ownerReferences,
	}

	objectMetaJobWithAllReferences = metav1.ObjectMeta{
		Name:      "job",
		Namespace: metav1.NamespaceDefault,
		Labels: map[string]string{
			jobconfig.LabelKeyJobConfigUID: string(objectMetaJobConfig.UID),
		},
		OwnerReferences: ownerReferences,
	}

	jobTemplateSpecBasic = v1alpha1.JobTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"labels.furiko.io/custom-label": "123",
			},
		},
		Spec: v1alpha1.JobTemplate{
			TaskTemplate: v1alpha1.TaskTemplate{
				Pod: &podTemplateSpecBasic,
			},
			TaskPendingTimeoutSeconds: pointer.Int64(1800),
			MaxAttempts:               pointer.Int64(5),
			RetryDelaySeconds:         pointer.Int64(60),
		},
	}

	jobTemplateSpecContainerUpdatedCommand = v1alpha1.JobTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"labels.furiko.io/custom-label": "123",
			},
		},
		Spec: v1alpha1.JobTemplate{
			TaskTemplate: v1alpha1.TaskTemplate{
				Pod: &podTemplateSpecUpdatedCommand,
			},
			TaskPendingTimeoutSeconds: pointer.Int64(1800),
			MaxAttempts:               pointer.Int64(5),
			RetryDelaySeconds:         pointer.Int64(60),
		},
	}

	jobTemplateSpecMoreRetries = v1alpha1.JobTemplateSpec{
		ObjectMeta: jobTemplateSpecBasic.ObjectMeta,
		Spec: v1alpha1.JobTemplate{
			TaskTemplate:      jobTemplateSpecBasic.Spec.TaskTemplate,
			MaxAttempts:       pointer.Int64(10),
			RetryDelaySeconds: jobTemplateSpecBasic.Spec.RetryDelaySeconds,
		},
	}

	jobTemplateSpecTooManyRetries = v1alpha1.JobTemplateSpec{
		ObjectMeta: jobTemplateSpecBasic.ObjectMeta,
		Spec: v1alpha1.JobTemplate{
			TaskTemplate:      jobTemplateSpecBasic.Spec.TaskTemplate,
			MaxAttempts:       pointer.Int64(100),
			RetryDelaySeconds: jobTemplateSpecBasic.Spec.RetryDelaySeconds,
		},
	}

	concurrencySpecBasic = v1alpha1.ConcurrencySpec{
		Policy: v1alpha1.ConcurrencyPolicyForbid,
	}

	cronScheduleBasic = v1alpha1.CronSchedule{
		Expression: "5 10 * * *",
		Timezone:   "Asia/Singapore",
	}

	cronScheduleWithHash = v1alpha1.CronSchedule{
		Expression: "H 10 * * *",
	}

	invalidCronExpression = "500 10 * * *"

	cronScheduleInvalidExpression = v1alpha1.CronSchedule{
		Expression: invalidCronExpression,
		Timezone:   "Asia/Singapore",
	}

	cronScheduleInvalidTimezone = v1alpha1.CronSchedule{
		Expression: "5 10 * * *",
		Timezone:   "Invalid/Time/Zone",
	}

	scheduleSpecBasic = v1alpha1.ScheduleSpec{
		Cron:     &cronScheduleBasic,
		Disabled: true,
	}

	scheduleSpecWithHash = v1alpha1.ScheduleSpec{
		Cron: &cronScheduleWithHash,
	}

	concurrencySpecMissingPolicy = v1alpha1.ConcurrencySpec{
		Policy: "",
	}

	concurrencySpecInvalidPolicy = v1alpha1.ConcurrencySpec{
		Policy: "invalid",
	}

	scheduleSpecInvalidCronSchedule = v1alpha1.ScheduleSpec{
		Cron: &cronScheduleInvalidExpression,
	}

	scheduleSpecInvalidTimezone = v1alpha1.ScheduleSpec{
		Cron: &cronScheduleInvalidTimezone,
	}

	optionSpecBasic = v1alpha1.OptionSpec{
		Options: []v1alpha1.Option{
			{
				Type:  v1alpha1.OptionTypeBool,
				Name:  "bool_option",
				Label: "My Option",
				Bool: &v1alpha1.BoolOptionConfig{
					Default: true,
					Format:  v1alpha1.BoolOptionFormatTrueFalse,
				},
			},
		},
	}

	optionSpecInvalidName = v1alpha1.OptionSpec{
		Options: []v1alpha1.Option{
			{
				Type: v1alpha1.OptionTypeBool,
				Name: "invalid name here!",
				Bool: &v1alpha1.BoolOptionConfig{
					Format: v1alpha1.BoolOptionFormatTrueFalse,
				},
			},
		},
	}

	startPolicyBasic = v1alpha1.StartPolicySpec{
		ConcurrencyPolicy: v1alpha1.ConcurrencyPolicyAllow,
	}

	startPolicyMissingConcurrencyPolicy = v1alpha1.StartPolicySpec{
		ConcurrencyPolicy: "",
	}

	startPolicyInvalidConcurrencyPolicy = v1alpha1.StartPolicySpec{
		ConcurrencyPolicy: "invalid",
	}
)

func TestValidateJobConfig(t *testing.T) {
	tests := []struct {
		name    string
		rjc     *v1alpha1.JobConfig
		cfgs    map[configv1alpha1.ConfigName]runtime.Object
		wantErr string
	}{
		{
			name: "basic job config",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
				},
			},
		},
		{
			name: "name too long",
			rjc: &v1alpha1.JobConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: "01234567890123456789012345678901234567890123456789",
				},
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
				},
			},
			wantErr: "metadata.name: Invalid value: \"01234567890123456789012345678901234567890123456789\": cannot be more than 49 characters",
		},
		{
			name: "extended job config",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
					Schedule:    &scheduleSpecBasic,
					Option:      &optionSpecBasic,
				},
			},
		},
		{
			name: "can validate hash cron schedule",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
					Schedule:    &scheduleSpecWithHash,
				},
			},
		},
		{
			name: "cannot validate hash cron schedule if hashNames is disabled",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
					Schedule:    &scheduleSpecWithHash,
				},
			},
			cfgs: map[configv1alpha1.ConfigName]runtime.Object{
				configv1alpha1.CronExecutionConfigName: &configv1alpha1.CronExecutionConfig{
					CronHashNames: pointer.Bool(false),
				},
			},
			wantErr: "spec.schedule.cron.expression: Invalid value: \"H 10 * * *\": cannot parse cron schedule",
		},
		{
			name: "missing concurrency.policy",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecMissingPolicy,
				},
			},
			wantErr: "spec.concurrency.policy: Required value",
		},
		{
			name: "invalid concurrency.policy",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecInvalidPolicy,
				},
			},
			wantErr: "spec.concurrency.policy: Unsupported value: \"invalid\"",
		},
		{
			name: "invalid concurrency.maxConcurrency",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template: jobTemplateSpecBasic,
					Concurrency: v1alpha1.ConcurrencySpec{
						Policy:         v1alpha1.ConcurrencyPolicyForbid,
						MaxConcurrency: pointer.Int64(0),
					},
				},
			},
			wantErr: "spec.concurrency.maxConcurrency: Invalid value: 0",
		},
		{
			name: "cannot use concurrency.maxConcurrency with Allow",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template: jobTemplateSpecBasic,
					Concurrency: v1alpha1.ConcurrencySpec{
						Policy:         v1alpha1.ConcurrencyPolicyAllow,
						MaxConcurrency: pointer.Int64(3),
					},
				},
			},
			wantErr: "spec.concurrency.maxConcurrency: Forbidden: cannot specify maxConcurrency with Allow",
		},
		{
			name: "schedule without any schedule types",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
					Schedule: &v1alpha1.ScheduleSpec{
						Disabled: true,
					},
				},
			},
			wantErr: "spec.schedule: Required value: at least one schedule type must be specified",
		},
		{
			name: "invalid schedule.cron.expression",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
					Schedule:    &scheduleSpecInvalidCronSchedule,
				},
			},
			wantErr: "spec.schedule.cron.expression: Invalid value: \"500 10 * * *\": cannot parse cron schedule",
		},
		{
			name: "invalid schedule.cron.expressions",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
					Schedule: &v1alpha1.ScheduleSpec{
						Cron: &v1alpha1.CronSchedule{
							Expressions: []string{
								invalidCronExpression,
							},
						},
					},
				},
			},
			wantErr: "spec.schedule.cron.expressions[0]: Invalid value: \"500 10 * * *\": cannot parse cron schedule",
		},
		{
			name: "no expression specified",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
					Schedule: &v1alpha1.ScheduleSpec{
						Cron: &v1alpha1.CronSchedule{},
					},
				},
			},
			wantErr: "spec.schedule.cron: Required value: expression is required",
		},
		{
			name: "too many expressions specified",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
					Schedule: &v1alpha1.ScheduleSpec{
						Cron: &v1alpha1.CronSchedule{
							Expression: "* * * * *",
							Expressions: []string{
								"* * * * *",
							},
						},
					},
				},
			},
			wantErr: "spec.schedule.cron: Too many: 2: must have at most 1 items",
		},
		{
			name: "invalid schedule.cron.timezone",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
					Schedule:    &scheduleSpecInvalidTimezone,
				},
			},
			wantErr: "spec.schedule.cron.timezone: Invalid value: \"Invalid/Time/Zone\": cannot parse timezone",
		},
		{
			name: "maxAttempts too large",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecTooManyRetries,
					Concurrency: concurrencySpecBasic,
					Schedule:    &scheduleSpecBasic,
				},
			},
			wantErr: "spec.template.spec.maxAttempts: Invalid value: 100: must be less than or equal to 50",
		},
		{
			name: "invalid options",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Template:    jobTemplateSpecBasic,
					Concurrency: concurrencySpecBasic,
					Option:      &optionSpecInvalidName,
				},
			},
			wantErr: `spec.option.options[0]: Invalid value: "invalid name here!": must match regex: ^[a-zA-Z_0-9.-]+$`,
		},
		{
			name: "require task template",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Concurrency: concurrencySpecBasic,
					Template: v1alpha1.JobTemplateSpec{
						Spec: v1alpha1.JobTemplate{
							TaskTemplate: v1alpha1.TaskTemplate{},
						},
					},
				},
			},
			wantErr: "spec.template.spec.taskTemplate: Required value",
		},
		{
			name: "invalid pod template",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Concurrency: concurrencySpecBasic,
					Template: v1alpha1.JobTemplateSpec{
						Spec: v1alpha1.JobTemplate{
							TaskTemplate: v1alpha1.TaskTemplate{
								Pod: &podTemplateSpecEmpty,
							},
						},
					},
				},
			},
			wantErr: "spec.template.spec.taskTemplate.pod.spec.containers: Required value",
		},
		{
			name: "cannot use restartPolicy Always",
			rjc: &v1alpha1.JobConfig{
				Spec: v1alpha1.JobConfigSpec{
					Concurrency: concurrencySpecBasic,
					Template: v1alpha1.JobTemplateSpec{
						Spec: v1alpha1.JobTemplate{
							TaskTemplate: v1alpha1.TaskTemplate{
								Pod: &v1alpha1.PodTemplateSpec{
									Spec: corev1.PodSpec{
										Containers:    []corev1.Container{{Name: "container", Image: "alpine"}},
										RestartPolicy: corev1.RestartPolicyAlways,
									},
								},
							},
						},
					},
				},
			},
			wantErr: "spec.template.spec.taskTemplate.pod.spec.restartPolicy: Invalid value: \"Always\": restartPolicy cannot be Always",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.rjc.DeepCopy()
			ctrlContext := mock.NewContext()
			ctrlContext.MockConfigs().SetConfigs(tt.cfgs)
			if err := ctrlContext.Start(context.Background()); err != nil {
				t.Fatal(err)
			}
			validator := validation.NewValidator(ctrlContext)
			err := validator.ValidateJobConfig(tt.rjc).ToAggregate()
			if checkError(err, tt.wantErr) {
				t.Errorf("ValidateJobConfig() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !cmp.Equal(tt.rjc, original) {
				t.Errorf("ValidateJobConfig() mutated input\ndiff = %v", cmp.Diff(original, tt.rjc))
			}
		})
	}
}

func TestValidateJob(t *testing.T) {
	tests := []struct {
		name    string
		rj      *v1alpha1.Job
		wantErr string
	}{
		{
			name: "basic job",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
		},
		{
			name: "name too long",
			rj: &v1alpha1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "0123456789012345678901234567890123456789012345678901234567890",
					Labels: objectMetaJob.Labels,
				},
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			wantErr: "metadata.name: Invalid value: \"0123456789012345678901234567890123456789012345678901234567890\": cannot be more than 60 characters",
		},
		{
			name: "missing Type",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     "",
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			wantErr: "spec.type: Required value",
		},
		{
			name: "invalid Type",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     "invalid",
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			wantErr: "spec.type: Unsupported value: \"invalid\"",
		},
		{
			name: "job with startPolicy",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:        v1alpha1.JobTypeAdhoc,
					Template:    &jobTemplateSpecBasic.Spec,
					StartPolicy: &startPolicyBasic,
				},
			},
		},
		{
			name: "missing schedule.concurrencyPolicy",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:        v1alpha1.JobTypeAdhoc,
					Template:    &jobTemplateSpecBasic.Spec,
					StartPolicy: &startPolicyMissingConcurrencyPolicy,
				},
			},
			wantErr: "spec.startPolicy.concurrencyPolicy: Required value",
		},
		{
			name: "invalid schedule.concurrencyPolicy",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:        v1alpha1.JobTypeAdhoc,
					Template:    &jobTemplateSpecBasic.Spec,
					StartPolicy: &startPolicyInvalidConcurrencyPolicy,
				},
			},
			wantErr: "spec.startPolicy.concurrencyPolicy: Unsupported value: \"invalid\"",
		},
		{
			name: "invalid ttlSecondsAfterFinished",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:                    v1alpha1.JobTypeAdhoc,
					Template:                &jobTemplateSpecBasic.Spec,
					TTLSecondsAfterFinished: pointer.Int64(-300),
				},
			},
			wantErr: "spec.ttlSecondsAfterFinished: Invalid value: -300",
		},
		{
			name: "maxAttempts too large",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecTooManyRetries.Spec,
				},
			},
			wantErr: "spec.template.maxAttempts: Invalid value: 100: must be less than or equal to 50",
		},
		{
			name: "require task template",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						TaskTemplate: v1alpha1.TaskTemplate{},
					},
				},
			},
			wantErr: "spec.template.taskTemplate: Required value",
		},
		{
			name: "invalid pod template",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &podTemplateSpecEmpty,
						},
					},
				},
			},
			wantErr: "spec.template.taskTemplate.pod.spec.containers: Required value",
		},
		{
			name: "cannot use restartPolicy Always",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &v1alpha1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers:    []corev1.Container{{Name: "container", Image: "alpine"}},
									RestartPolicy: corev1.RestartPolicyAlways,
								},
							},
						},
					},
				},
			},
			wantErr: "spec.template.taskTemplate.pod.spec.restartPolicy: Invalid value: \"Always\": restartPolicy cannot be Always",
		},
		{
			name: "required template.parallelism.completionStrategy",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						Parallelism: &v1alpha1.ParallelismSpec{
							WithCount:          pointer.Int64(3),
							CompletionStrategy: "",
						},
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &podTemplateSpecBasic,
						},
					},
				},
			},
			wantErr: "spec.template.parallelism.completionStrategy: Required value",
		},
		{
			name: "invalid template.parallelism.completionStrategy",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						Parallelism: &v1alpha1.ParallelismSpec{
							WithCount:          pointer.Int64(3),
							CompletionStrategy: "invalid",
						},
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &podTemplateSpecBasic,
						},
					},
				},
			},
			wantErr: "spec.template.parallelism.completionStrategy: Unsupported value: \"invalid\"",
		},
		{
			name: "required template.parallelism type",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						Parallelism: &v1alpha1.ParallelismSpec{
							CompletionStrategy: v1alpha1.AllSuccessful,
						},
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &podTemplateSpecBasic,
						},
					},
				},
			},
			wantErr: "spec.template.parallelism: Required value: must specify a parallelism type",
		},
		{
			name: "negative template.parallelism.withCount",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						Parallelism: &v1alpha1.ParallelismSpec{
							WithCount:          pointer.Int64(-1),
							CompletionStrategy: v1alpha1.AllSuccessful,
						},
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &podTemplateSpecBasic,
						},
					},
				},
			},
			wantErr: "spec.template.parallelism.withCount: Invalid value: -1: must be greater than 0",
		},
		{
			name: "zero template.parallelism.withCount",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						Parallelism: &v1alpha1.ParallelismSpec{
							WithCount:          pointer.Int64(0),
							CompletionStrategy: v1alpha1.AllSuccessful,
						},
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &podTemplateSpecBasic,
						},
					},
				},
			},
			wantErr: "spec.template.parallelism.withCount: Invalid value: 0: must be greater than 0",
		},
		{
			name: "empty string in template.parallelism.withKeys",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						Parallelism: &v1alpha1.ParallelismSpec{
							WithKeys:           []string{"key1", ""},
							CompletionStrategy: v1alpha1.AllSuccessful,
						},
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &podTemplateSpecBasic,
						},
					},
				},
			},
			wantErr: "spec.template.parallelism.withKeys[1]: Required value: key cannot be empty",
		},
		{
			name: "empty string in template.parallelism.withMatrix",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						Parallelism: &v1alpha1.ParallelismSpec{
							WithMatrix: map[string][]string{
								"": {"value1", "value2"},
							},
							CompletionStrategy: v1alpha1.AllSuccessful,
						},
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &podTemplateSpecBasic,
						},
					},
				},
			},
			wantErr: `spec.template.parallelism.withMatrix: Invalid value: "": withMatrix key must match regexp`,
		},
		{
			name: "invalid key in template.parallelism.withMatrix",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						Parallelism: &v1alpha1.ParallelismSpec{
							WithMatrix: map[string][]string{
								"invalid key": {"value1", "value2"},
							},
							CompletionStrategy: v1alpha1.AllSuccessful,
						},
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &podTemplateSpecBasic,
						},
					},
				},
			},
			wantErr: `spec.template.parallelism.withMatrix: Invalid value: "invalid key": withMatrix key must match regexp`,
		},
		{
			name: "invalid value in template.parallelism.withMatrix",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					Template: &v1alpha1.JobTemplate{
						Parallelism: &v1alpha1.ParallelismSpec{
							WithMatrix: map[string][]string{
								"key": {"value1", ""},
							},
							CompletionStrategy: v1alpha1.AllSuccessful,
						},
						TaskTemplate: v1alpha1.TaskTemplate{
							Pod: &podTemplateSpecBasic,
						},
					},
				},
			},
			wantErr: "spec.template.parallelism.withMatrix[key]: Required value: value cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.rj.DeepCopy()
			ctrlContext := mock.NewContext()
			if err := ctrlContext.Start(context.Background()); err != nil {
				t.Fatal(err)
			}
			validator := validation.NewValidator(ctrlContext)
			err := validator.ValidateJob(tt.rj).ToAggregate()
			if checkError(err, tt.wantErr) {
				t.Errorf("ValidateJob() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !cmp.Equal(tt.rj, original) {
				t.Errorf("ValidateJob() mutated input\ndiff = %v", cmp.Diff(original, tt.rj))
			}
		})
	}
}

func TestValidateJobUpdate(t *testing.T) {
	tests := []struct {
		name    string
		oldRj   *v1alpha1.Job
		newRj   *v1alpha1.Job
		now     *time.Time
		wantErr string
	}{
		{
			name: "can update StartPolicy if not started",
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:        v1alpha1.JobTypeAdhoc,
					Template:    &jobTemplateSpecBasic.Spec,
					StartPolicy: &startPolicyBasic,
				},
			},
		},
		{
			name: "cannot update StartPolicy if started",
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
				Status: v1alpha1.JobStatus{
					StartTime: startTime,
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:        v1alpha1.JobTypeAdhoc,
					Template:    &jobTemplateSpecBasic.Spec,
					StartPolicy: &startPolicyBasic,
				},
				Status: v1alpha1.JobStatus{
					StartTime: startTime,
				},
			},
			wantErr: "spec.startPolicy: Invalid value: v1alpha1.StartPolicySpec{ConcurrencyPolicy:\"Allow\", StartAfter:<nil>}: cannot update startPolicy once Job is started",
		},
		{
			name: "immutable label JobConfig UID",
			oldRj: &v1alpha1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name: objectMetaJob.Name,
					Labels: map[string]string{
						jobconfig.LabelKeyJobConfigUID: string(objectMetaJobConfig.UID),
					},
				},
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name: objectMetaJob.Name,
					Labels: map[string]string{
						jobconfig.LabelKeyJobConfigUID: "abc",
					},
				},
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			wantErr: "metadata.labels[execution.furiko.io/job-config-uid]: Invalid value: \"abc\": field is immutable",
		},
		{
			name: "immutable field ConfigName",
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					ConfigName: "jobconfig",
					Type:       v1alpha1.JobTypeAdhoc,
					Template:   &jobTemplateSpecBasic.Spec,
				},
			},
			wantErr: "spec.configName: Invalid value: \"jobconfig\": field is immutable",
		},
		{
			name: "immutable field type",
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeScheduled,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			wantErr: "spec.type: Invalid value: \"Scheduled\": field is immutable",
		},
		{
			name: "immutable field task",
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecContainerUpdatedCommand.Spec,
				},
			},
			wantErr: "spec.template.taskTemplate: Invalid value",
		},
		{
			name: "immutable field maxAttempts",
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecMoreRetries.Spec,
				},
			},
			wantErr: "spec.template.maxAttempts: Invalid value: 10: field is immutable",
		},
		{
			name: "can set KillTimestamp",
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:          v1alpha1.JobTypeAdhoc,
					Template:      &jobTemplateSpecBasic.Spec,
					KillTimestamp: testutils.Mkmtimep("2021-02-09T04:10:00Z"),
				},
			},
		},
		{
			name: "can set KillTimestamp in the past",
			now:  testutils.Mktimep("2021-02-09T04:12:00Z"),
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:          v1alpha1.JobTypeAdhoc,
					Template:      &jobTemplateSpecBasic.Spec,
					KillTimestamp: testutils.Mkmtimep("2021-02-09T04:10:00Z"),
				},
			},
		},
		{
			name: "can update KillTimestamp that was in the future",
			now:  testutils.Mktimep("2021-02-09T04:08:00Z"),
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:          v1alpha1.JobTypeAdhoc,
					Template:      &jobTemplateSpecBasic.Spec,
					KillTimestamp: testutils.Mkmtimep("2021-02-09T04:10:00Z"),
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:          v1alpha1.JobTypeAdhoc,
					Template:      &jobTemplateSpecBasic.Spec,
					KillTimestamp: testutils.Mkmtimep("2021-02-09T04:15:00Z"),
				},
			},
		},
		{
			name: "can shorten KillTimestamp to the past that was originally in the future",
			now:  testutils.Mktimep("2021-02-09T04:08:00Z"),
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:          v1alpha1.JobTypeAdhoc,
					Template:      &jobTemplateSpecBasic.Spec,
					KillTimestamp: testutils.Mkmtimep("2021-02-09T04:10:00Z"),
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:          v1alpha1.JobTypeAdhoc,
					Template:      &jobTemplateSpecBasic.Spec,
					KillTimestamp: testutils.Mkmtimep("2021-02-09T04:05:00Z"),
				},
			},
		},
		{
			name: "cannot update KillTimestamp that was originally in the past",
			now:  testutils.Mktimep("2021-02-09T04:12:00Z"),
			oldRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:          v1alpha1.JobTypeAdhoc,
					Template:      &jobTemplateSpecBasic.Spec,
					KillTimestamp: testutils.Mkmtimep("2021-02-09T04:10:00Z"),
				},
			},
			newRj: &v1alpha1.Job{
				ObjectMeta: objectMetaJob,
				Spec: v1alpha1.JobSpec{
					Type:          v1alpha1.JobTypeAdhoc,
					Template:      &jobTemplateSpecBasic.Spec,
					KillTimestamp: testutils.Mkmtimep("2021-02-09T04:15:00Z"),
				},
			},
			wantErr: "spec.killTimestamp: Invalid value: 2021-02-09 04:15:00 +0000 UTC: field is immutable once passed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClock := fakeclock.NewFakeClock(time.Now())
			validation.Clock = fakeClock
			if tt.now != nil {
				fakeClock.SetTime(*tt.now)
			}
			oldOriginal := tt.oldRj.DeepCopy()
			newOriginal := tt.newRj.DeepCopy()
			ctrlContext := mock.NewContext()
			if err := ctrlContext.Start(context.Background()); err != nil {
				t.Fatal(err)
			}
			validator := validation.NewValidator(ctrlContext)
			err := validator.ValidateJobUpdate(tt.oldRj, tt.newRj).ToAggregate()
			if checkError(err, tt.wantErr) {
				t.Errorf("ValidateJobUpdate() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !cmp.Equal(tt.oldRj, oldOriginal) {
				t.Errorf("ValidateJob() mutated input\ndiff = %v", cmp.Diff(oldOriginal, tt.oldRj))
			}
			if !cmp.Equal(tt.newRj, newOriginal) {
				t.Errorf("ValidateJob() mutated input\ndiff = %v", cmp.Diff(newOriginal, tt.newRj))
			}
		})
	}
}

func TestValidateJobCreate(t *testing.T) {
	tests := []struct {
		name    string
		cfgs    controllercontext.ConfigsMap
		rj      *v1alpha1.Job
		rjcs    []*v1alpha1.JobConfig
		wantErr string
	}{
		{
			name: "can create Job without JobConfig",
			rj: &v1alpha1.Job{
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
		},
		{
			name: "cannot look up JobConfig",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJobWithAllReferences,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					StartPolicy: &v1alpha1.StartPolicySpec{
						ConcurrencyPolicy: v1alpha1.ConcurrencyPolicyForbid,
					},
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			wantErr: `metadata.ownerReferences[0]: Not found: "jobconfig-sample"`,
		},
		{
			name: "missing UID label to ensure JobConfig reference",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJobWithNoLabelUID,
				Spec: v1alpha1.JobSpec{
					Type: v1alpha1.JobTypeAdhoc,
					StartPolicy: &v1alpha1.StartPolicySpec{
						ConcurrencyPolicy: v1alpha1.ConcurrencyPolicyForbid,
					},
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			rjcs: []*v1alpha1.JobConfig{
				{
					ObjectMeta: objectMetaJobConfig,
					Status: v1alpha1.JobConfigStatus{
						Active: 5,
					},
				},
			},
			wantErr: "metadata.labels[execution.furiko.io/job-config-uid]: Required value: label must be specified if ownerReferences is specified",
		},
		{
			name: "can create Job with no startPolicy regardless",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJobWithAllReferences,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
				},
			},
			rjcs: []*v1alpha1.JobConfig{
				{
					ObjectMeta: objectMetaJobConfig,
					Status: v1alpha1.JobConfigStatus{
						Active: 5,
					},
				},
			},
		},
		{
			name: "can create Job with startPolicy.concurrencyPolicy Allow regardless",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJobWithAllReferences,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
					StartPolicy: &v1alpha1.StartPolicySpec{
						ConcurrencyPolicy: v1alpha1.ConcurrencyPolicyAllow,
					},
				},
			},
			rjcs: []*v1alpha1.JobConfig{
				{
					ObjectMeta: objectMetaJobConfig,
					Status: v1alpha1.JobConfigStatus{
						Active: 5,
					},
				},
			},
		},
		{
			name: "can create Job with startPolicy.concurrencyPolicy Enqueue regardless",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJobWithAllReferences,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
					StartPolicy: &v1alpha1.StartPolicySpec{
						ConcurrencyPolicy: v1alpha1.ConcurrencyPolicyEnqueue,
					},
				},
			},
			rjcs: []*v1alpha1.JobConfig{
				{
					ObjectMeta: objectMetaJobConfig,
					Status: v1alpha1.JobConfigStatus{
						Active: 5,
					},
				},
			},
		},
		{
			name: "cannot create Job with concurrencyPolicy Forbid with Active",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJobWithAllReferences,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
					StartPolicy: &v1alpha1.StartPolicySpec{
						ConcurrencyPolicy: v1alpha1.ConcurrencyPolicyForbid,
					},
				},
			},
			rjcs: []*v1alpha1.JobConfig{
				{
					ObjectMeta: objectMetaJobConfig,
					Status: v1alpha1.JobConfigStatus{
						Active: 5,
					},
				},
			},
			wantErr: "spec.startPolicy.concurrencyPolicy: Forbidden: jobconfig-sample currently has 5 active job(s), but concurrency policy forbids exceeding maximum concurrency of 1",
		},
		{
			name: "can create Job with startPolicy.concurrencyPolicy Forbid with no Active",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJobWithAllReferences,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
					StartPolicy: &v1alpha1.StartPolicySpec{
						ConcurrencyPolicy: v1alpha1.ConcurrencyPolicyForbid,
					},
				},
			},
			rjcs: []*v1alpha1.JobConfig{
				{
					ObjectMeta: objectMetaJobConfig,
					Status: v1alpha1.JobConfigStatus{
						Active: 0,
					},
				},
			},
		},
		{
			name: "cannot create Job with too many Queued",
			cfgs: controllercontext.ConfigsMap{
				configv1alpha1.JobConfigExecutionConfigName: &configv1alpha1.JobConfigExecutionConfig{
					MaxEnqueuedJobs: pointer.Int64(5),
				},
			},
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJobWithAllReferences,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
					StartPolicy: &v1alpha1.StartPolicySpec{
						ConcurrencyPolicy: v1alpha1.ConcurrencyPolicyEnqueue,
					},
				},
			},
			rjcs: []*v1alpha1.JobConfig{
				{
					ObjectMeta: objectMetaJobConfig,
					Status: v1alpha1.JobConfigStatus{
						Queued: 5,
					},
				},
			},
			wantErr: "spec.startPolicy: Forbidden: cannot create new Job for JobConfig jobconfig-sample, which would exceed maximum queue length of 5",
		},
		{
			name: "can create Job with concurrencyPolicy Forbid with Active and maxConcurrency",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJobWithAllReferences,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
					StartPolicy: &v1alpha1.StartPolicySpec{
						ConcurrencyPolicy: v1alpha1.ConcurrencyPolicyForbid,
					},
				},
			},
			rjcs: []*v1alpha1.JobConfig{
				{
					ObjectMeta: objectMetaJobConfig,
					Spec: v1alpha1.JobConfigSpec{
						Concurrency: v1alpha1.ConcurrencySpec{
							Policy:         v1alpha1.ConcurrencyPolicyForbid,
							MaxConcurrency: pointer.Int64(3),
						},
					},
					Status: v1alpha1.JobConfigStatus{
						Active: 1,
					},
				},
			},
		},
		{
			name: "cannot create Job with concurrencyPolicy Forbid with Active and exceed maxConcurrency",
			rj: &v1alpha1.Job{
				ObjectMeta: objectMetaJobWithAllReferences,
				Spec: v1alpha1.JobSpec{
					Type:     v1alpha1.JobTypeAdhoc,
					Template: &jobTemplateSpecBasic.Spec,
					StartPolicy: &v1alpha1.StartPolicySpec{
						ConcurrencyPolicy: v1alpha1.ConcurrencyPolicyForbid,
					},
				},
			},
			rjcs: []*v1alpha1.JobConfig{
				{
					ObjectMeta: objectMetaJobConfig,
					Spec: v1alpha1.JobConfigSpec{
						Concurrency: v1alpha1.ConcurrencySpec{
							Policy:         v1alpha1.ConcurrencyPolicyForbid,
							MaxConcurrency: pointer.Int64(3),
						},
					},
					Status: v1alpha1.JobConfigStatus{
						Active: 3,
					},
				},
			},
			wantErr: "spec.startPolicy.concurrencyPolicy: Forbidden: jobconfig-sample currently has 3 active job(s), but concurrency policy forbids exceeding maximum concurrency of 3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalRj := tt.rj.DeepCopy()
			validator := setup(t, tt.cfgs, tt.rjcs)
			err := validator.ValidateJobCreate(tt.rj).ToAggregate()
			if checkError(err, tt.wantErr) {
				t.Errorf("ValidateJobCreate() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !cmp.Equal(tt.rj, originalRj) {
				t.Errorf("ValidateJob() mutated input\ndiff = %v", cmp.Diff(originalRj, tt.rj))
			}
		})
	}
}

func setup(t *testing.T, cfgs map[configv1alpha1.ConfigName]runtime.Object, rjcs []*v1alpha1.JobConfig) *validation.Validator {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrlContext := mock.NewContext()
	ctrlContext.MockConfigs().SetConfigs(cfgs)
	hasSynced := ctrlContext.Informers().Furiko().Execution().V1alpha1().JobConfigs().Informer().HasSynced
	if err := ctrlContext.Start(ctx); err != nil {
		t.Fatal(err)
	}

	for _, rjc := range rjcs {
		_, err := ctrlContext.Clientsets().Furiko().ExecutionV1alpha1().JobConfigs(rjc.Namespace).
			Create(ctx, rjc, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("cannot create JobConfig: %v", err)
		}
	}

	if !cache.WaitForCacheSync(ctx.Done(), hasSynced) {
		t.Fatalf("cannot sync caches")
	}

	return validation.NewValidator(ctrlContext)
}

func checkError(err error, wantErr string) bool {
	// Mismatched error
	if (err == nil) != (wantErr == "") {
		return true
	}

	// Wrong error
	if err != nil && !strings.HasPrefix(err.Error(), wantErr) {
		return true
	}

	return false
}
