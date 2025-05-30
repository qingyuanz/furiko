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

package job_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	execution "github.com/furiko-io/furiko/apis/execution/v1alpha1"
	"github.com/furiko-io/furiko/pkg/execution/tasks"
	jobutil "github.com/furiko-io/furiko/pkg/execution/util/job"
)

func TestGenerateTaskRefs(t *testing.T) {
	tests := []struct {
		name     string
		existing []execution.TaskRef
		tasks    []tasks.Task
		want     []execution.TaskRef
	}{
		{
			name:     "no tasks",
			existing: nil,
			tasks:    nil,
			want:     nil,
		},
		{
			name:     "add new task",
			existing: nil,
			tasks: []tasks.Task{
				&stubTask{
					taskRef: execution.TaskRef{
						Name:              "task1",
						CreationTimestamp: createTime,
					},
				},
			},
			want: []execution.TaskRef{
				{
					Name:              "task1",
					CreationTimestamp: createTime,
				},
			},
		},
		{
			name: "update existing task",
			existing: []execution.TaskRef{
				{
					Name:              "task1",
					CreationTimestamp: createTime,
				},
			},
			tasks: []tasks.Task{
				&stubTask{
					taskRef: execution.TaskRef{
						Name:              "task1",
						CreationTimestamp: createTime,
						FinishTimestamp:   &finishTime,
						Status: execution.TaskStatus{
							State:  execution.TaskTerminated,
							Result: execution.TaskSucceeded,
						},
					},
				},
			},
			want: []execution.TaskRef{
				{
					Name:              "task1",
					CreationTimestamp: createTime,
					FinishTimestamp:   &finishTime,
					Status: execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskSucceeded,
					},
					DeletedStatus: &execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskSucceeded,
					},
				},
			},
		},
		{
			name: "add new without existing task",
			existing: []execution.TaskRef{
				{
					Name:              "task1",
					CreationTimestamp: createTime,
					FinishTimestamp:   &finishTime,
					Status: execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskSucceeded,
					},
					DeletedStatus: &execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskSucceeded,
					},
				},
			},
			tasks: []tasks.Task{
				&stubTask{
					taskRef: execution.TaskRef{
						Name:              "task2",
						CreationTimestamp: createTime2,
					},
				},
			},
			want: []execution.TaskRef{
				{
					Name:              "task1",
					CreationTimestamp: createTime,
					FinishTimestamp:   &finishTime,
					Status: execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskSucceeded,
					},
					DeletedStatus: &execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskSucceeded,
					},
				},
				{
					Name:              "task2",
					CreationTimestamp: createTime2,
				},
			},
		},
		{
			name: "existing task transitioned from finished back to running",
			existing: []execution.TaskRef{
				{
					Name:              "task1",
					CreationTimestamp: createTime,
					RunningTimestamp:  &startTime,
					FinishTimestamp:   &finishTime,
					Status: execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskFailed,
					},
					DeletedStatus: &execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskFailed,
					},
				},
			},
			tasks: []tasks.Task{
				&stubTask{
					taskRef: execution.TaskRef{
						Name:              "task1",
						CreationTimestamp: createTime,
						RunningTimestamp:  &startTime,
						Status: execution.TaskStatus{
							State: execution.TaskStarting,
						},
					},
				},
			},
			want: []execution.TaskRef{
				{
					Name:              "task1",
					CreationTimestamp: createTime,
					RunningTimestamp:  &startTime,
					FinishTimestamp:   &finishTime,
					Status: execution.TaskStatus{
						State: execution.TaskStarting,
					},
					DeletedStatus: &execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskFailed,
					},
				},
			},
		},
		{
			name:     "multiple TaskRefs with same CreationTimestamp",
			existing: []execution.TaskRef{},
			tasks: []tasks.Task{
				&stubTask{
					taskRef: execution.TaskRef{
						Name:              "task2",
						CreationTimestamp: createTime,
						RunningTimestamp:  &startTime,
						Status: execution.TaskStatus{
							State: execution.TaskStarting,
						},
					},
				},
				&stubTask{
					taskRef: execution.TaskRef{
						Name:              "task1",
						CreationTimestamp: createTime,
						RunningTimestamp:  &startTime,
						Status: execution.TaskStatus{
							State: execution.TaskStarting,
						},
					},
				},
			},
			want: []execution.TaskRef{
				{
					Name:              "task1",
					CreationTimestamp: createTime,
					RunningTimestamp:  &startTime,
					Status: execution.TaskStatus{
						State: execution.TaskStarting,
					},
				},
				{
					Name:              "task2",
					CreationTimestamp: createTime,
					RunningTimestamp:  &startTime,
					Status: execution.TaskStatus{
						State: execution.TaskStarting,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		cmpOpts := cmp.Options{
			cmpopts.EquateEmpty(),
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := jobutil.GenerateTaskRefs(tt.existing, tt.tasks); !cmp.Equal(got, tt.want, cmpOpts) {
				t.Errorf("GenerateTaskRefs() = not equal\ndiff: %v", cmp.Diff(tt.want, got, cmpOpts))
			}
		})
	}
}

func TestUpdateTaskRefDeletedStatusIfNotSet(t *testing.T) {
	createStatus := func(deletedStatus *execution.TaskStatus) execution.JobStatus {
		return execution.JobStatus{
			CreatedTasks: 1,
			Tasks: []execution.TaskRef{
				{
					Name:              "task1",
					CreationTimestamp: createTime,
					Status: execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskSucceeded,
					},
					DeletedStatus: deletedStatus,
				},
			},
		}
	}

	type args struct {
		rj       *execution.Job
		taskName string
		status   execution.TaskStatus
	}
	tests := []struct {
		name string
		args args
		want *execution.Job
	}{
		{
			name: "Task ref does not exist",
			args: args{
				rj: &execution.Job{
					Status: createStatus(&execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskSucceeded,
					}),
				},
				taskName: "task2",
				status: execution.TaskStatus{
					State:  execution.TaskTerminated,
					Result: execution.TaskKilled,
				},
			},
			want: &execution.Job{
				Status: createStatus(&execution.TaskStatus{
					State:  execution.TaskTerminated,
					Result: execution.TaskSucceeded,
				}),
			},
		},
		{
			name: "Don't update non-nil DeletedStatus",
			args: args{
				rj: &execution.Job{
					Status: createStatus(&execution.TaskStatus{
						State:  execution.TaskTerminated,
						Result: execution.TaskSucceeded,
					}),
				},
				taskName: "task1",
				status: execution.TaskStatus{
					State:  execution.TaskTerminated,
					Result: execution.TaskKilled,
				},
			},
			want: &execution.Job{
				Status: createStatus(&execution.TaskStatus{
					State:  execution.TaskTerminated,
					Result: execution.TaskSucceeded,
				}),
			},
		},
		{
			name: "Update nil DeletedStatus",
			args: args{
				rj: &execution.Job{
					Status: createStatus(nil),
				},
				taskName: "task1",
				status: execution.TaskStatus{
					State:  execution.TaskTerminated,
					Result: execution.TaskKilled,
				},
			},
			want: &execution.Job{
				Status: createStatus(&execution.TaskStatus{
					State:  execution.TaskTerminated,
					Result: execution.TaskKilled,
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.args.rj.DeepCopy()

			got := jobutil.UpdateTaskRefDeletedStatusIfNotSet(tt.args.rj, tt.args.taskName, tt.args.status)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("UpdateTaskRefDeletedStatus() not equal:\ndiff: %v", cmp.Diff(tt.want, got))
			}

			if !cmp.Equal(original, tt.args.rj) {
				t.Errorf("Job was mutated:\ndiff: %v", cmp.Diff(original, tt.args.rj))
			}
		})
	}
}

