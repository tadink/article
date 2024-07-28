package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"

	"GinTest/config"
	"GinTest/db"
	"GinTest/frontend"
	"GinTest/global"
	"GinTest/log"
)

var sites *sync.Map

func init() {
	log.InitLog()
	err := config.Init()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	err = db.Init()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	sites, err = db.LoadSites()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
func main() {
	if len(os.Args) < 2 {
		serverStart()
		return
	}
	switch os.Args[1] {
	case "start":
		handleStartCmd()
	case "stop":
		handleStopCmd()
	case "restart":
		handleStopCmd()
		handleStartCmd()
	default:
		fmt.Println("unknow command")
	}

}
func handleStartCmd() {
	cmd := exec.Command(os.Args[0])
	err := cmd.Start()
	if err != nil {
		fmt.Println("start error:", err.Error())
		return
	}
	pid := fmt.Sprintf("%d", cmd.Process.Pid)
	err = os.WriteFile("pid", []byte(pid), os.ModePerm)
	if err != nil {
		fmt.Println("写入pid文件错误", err.Error())
		cmd.Process.Kill()
		return
	}
	fmt.Println("启动成功", pid)
}
func handleStopCmd() {
	data, err := os.ReadFile("pid")
	if err != nil {
		fmt.Println("read pid error", err.Error())
		return
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		fmt.Println("read pid error", err.Error())
		return
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("find process error", err.Error())
		return
	}
	if runtime.GOOS == "windows" {
		err = process.Signal(syscall.SIGKILL)
	} else {
		err = process.Signal(syscall.SIGTERM)
	}
	if err != nil {
		fmt.Println("process.Signal error", err.Error())
		return
	}
	fmt.Println("程序已经停止")
}
func serverStart() {
	frontend.Start(sites)
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()
	<-ctx.Done()
	slog.Info("Shutdown Server ...")
	if err := frontend.Shutdown(); err != nil {
		slog.Error("Server Shutdown:", "error", err.Error())
	}
	for _, cleanup := range global.Cleanups {
		err := cleanup()
		if err != nil {
			slog.Error(err.Error())
		}
	}
	slog.Info("Server exiting")
}
