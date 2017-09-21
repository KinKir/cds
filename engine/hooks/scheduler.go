package hooks

import (
	"context"
	"sort"
	"time"

	"github.com/ovh/cds/sdk/log"
)

// Entry point of the internal scheduler
func (s *Service) runScheduler(c context.Context) error {
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	go func() {
		if err := s.dequeueTaskExecutions(ctx); err != nil {
			log.Error("Hooks> runScheduler> dequeueLongRunningTasks> %v", err)
			cancel()
		}
	}()

	go func() {
		if err := s.retryTaskExecutionsRoutine(ctx); err != nil {
			log.Error("Hooks> runScheduler> retryTaskExecutionsRoutine> %v", err)
			cancel()
		}
	}()

	go func() {
		if err := s.deleteTaskExecutionsRoutine(ctx); err != nil {
			log.Error("Hooks> runScheduler> deleteTaskExecutionsRoutine> %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	return ctx.Err()
}

// Every x seconds, the scehduler try to relaunch all tasks which have never been processed, or in error
func (s *Service) retryTaskExecutionsRoutine(c context.Context) error {
	tick := time.NewTicker(time.Duration(s.Cfg.RetryDelay) * time.Second)
	for {
		select {
		case <-c.Done():
			tick.Stop()
			return c.Err()
		case <-tick.C:
			tasks, err := s.Dao.FindAllTasks()
			if err != nil {
				log.Error("Hooks> retryTaskExecutionsRoutine > Unable to find all tasks: %v", err)
				continue
			}
			for _, t := range tasks {
				execs, err := s.Dao.FindAllTaskExecutions(&t)
				if err != nil {
					log.Error("Hooks> retryTaskExecutionsRoutine > Unable to find all task executions (%s): %v", t.UUID, err)
					continue
				}
				for _, e := range execs {
					if e.ProcessingTimestamp == 0 && e.Timestamp <= time.Now().UnixNano() {
						s.Dao.EnqueueTaskExecution(&e)
						continue
					}
					if e.NbErrors < s.Cfg.RetryError && e.LastError != "" {
						s.Dao.EnqueueTaskExecution(&e)
						continue
					}
				}
			}
		}
	}
}

// Every 60 seconds, old executions of each task are deleted
func (s *Service) deleteTaskExecutionsRoutine(c context.Context) error {
	tick := time.NewTicker(time.Duration(60) * time.Second)
	for {
		select {
		case <-c.Done():
			tick.Stop()
			return c.Err()
		case <-tick.C:
			tasks, err := s.Dao.FindAllTasks()
			if err != nil {
				log.Error("Hooks> deleteTaskExecutionsRoutine > Unable to find all tasks: %v", err)
				continue
			}
			for _, t := range tasks {
				execs, err := s.Dao.FindAllTaskExecutions(&t)
				if err != nil {
					log.Error("Hooks> deleteTaskExecutionsRoutine > Unable to find all task executions (%s): %v", t.UUID, err)
					continue
				}
				sort.Slice(execs, func(i, j int) bool {
					return execs[i].Timestamp > execs[j].Timestamp
				})

				for i, e := range execs {
					if i >= s.Cfg.ExecutionHistory && e.ProcessingTimestamp != 0 {
						s.Dao.DeleteTaskExecution(&e)
					}
				}
			}
		}
	}
}

// Get from queue task execution
func (s *Service) dequeueTaskExecutions(c context.Context) error {
	for {
		if c.Err() != nil {
			return c.Err()
		}

		var taskKey string
		s.Cache.DequeueWithContext(c, schedulerQueueKey, &taskKey)

		log.Debug("Dequeuing a task execution: %v", taskKey)

		var t = TaskExecution{}
		if !s.Cache.Get(taskKey, &t) {
			continue
		}

		//Execute
		if err := s.doTask(c, &t); err != nil {
			log.Error("Hooks> doLongRunningTask failed: %v", err)
			t.LastError = err.Error()
			t.NbErrors++
		}
		//Save the execution
		t.ProcessingTimestamp = time.Now().UnixNano()
		s.Dao.SaveTaskExecution(&t)

		continue
	}
}
