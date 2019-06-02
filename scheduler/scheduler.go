package scheduler

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	tasks     = map[string]Task{}
	tasksLock = sync.RWMutex{}
)

type Task struct {
	Name            string
	Fn              func()
	Period          time.Duration
	FirstCallOnInit bool
	timer           *time.Timer
}

func ScheduleTask(t Task) {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	tasks[t.Name] = t
	go func(t Task) {
		t.timer = time.NewTimer(t.Period)

		for range t.timer.C {
			t.Fn()
		}
	}(t)

	if t.FirstCallOnInit {
		t.Fn()
	}

	log.Infof("Task '%s' scheduled every %s", t.Name, t.Period.String())
}

func ListTasks() []string {
	tasksLock.RLock()
	defer tasksLock.RUnlock()

	tasksList := make([]string, len(tasks))
	i := 0
	for k := range tasks {
		tasksList[i] = k
		i++
	}

	return tasksList
}

func Stop() {
	tasksLock.RLock()
	defer tasksLock.RUnlock()

	for name, t := range tasks {
		log.Infof("Stopping task '%s'", name)
		t.timer.Stop()
	}
}
