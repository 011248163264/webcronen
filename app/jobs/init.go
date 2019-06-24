package jobs

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/linhux/webcronen/app/models"
	"os/exec"
	"time"
)

func InitJobs() {
	list, _ := models.TaskGetList(1, 1000000, "status", 1)
	for _, task := range list {
		job, err := NewJobFromTask(task)
		if err != nil {
			beego.Error("InitJobs:", err.Error())
			continue
		}
		AddJob(task.CronSpec, job)
	}
}

func runCmdWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		beego.Warn(fmt.Sprintf("Task execution time exceeds %d seconds，The process will be forced to kill: %d", int(timeout/time.Second), cmd.Process.Pid))
		go func() {
			<-done // Read the above go routine data，Avoid blocking and can't quit
		}()
		if err = cmd.Process.Kill(); err != nil {
			beego.Error(fmt.Sprintf("The process can't kill: %d, Error message: %s", cmd.Process.Pid, err))
		}
		return err, true
	case err = <-done:
		return err, false
	}
}
