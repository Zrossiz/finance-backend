package cronscheduler

import (
	"context"
	"fmt"

	"github.com/Zrossiz/finance-backend/internal/apperrors"
	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	cr    *cron.Cron
	Tasks []Task
}

type Task struct {
	CronTaskID cron.EntryID
	Name       string
	Schedule   string
	Handler    func()
}

func New(cr *cron.Cron) *CronScheduler {
	return &CronScheduler{
		cr:    cr,
		Tasks: make([]Task, 0),
	}
}

func (c *CronScheduler) StartOnInit(tasks []Task) {
	for _, v := range tasks {
		go v.Handler()
	}
}

func (c *CronScheduler) AddTask(task Task) error {
	if task.Name == "" {
		return fmt.Errorf("cron task name cannot be empty")
	}

	if task.Handler == nil {
		return fmt.Errorf("cron handler cannot be nil")
	}

	for _, v := range c.Tasks {
		if v.Name == task.Name {
			return fmt.Errorf("cron task with name: %s already exist", task.Name)
		}
	}

	id, err := c.cr.AddFunc(task.Schedule, task.Handler)
	if err != nil {
		return err
	}

	task.CronTaskID = id
	c.Tasks = append(c.Tasks, task)

	return nil
}

func (c *CronScheduler) RemoveTask(name string) error {
	for i := 0; i < len(c.Tasks); i++ {
		if c.Tasks[i].Name == name {
			c.cr.Remove(c.Tasks[i].CronTaskID)
			c.Tasks = append(c.Tasks[:i], c.Tasks[i+1:]...)
			return nil
		}
	}

	return apperrors.ErrNotFound
}

func (c *CronScheduler) Stop() context.Context {
	ctx := c.cr.Stop()

	return ctx
}
