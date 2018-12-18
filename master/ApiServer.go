package master

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"
	"zlc_sys/common"
)

type ApiServer struct {
	httpServer *http.Server
}

//单例
var (
	G_apiServer *ApiServer
)

//保存任务接口 保存到etcd
// post job = {"name":"job1","command":"echo hello","cornExpr:*/5******"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {

	var err error
	var postJob string
	var job common.Job
	var oldjob *common.Job
	var bytes []byte
	//1,post 表单
	// if err = req.ParseForm(); err != nil {

	// 	goto ERR
	// }

	//2.取表单中的job字段
	//postJob = req.PostForm.Get("job")

	postJob = req.PostFormValue("job")

	//3,反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	//fmt.Println(job)
	//4.保存到etcd
	if oldjob, err = G_jobMgr.SaveJob(&job); err != nil {

		goto ERR
	}

	//5 返回应答
	if bytes, err = common.BuildResponse(0, "success", oldjob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:

	//异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

//任务删除
// post /job/delete name = job1

func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	var err error
	var name string
	var oldjob *common.Job
	var bytes []byte
	// if err = req.ParseForm(); err != nil {
	// 	goto ERR
	// }
	// name = req.PostForm.Get("name")
	name = req.PostFormValue("name")

	if oldjob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "sussces", oldjob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	//异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func handleJobList(resp http.ResponseWriter, req *http.Request) {
	var jobList []*common.Job
	var err error
	var bytes []byte

	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "sussces", jobList); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	//异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)
	// if err = req.ParseForm(); err != nil {
	// 	goto ERR
	// }

	//name = req.PostForm.Get("name")
	name = req.PostFormValue("name")
	// fmt.Println(name)
	G_jobMgr.KillJob(name)

	if bytes, err = common.BuildResponse(0, "sussces", nil); err == nil {
		resp.Write(bytes)
	}
	return
	// ERR:
	// 	//异常应答
	// 	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
	// 		resp.Write(bytes)
	// 	}
}

//初始化服务
func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)
	//配置路由
	mux = http.NewServeMux()

	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	//启动服务端
	go httpServer.Serve(listener)

	return
}
