package executor

import (
	"errors"
	"testing"
	"time"
)

type concreteTask struct {
	duration   time.Duration
	name       string
	shouldFail bool
}

func newConcreteTask(name string, duration time.Duration, shouldFail bool) concreteTask {
	return concreteTask{
		duration:   duration,
		name:       name,
		shouldFail: shouldFail,
	}
}

func (t concreteTask) Run() error {
	time.Sleep(t.duration * time.Millisecond)
	if t.shouldFail {
		return errors.New("this task fails")
	}
	return nil
}

func (t concreteTask) GetName() string {
	return t.name
}

func TestParallelExecutor_Execute(t *testing.T) {
	type fields struct {
		concurrencyLimit int64
		tasks            []Task
		continueOnError  bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Parallel executor that should continue on any task error",
			fields: fields{
				concurrencyLimit: 3,
				continueOnError:  true,
				tasks: []Task{
					newConcreteTask("Success task 1", 3, false),
					newConcreteTask("Success task 2", 4, false),
					newConcreteTask("Success task 3", 2, false),
					newConcreteTask("Success task 4", 1, false),
					newConcreteTask("Failure task 1", 2, true),
					newConcreteTask("Failure task 2", 2, true),
					newConcreteTask("Success task 5", 1, false),
				},
			},
			wantErr: true,
		},
		{
			name: "Parallel executor that should NOT continue on any task error",
			fields: fields{
				concurrencyLimit: 3,
				continueOnError:  false,
				tasks: []Task{
					newConcreteTask("Success task 1", 3, false),
					newConcreteTask("Failure task 1", 2, true),
					newConcreteTask("Success task 2", 1, false),
					newConcreteTask("Success task 3", 2, false),
					newConcreteTask("Success task 4", 3, false),
					newConcreteTask("Failure task 2", 3, true),
					newConcreteTask("Success task 5", 3, false),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Parallel{
				concurrencyLimit: tt.fields.concurrencyLimit,
				tasks:            tt.fields.tasks,
				continueOnError:  tt.fields.continueOnError,
			}
			if err := e.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