func TestUpdateJobTaskRefs(t *testing.T) {
	tests := []struct {
		name  string
		rj    *execution.Job
		tasks []tasks.Task
		want  *execution.Job
	}{
		{
			name: "no tasks",
			rj:   &execution.Job{},
			want: &execution.Job{
				Status: execution.JobStatus{
					Tasks: []execution.TaskRef{},
				},
			},
		},
		{
			name: "single created task",
			rj:   &execution.Job{},
			tasks: []tasks.Task{
				&stubTask{
					taskRef: execution.TaskRef{
						Name:              "task1",
						CreationTimestamp: createTime,
						RetryIndex:        0,
					},
				},
			},
			want: &execution.Job{
				Status: execution.JobStatus{
					Tasks: []execution.TaskRef{
						{
							Name:              "task1",
							CreationTimestamp: createTime,
							RetryIndex:        0,
							Status:            execution.TaskStatus{},
						},
					},
					CreatedTasks: 1,
				},
			},
		},
		{
			name: "single running task",
			rj:   &execution.Job{},
			tasks: []tasks.Task{
				&stubTask{
					taskRef: execution.TaskRef{
						Name:              "task1",
						CreationTimestamp: createTime,
						RunningTimestamp:  &startTime,
						RetryIndex:        0,
						Status: execution.TaskStatus{
							State: execution.TaskRunning,
						},
					},
				},
			},
			want: &execution.Job{
				Status: execution.JobStatus{
					Tasks: []execution.TaskRef{
						{
							Name:              "task1",
							CreationTimestamp: createTime,
							RunningTimestamp:  &startTime,
							RetryIndex:        0,
							Status: execution.TaskStatus{
								State: execution.TaskRunning,
							},
						},
					},
					CreatedTasks: 1,
					RunningTasks: 1,
				},
			},
		},
		{
			name: "single killing task",
			rj:   &execution.Job{},
			tasks: []tasks.Task{
				&stubTask{
					taskRef: execution.TaskRef{
						Name:              "task1",
						CreationTimestamp: createTime,
						RunningTimestamp:  &startTime,
						RetryIndex:        0,
						Status: execution.TaskStatus{
							State: execution.TaskKilling,
						},
					},
				},
			},
			want: &execution.Job{
				Status: execution.JobStatus{
					Tasks: []execution.TaskRef{
						{
							Name:              "task1",
							CreationTimestamp: createTime,
							RunningTimestamp:  &startTime,
							RetryIndex:        0,
							Status: execution.TaskStatus{
								State: execution.TaskKilling,
							},
						},
					},
					CreatedTasks: 1,
					RunningTasks: 1,
				},
			},
		},
		{
			name: "single finished task",
			rj:   &execution.Job{},
			tasks: []tasks.Task{
				&stubTask{
					taskRef: execution.TaskRef{
						Name:              "task1",
						CreationTimestamp: createTime,
						RunningTimestamp:  &startTime,
						FinishTimestamp:   &finishTime,
						RetryIndex:        0,
						Status: execution.TaskStatus{
							State:  execution.TaskTerminated,
							Result: execution.TaskKilled,
						},
					},
				},
			},
			want: &execution.Job{
				Status: execution.JobStatus{
					Tasks: []execution.TaskRef{
						{
							Name:              "task1",
							CreationTimestamp: createTime,
							RunningTimestamp:  &startTime,
							FinishTimestamp:   &finishTime,
							RetryIndex:        0,
							Status: execution.TaskStatus{
								State:  execution.TaskTerminated,
								Result: execution.TaskKilled,
							},
							DeletedStatus: &execution.TaskStatus{
								State:  execution.TaskTerminated,
								Result: execution.TaskKilled,
							},
						},
					},
					CreatedTasks: 1,
					RunningTasks: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jobutil.UpdateJobTaskRefs(tt.rj, tt.tasks); !cmp.Equal(tt.want, got) {
				t.Errorf("UpdateJobTaskRefs() not equal\ndiff = %v", cmp.Diff(tt.want, got))
			}
		})
	}
}
