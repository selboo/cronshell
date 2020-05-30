package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/gcfg.v1"
)

var (
	cfg Config

	conf   = "/etc/cronshell.conf"
	shell  = "/bin/sh"
	expire = 50
	conn   redis.Conn
	err    error
	logstr string
)

// Config ...
type Config struct {
	Log struct {
		Logfile    string
		Outputfile string
	}
	Redis struct {
		Host string
		Port int64
	}
}

// Init ...
func Init() error {

	if len(os.Args) <= 2 {
		e := fmt.Errorf("args fail ")
		return e
	}

	if err = gcfg.ReadFileInto(&cfg, conf); err != nil {
		e := fmt.Errorf("config fail %v", err)
		return e
	}

	conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port))
	if err != nil {
		e := fmt.Errorf("redis conn fail %v", err)
		return e
	}

	return nil
}

func command(cmdStr string) (string, error) {
	cmd := exec.Command(shell, "-c", cmdStr)
	opBytes, err := cmd.Output()
	if err != nil {
		e := fmt.Errorf("command fail %v", err)
		return "", e
	}
	return string(opBytes), nil
}

func xmd5(cmdStr string) string {
	data := []byte(cmdStr)
	h := md5.New()
	h.Write(data)
	output := h.Sum(nil)
	return fmt.Sprintf("%x", output)
}

func logs(s string) {
	logFile, _ := os.OpenFile(cfg.Log.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer logFile.Close()

	debugLog := log.New(logFile, "[CRONSHELL] ", log.LstdFlags)
	debugLog.Println(s)
}

func main() {

	if err := Init(); err != nil {
		log.Fatalf("init fail %v", err)
		return
	}
	args := os.Args

	cmdstr := strings.Join(args[2:], " ")

	// 执行的命令MD5
	cmdmd5 := xmd5(cmdstr)

	ok, err := conn.Do("set", cmdmd5, cmdstr, "EX", expire, "NX")
	if err != nil {
		log.Fatalf("redis fail %v", err)
		return
	}

	if ok == "OK" {
		cmdres, err := command(cmdstr)
		if err != nil {
			logstr = fmt.Sprintf("EXEC: %v MD5: %v RES:\n%v", cmdstr, cmdmd5, err)
		} else {
			logstr = fmt.Sprintf("EXEC: %v MD5: %v RES:\n%v", cmdstr, cmdmd5, cmdres)
		}
	} else {
		logstr = fmt.Sprintf("SKIP: %v MD5: %v", cmdstr, cmdmd5)
	}

	logs(logstr)

}
