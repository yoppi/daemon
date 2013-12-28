package daemon

import (
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"syscall"
)

type Daemon struct {
	Pid int
	Logfile string
}

func NewDaemon(logfile string) *Daemon {
	return &Daemon{Logfile: logfile}
}

func (d *Daemon) lock() {
	pidfile := path.Base(os.Args[0]) + ".pid"
	f, err := os.Create(pidfile)
	if err != nil {
		log.Fatalf("Cannot create %v\n", pidfile)
	}
	defer f.Close()
	f.WriteString(strconv.Itoa(d.Pid))
}

func (d *Daemon) redirectIO() {
	devnull, _ := os.OpenFile("/dev/null", os.O_RDWR, 0)
	devnullFd := int(devnull.Fd())

	syscall.Dup2(devnullFd, int(os.Stdin.Fd()))

	if len(d.Logfile) > 0 {
		logfile, err := os.Create(d.Logfile)
		if err != nil {
			log.Printf("Cannot create %v\n", d.Logfile)
			return
		}
		logFd := int(logfile.Fd())
		syscall.Dup2(logFd, int(os.Stdout.Fd()))
		syscall.Dup2(logFd, int(os.Stderr.Fd()))
		defer logfile.Close()
	} else {
		syscall.Dup2(devnullFd, int(os.Stdout.Fd()))
		syscall.Dup2(devnullFd, int(os.Stderr.Fd()))
	}
}

func (d *Daemon) setSID() int {
	ret, err := syscall.Setsid()
	if err != nil {
		log.Fatalf("syscall.Setsid() = %v\n", err)
	}
	if ret < 0 {
		return -1
	}
	return ret
}

func (d *Daemon) safeFork() {
	ret, ret2, err := syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		log.Fatalf("syscall.RawSyscall() = %v\n", err)
	}
	if ret2 < 0 {
		log.Fatalf("syscall.RawSyscall() = %v\n", ret2)
	}

	d.Pid = int(ret)

	if runtime.GOOS == "darwin" && ret2 == 1 {
		ret = 0
	}

	if ret > 0 {
		os.Exit(0)
	}

	d.lock()
}

func Daemonize(logfile string) *Daemon {
	if syscall.Getppid() == 1 {
		return nil
	}

	daemon := NewDaemon(logfile)
	daemon.safeFork()
	daemon.setSID()
	daemon.redirectIO()

	syscall.Umask(0)

	return daemon
}

