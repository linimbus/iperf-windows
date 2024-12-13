package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/astaxie/beego/logs"
)

type IperfServer struct {
	running  bool
	exitCode int
	stdOut   io.ReadCloser
	stdErr   io.ReadCloser
	cancel   context.CancelFunc
}

func ExecuteAsync(binary string, cmd []string) (stdOut io.ReadCloser, stdErr io.ReadCloser, cancel context.CancelFunc, exitCode chan int, err error) {
	exitCode = make(chan int)

	ctx, cancel := context.WithCancel(context.Background())
	exe := exec.CommandContext(ctx, binary, cmd...)
	exe.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	stdOut, err = exe.StdoutPipe()
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}
	stdErr, err = exe.StderrPipe()
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}
	err = exe.Start()
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}
	go func() {
		if err := exe.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					exitCode <- status.ExitStatus()
				}
			}
		} else {
			exitCode <- 0
		}
	}()
	return stdOut, stdErr, cancel, exitCode, nil
}

func (s *IperfServer) Shutdown() {
	if s.running && s.cancel != nil {
		s.cancel()
		time.Sleep(100 * time.Millisecond)
	}
}

func ReaderScan(prefix string, buff io.ReadCloser) {
	if buff == nil {
		logs.Info("unable to read, ReadCloser is nil")
		return
	}

	scanner := bufio.NewScanner(buff)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		text := scanner.Text()
		logs.Info("%s -> %s", prefix, text)
	}
}

func ServerStartup() (*IperfServer, error) {
	value, err := json.Marshal(configCache)
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
	} else {
		logs.Info("iperf server options %s", string(value))
	}

	builder := strings.Builder{}
	builder.WriteString(" -s")
	fmt.Fprintf(&builder, " --port %d", configCache.ServerPort)

	if configCache.ServerInterval > 0 {
		fmt.Fprintf(&builder, " --interval %d", configCache.ServerInterval)
	}

	if configCache.ServerJsonFormat {
		fmt.Fprintf(&builder, " --json %d", configCache.ServerInterval)
	}

	stdout, stdErr, cancel, exitCode, err := ExecuteAsync(filepath.Join(ToolDirGet(), "iperf3.exe"), strings.Fields(builder.String()))
	if err != nil {
		logs.Warning("iperf server startup failed, %s", err.Error())
		return nil, err
	}

	srv := new(IperfServer)
	srv.stdOut = stdout
	srv.stdErr = stdErr
	srv.cancel = cancel
	srv.running = true

	go ReaderScan("IPerf3 server stdout: ", stdout)
	go ReaderScan("IPerf3 server stderr: ", stdErr)

	go func() {
		exitCode := <-exitCode

		srv.exitCode = exitCode
		srv.running = false
	}()

	return srv, nil
}

func ClientStartup() (*IperfServer, error) {

	value, err := json.Marshal(configCache)
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
	} else {
		logs.Info("iperf client options %s", string(value))
	}

	builder := strings.Builder{}
	fmt.Fprintf(&builder, "-c %s", configCache.ClientAddress)
	fmt.Fprintf(&builder, " -p %d", configCache.ClientPort)
	fmt.Fprintf(&builder, " -t %d", configCache.ClientRunTime)
	fmt.Fprintf(&builder, " -P %d", configCache.ClientStreams)

	if configCache.ClientOmitSec > 0 {
		fmt.Fprintf(&builder, " -O %d", configCache.ClientOmitSec)
	}

	if configCache.ClientBandwidth > 0 {
		fmt.Fprintf(&builder, " -b %d%c", configCache.ClientBandwidth, configCache.ClientBandwidthUnit[0])
	}

	if configCache.ClientProtocol == "udp" {
		builder.WriteString(" -u")
	}

	if configCache.ClientNoDelay {
		builder.WriteString(" -N")
	}

	if configCache.ClientZeroCopy {
		builder.WriteString(" -Z")
	}

	if configCache.ClientPayload > 0 {
		fmt.Fprintf(&builder, " -l %d", configCache.ClientPayload)
	}

	if configCache.ClientJsonFormat {
		builder.WriteString(" -J")
	}

	builder.WriteString(" --get-server-output")

	stdout, stdErr, cancel, exitCode, err := ExecuteAsync(filepath.Join(ToolDirGet(), "iperf3.exe"), strings.Fields(builder.String()))
	if err != nil {
		logs.Warning("iperf server startup failed, %s", err.Error())
		return nil, err
	}

	srv := new(IperfServer)
	srv.stdOut = stdout
	srv.stdErr = stdErr
	srv.cancel = cancel
	srv.running = true

	go ReaderScan("IPerf3 client stdout: ", stdout)
	go ReaderScan("IPerf3 client stderr: ", stdErr)

	go func() {
		exitCode := <-exitCode

		srv.exitCode = exitCode
		srv.running = false
	}()

	return srv, nil
}
