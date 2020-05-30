package main

import (
	"com.phpstu/webapp/src/config"
	"com.phpstu/webapp/src/models"
	"com.phpstu/webapp/src/routers"
	"context"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	//版本
	version string
	//构建时间
	built string
	//golang版本
	goversion string
	//构建系统信息
	osversion string
	//构建的git提交信息
	gitcommit string
)

func main() {
	var envPrefix = flag.String("env_prefix", "WEBAPP", "请输入服务名称 如:webapp")
	var cfg = flag.String("config", "", "请输入配置文件名称 如:config")
	var pvs = flag.Bool("version", false, "打印版本信息")
	flag.Parse()

	//版本打印
	if *pvs {
		versionInfo := "Version: " + version + "\nBuilt: " + built + "\nGo Version: " + goversion + "\nBuiltOs: " + osversion + "\nGit commit: " + gitcommit
		println(versionInfo)
		os.Exit(0)
	}
	//初始化配置文件 日志
	if err := config.Init(*cfg, *envPrefix); err != nil {
		panic(err)
	}
	models.Db = &models.Databases{}
	//初始化数据库
	if err := models.Db.InitMySql(); err != nil {
		panic(err)
	}
	defer models.Db.CloseMysql()
	//初始化redis
	if err := models.Db.InitRedis(); err != nil {
		panic(err)
	}
	defer models.Db.CloseRedis()

	done := make(chan struct{}, 1)
	//接受退出信息
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	server := newServer()
	go gracefulStopServer(server, quit, done)

	//启动监测
	go func() {
		err := checkHealth()
		if err != nil {
			panic(err)
		}
	}()

	zap.L().Error(server.ListenAndServe().Error())
	<-done
	zap.L().Info("Server is closed")

}

func newServer() *http.Server {
	gin.SetMode(viper.GetString("runmode"))
	g := gin.New()
	//加载路由
	routers.Load(g)

	return &http.Server{
		Addr:           viper.GetString("addr"),
		Handler:        g,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func gracefulStopServer(server *http.Server, quit <-chan os.Signal, done chan<- struct{}) {
	<-quit
	zap.L().Info("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	err := server.Shutdown(ctx)
	if err != nil {
		zap.L().Error("Could not gracefully shutdown the server", zap.Error(err))
	}
	close(done)
}

func checkHealth() error {
	tryNum := viper.GetInt("health_count")

	for i := 0; i < tryNum; i++ {
		time.Sleep(3 * time.Second)
		request, err := http.NewRequest(http.MethodGet, viper.GetString("url")+"/health", nil)
		if err != nil {
			return err
		}
		client := http.DefaultClient
		response, rerr := client.Do(request)
		if rerr == nil && response.StatusCode == http.StatusOK {
			return nil
		}
		go func() {
			defer response.Body.Close()
		}()
		zap.L().Info("one Second retry")
	}

	return errors.New("server no start,No connect to router")
}
