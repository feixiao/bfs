package main

import (
	"github.com/feixiao/bfs/proxy/conf"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	log "github.com/golang/glog"
)

const (
	version = "1.0.0"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "./proxy.toml", " set directory config file path")
}

func main() {
	var (
		c   *conf.Config
		err error
	)
	flag.Parse()
	defer log.Flush()
	log.Infof("bfs proxy [version: %s] start", version)
	if c, err = conf.NewConfig(configFile); err != nil {
		log.Errorf("NewConfig(\"%s\") error(%v)", configFile, err)
		panic(err)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	// init http
	if err = StartAPI(c); err != nil {
		log.Errorf("http.Init() error(%v)", err)
		panic(err)
	}
	if c.PprofEnable {
		log.Infof("init http pprof...")
		StartPprof(c.PprofListen)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-ch
		log.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
