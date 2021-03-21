package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tianzuoan/go-blog/global"
	"github.com/tianzuoan/go-blog/internal/model"
	routers "github.com/tianzuoan/go-blog/internal/routes"
	"github.com/tianzuoan/go-blog/pkg/logger"
	"github.com/tianzuoan/go-blog/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"time"
)

// @title 枫叶博客系统
// @version 1.0
// @description 枫叶Go语言学习
// @termsOfService https://github.com/tianzuoan
func main() {
	//a := []int{2: 1, 8: 6}
	//fmt.Println(a)
	//return
	gin.SetMode(global.ServerSetting.RunMode)
	r := routers.NewRouter()
	s := http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        r,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	global.Logger.Infof("%s cheduzi   kkk", "测试看下效果！")
	_ = s.ListenAndServe()
}

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("global init setting failed! err:", err)
	}
	err = setupDb()
	if err != nil {
		log.Fatalf("global init db failed! err:", err)
	}

	setupLogger()
}

func setupSetting() error {
	setting2, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting2.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	err = setting2.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting2.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupDb() error {
	db, err := model.NewDBEngine(*global.DatabaseSetting)
	if err != nil {
		return err
	}
	global.DBEngine = db
	return nil
}

func setupLogger() {
	//初始化日志组件
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "【鱼水之恋】", log.LstdFlags).WithCallersFrame()
}
