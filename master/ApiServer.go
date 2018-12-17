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
// post job = {"name":"job1","command":"echo hello","cornExpr:*********"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {

	var err error
	var postJob string
	var job common.Job
	//1,post 表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//2.取表单中的job字段
	postJob = req.PostForm.Get("job")
	//3,反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	//4.保存到etcd
ERR:
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

	mux.HandleFunc("job/save", handleJobSave)

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
