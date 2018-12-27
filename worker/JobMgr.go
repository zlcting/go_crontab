package worker

import (
	"context"
	"time"
	"zlc_sys/common"

	"go.etcd.io/etcd/mvcc/mvccpb"

	"go.etcd.io/etcd/clientv3"
)

// 任务管理器
type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	G_jobMgr *JobMgr
)

//监听任务变化
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {

	}

	for _, kvpair = range getResp.Kvs {

		if job, err = common.UnpackJob(kvpair.Value); err == nil {

			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			//TODO 把任务同步给调度协程
			//fmt.Println(jobEvent)
			G_scheduler.PushJobEnvent(jobEvent)

		}
	}
	//从该revision 向后监听变化事件
	go func() {
		//从get时刻的后续版本监听
		watchStartRevision = getResp.Header.Revision + 1
		//启动监听 cron/josbs 目录的后续变化
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {

				switch watchEvent.Type {
				case mvccpb.PUT: //任务保存

					//反序列化job
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE:
					//delete /cron/job/
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))
					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)

				}
				//todo 推给scheduler
				//fmt.Println(*jobEvent)
				G_scheduler.PushJobEnvent(jobEvent)

			}
		}

	}()
	return
}

//初始化
func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	//初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}
	if client, err = clientv3.New(config); err != nil {
		return
	}

	//得到kv和lease的api
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)
	//赋值单例
	G_jobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}
	G_jobMgr.watchJobs()
	return
}

//创建任务执行锁
func (jobMgr *JobMgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	//返回一把锁

	jobLock = InitJobLock(jobName, jobMgr.kv, jobMgr.lease)
	return
}
