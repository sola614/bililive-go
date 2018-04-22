package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hr3lxphr6j/bililive-go/src/configs"
	"github.com/hr3lxphr6j/bililive-go/src/instance"
	"github.com/hr3lxphr6j/bililive-go/src/lib/events"
	"github.com/hr3lxphr6j/bililive-go/src/lib/utils"
	"github.com/hr3lxphr6j/bililive-go/src/listeners"
	"github.com/hr3lxphr6j/bililive-go/src/log"
	"github.com/hr3lxphr6j/bililive-go/src/recorders"
	"os"
)

const (
	AppName     = "BiliLive"
	AppVersion  = "0.01"
	CommandName = "bililive"
)

var (
	h bool   // 帮助
	v bool   // 版本信息
	c string // 配置文件
)

func parse() {
	flag.BoolVar(&h, "h", false, "show help info")
	flag.BoolVar(&v, "v", false, "show version")
	flag.StringVar(&c, "c", "./config.yml", "config file")
	flag.Parse()
}

func help() {
	version()
	fmt.Fprintf(os.Stderr, "Usage: %s [-hv] [-c filename]\n\nOptions:\n  -h:\tthis help\n  -v:\tshow version and exit\n  -c:\tset configuration file (default: ./config.yml)\n", CommandName)
}

func version() {
	fmt.Fprintf(os.Stderr, "%s Version: %s\n", AppName, AppVersion)
}

func main() {
	if !utils.IsFFmpegExist() {
		fmt.Fprintf(os.Stderr, "FFmpeg not found, Please Check.\n")
		os.Exit(3)
	}

	parse()
	if h {
		help()
		return
	}
	if v {
		version()
		return
	}
	inst := new(instance.Instance)
	ctx := context.WithValue(context.Background(), instance.InstanceKey, inst)

	config, err := configs.NewConfig(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
	inst.Config = config

	events.NewIEventDispatcher(ctx)
	log.NewLogger(ctx)
	listeners.NewIListenerManager(ctx)
	recorders.NewIRecorderManager(ctx)

	inst.ListenerManager.Start(ctx)
	inst.RecorderManager.Start(ctx)

	inst.WaitGroup.Wait()
}
