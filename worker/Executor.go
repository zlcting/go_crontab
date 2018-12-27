package worker

import (
	"context"
	"os/exec"
	"time"
	"zlc_sys/common"
)

type Executor struct {
}

var (
	G_executor *Executor
)

func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			err     error
			output  []byte
			result  *common.JobExecuteResult
			jobLock *JobLock
		)
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}

		//todo获取分布式锁
		//初始化一把锁
		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)
		//任务开始时间
		result.StartTime = time.Now()

		err = jobLock.TryLock() //上锁

		defer jobLock.Unlock() //释放锁

		if err != nil { //上锁失败
			result.Err = err
			result.EndTime = time.Now()
		} else {
			//上锁成功后 重置任务启动时间
			result.StartTime = time.Now()

			//执行shell命令
			cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
			//执行输出
			output, err = cmd.CombinedOutput()
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}
		//任务执行完成 把执行的结果返回scheduler scheduler会从exttingTable 中删除执行记录
		G_scheduler.PushJobResult(result)
	}()
}

func InitExecutor() (err error) {

	G_executor = &Executor{}
	return
}
