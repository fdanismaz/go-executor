# go-executor
```
go get github.com/fdanismaz/go-executor/executor
```

## Parallel Executor
The parallel executor is capable of running multiple independent tasks in separate go routines
asynchronously. Independent here means that the execution of a task does not affect the execution
of any other.

```go
func NewParallel(concurrencyLimit int64, continueOnError bool) Executor {
    ...
}
```
The `NewParallel` function creates takes the two following parameters:
- `concurrencyLimit`: limits the number of tasks to be executed at the same time. If you have 100 
  tasks to be executed and if you want max. 3 tasks to be executed at the same time, in that
  case the 4th task will not be started unless one of the first three task is completed.
- `continueOnError`: give `false` to stop the executor when a task fails. Give `true` otherwise.


```go
exec := executor.NewParallel(3, true)
exec.Submit(task1)
exec.Submit(task2)
exec.Submit(task3)
exec.Submit(task4)
exec.Submit(task5)
...
err := exec.execute()
```

The execute function returns the error of the first failed task if the `continueOnError` is
set to `false`. It will be the error of the first failed task because the executor will stop
when a task fails. 

It returns the error of the last failed task if the `continueOnError` is set to `true`, because
in that case the executor won't stop when a task fails and each time a task is failed, it will
update the error variable of the executor.

It returns `nil` when all tasks are completed successfully.
