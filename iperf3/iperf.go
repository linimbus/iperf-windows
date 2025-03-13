package iperf3

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

type Connect struct {
	Socket     int64  `json:"socket"`
	LocalHost  string `json:"local_host"`
	LocalPort  int64  `json:"local_port"`
	RemoteHost string `json:"remote_host"`
	RemotePort int64  `json:"remote_port"`
}

type TimeStamp struct {
	Time     string `json:"time"`
	TimeSecs int64  `json:"timesecs"`
}

type Connecting struct {
	Host string `json:"host"`
	Port int64  `json:"port"`
}

type TestConfig struct {
	Protocol   string `json:"protocol"`
	NumStreams int64  `json:"num_streams"`
	BlkSize    int64  `json:"blksize"`
	Omit       int64  `json:"omit"`
	Duration   int64  `json:"duration"`
	Bytes      int64  `json:"bytes"`
	Blocks     int64  `json:"blocks"`
	Reverse    int64  `json:"reverse"`
}

type Start struct {
	Connected     []Connect  `json:"connected"`
	Version       string     `json:"version"`
	SystemInfo    string     `json:"system_info"`
	Timestamp     TimeStamp  `json:"timestamp"`
	ConnectingTo  Connecting `json:"connecting_to"`
	Cookie        string     `json:"cookie"`
	TcpMssDefault int64      `json:"tcp_mss_default"`
	TestStart     TestConfig `json:"test_start"`
}

type Stream struct {
	Socket       int64   `json:"socket"`
	Start        float64 `json:"start"`
	End          float64 `json:"end"`
	Seconds      float64 `json:"seconds"`
	Bytes        int64   `json:"bytes"`
	BitPerSecond float64 `json:"bits_per_second"`
	Omitted      bool    `json:"omitted"`
}

type Sum struct {
	Start        float64 `json:"start"`
	End          float64 `json:"end"`
	Seconds      float64 `json:"seconds"`
	Bytes        int64   `json:"bytes"`
	BitPerSecond float64 `json:"bits_per_second"`
	Omitted      bool    `json:"omitted"`
}

type Interval struct {
	Streams []Stream `json:"streams"`
	Sum     Sum      `json:"sum"`
}

type CpuUtilPercent struct {
	HostTotal    float64 `json:"host_total"`
	HostUser     float64 `json:"host_user"`
	HostSystem   float64 `json:"host_system"`
	RemoteTotal  float64 `json:"remote_total"`
	RemoteUser   float64 `json:"remote_user"`
	RemoteSystem float64 `json:"remote_system"`
}

type StreamResult struct {
	Sender   Stream `json:"sender"`
	Receiver Stream `json:"receiver"`
}

type End struct {
	Streams     []StreamResult `json:"streams"`
	SumSender   Sum            `json:"sum_sent"`
	SumReceiver Sum            `json:"sum_received"`
	CpuPercent  CpuUtilPercent `json:"cpu_utilization_percent"`
}

type Result struct {
	Start     Start      `json:"start"`
	End       End        `json:"end"`
	Intervals []Interval `json:"intervals"`
}

type IperfServer struct {
	running  bool
	exitCode int
	stdOut   io.ReadCloser
	stdErr   io.ReadCloser
	cancel   context.CancelFunc
}

func ExecuteAsync(binary string, cmd []string) (stdOut io.ReadCloser, stdErr io.ReadCloser, cancel context.CancelFunc, exitCode chan int, err error) {
	exitCode = make(chan int)

	logs.Info("ExecuteAsync %s %v", binary, cmd)

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
		logs.Info("iperf3.exe shutdown")
	}
}

func ReaderScan(prefix string, buff io.ReadCloser, outputDir string) {
	if buff == nil {
		logs.Info("unable to read, ReadCloser is nil")
		return
	}
	defer buff.Close()

	text := make([]byte, 0)

	scanner := bufio.NewScanner(buff)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		text = append(text, scanner.Bytes()...)
		if json.Valid(text) {
			jsonFmt, err := FormatJSON(text)
			if err != nil {
				logs.Warning("json format fail, %s", err.Error())
			} else {
				text = jsonFmt
			}

			logs.Info("%s -> %s", prefix, text)

			if outputDir != "" {
				SaveToFile(filepath.Join(outputDir, fmt.Sprintf("iperf3_%s.json", GetTimestamp())), text)
			}

			var result Result
			if err := json.Unmarshal(text, &result); err != nil {
				logs.Info("json unmarshal fail, %s", err.Error())
			}

			text = make([]byte, 0)
		}
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

	go ReaderScan("server_stdout: ", stdout, configCache.ServerLog)
	go ReaderScan("server_stderr: ", stdErr, configCache.ServerLog)

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
	fmt.Fprintf(&builder, " --interval 1")

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

	go ReaderScan("client stdout: ", stdout, configCache.ClientLog)
	go ReaderScan("client stderr: ", stdErr, configCache.ClientLog)

	go func() {
		exitCode := <-exitCode

		srv.exitCode = exitCode
		srv.running = false
	}()

	return srv, nil
}
