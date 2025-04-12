package iperf3

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
	stdOut   string
	stdErr   string
	cancel   context.CancelFunc
}

func ExecuteAsync(binary string, cmd []string) (*os.File, *os.File, context.CancelFunc, chan int, error) {
	logs.Info("ExecuteAsync %s %v", binary, cmd)

	stdoutTmp, err := os.CreateTemp("", "iperf3_win_stdout_*.json")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	stderrTmp, err := os.CreateTemp("", "iperf3_win_stderr_*.json")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	exe := exec.CommandContext(ctx, binary, cmd...)
	exe.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	exe.Stdout = stdoutTmp
	exe.Stderr = stderrTmp

	err = exe.Start()
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}

	exitCodeChan := make(chan int, 1)

	go func() {
		defer stdoutTmp.Close()
		defer stderrTmp.Close()

		exe.Wait()

		logs.Info("iperf3.exe exit code %d", exe.ProcessState.ExitCode())

		exitCodeChan <- exe.ProcessState.ExitCode()
	}()

	return stdoutTmp, stderrTmp, cancel, exitCodeChan, nil
}

func (s *IperfServer) Shutdown() {
	if s.running && s.cancel != nil {
		s.cancel()
		time.Sleep(100 * time.Millisecond)
		logs.Info("shutdown iperf3.exe")
	}
}

func ReadResult(filePath string, outputDir string) {
	text, err := os.ReadFile(filePath)
	if err != nil {
		logs.Error("read file %s failed, %s", filePath, err.Error())
		return
	}

	if !json.Valid(text) {
		logs.Error("json invalid, %s", string(text))
		return
	}

	text, err = FormatJSON(text)
	if err != nil {
		logs.Warning("json format fail, %s", err.Error())
		return
	}

	logs.Info("iperf3 result: %s", string(text))

	if outputDir != "" {
		SaveToFile(filepath.Join(outputDir, fmt.Sprintf("iperf3_%s.json", GetTimestamp())), text)
	}

	var result Result
	if err := json.Unmarshal(text, &result); err != nil {
		logs.Info("json unmarshal fail, %s", err.Error())
	}
	// TBD
}

func ServerStartup(index int) (*IperfServer, error) {
	value, err := json.Marshal(configCache)
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
	} else {
		logs.Info("iperf server options %s", string(value))
	}

	builder := strings.Builder{}
	builder.WriteString(" -s")

	fmt.Fprintf(&builder, " -B %s", configCache.ServerListen)
	fmt.Fprintf(&builder, " --port %d", configCache.ServerPort+index)

	if configCache.ServerInterval > 0 {
		fmt.Fprintf(&builder, " --interval %d", configCache.ServerInterval)
	}

	if configCache.ServerJsonFormat {
		fmt.Fprintf(&builder, " --json %d", configCache.ServerInterval)
	}

	fmt.Fprintf(&builder, " --forceflush")

	stdout, stdErr, cancel, exitCodeChan, err := ExecuteAsync(filepath.Join(ToolDirGet(), "iperf3.exe"), strings.Fields(builder.String()))
	if err != nil {
		logs.Warning("iperf server startup failed, %s", err.Error())
		return nil, err
	}

	srv := new(IperfServer)
	srv.stdOut = stdout.Name()
	srv.stdErr = stdErr.Name()
	srv.cancel = cancel
	srv.running = true

	go func() {
		exitCode := <-exitCodeChan

		logs.Info("iperf3.exe index: %d exit code %d", index, exitCode)

		ReadResult(stdout.Name(), configCache.ServerLog)
		ReadResult(stdErr.Name(), configCache.ServerLog)

		srv.exitCode = exitCode
		srv.running = false
	}()

	return srv, nil
}

func ClientStartup(cnt int) (*IperfServer, error) {

	value, err := json.Marshal(configCache)
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
	} else {
		logs.Info("iperf client run times %d with options %s", cnt, string(value))
	}

	builder := strings.Builder{}

	fmt.Fprintf(&builder, " -B %s", configCache.ClientListen)
	fmt.Fprintf(&builder, " -c %s", configCache.ClientAddress)
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

	if configCache.ClientWindows > 0 {
		fmt.Fprintf(&builder, " -w %d%c", configCache.ClientWindows, configCache.ClientWindowsUnit[0])
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

	if configCache.ClientReverseMode {
		builder.WriteString(" -R")
	}

	if configCache.ClientBidirectionalMode {
		builder.WriteString(" --bidir")
	}

	if configCache.ClientPayload > 0 {
		fmt.Fprintf(&builder, " -l %d", configCache.ClientPayload)
	}

	if configCache.ClientJsonFormat {
		builder.WriteString(" -J")
	}

	if configCache.ClientDscpValue > 0 {
		fmt.Fprintf(&builder, " --dscp %d", configCache.ClientDscpValue)
	}

	if configCache.ClientTypeService > 0 {
		fmt.Fprintf(&builder, " --tos %d", configCache.ClientTypeService)
	}

	if configCache.ClientDontFragment {
		builder.WriteString(" --dont-fragment")
	}

	if configCache.ClientMaxmumSegment {
		builder.WriteString(" --set-mss")
	}

	if configCache.ClientOnlyIPv4 {
		builder.WriteString(" --version4")
	}

	if configCache.ClientOnlyIPv6 {
		builder.WriteString(" --version6")
	}

	builder.WriteString(" --get-server-output")

	stdout, stdErr, cancel, exitCodeChan, err := ExecuteAsync(filepath.Join(ToolDirGet(), "iperf3.exe"), strings.Fields(builder.String()))
	if err != nil {
		logs.Warning("iperf client startup failed, %s", err.Error())
		return nil, err
	}

	srv := new(IperfServer)
	srv.stdOut = stdout.Name()
	srv.stdErr = stdErr.Name()
	srv.cancel = cancel
	srv.running = true

	go func() {
		exitCode := <-exitCodeChan

		ReadResult(stdout.Name(), configCache.ClientLog)
		ReadResult(stdErr.Name(), configCache.ClientLog)

		logs.Info("iperf3.exe client exit code %d", exitCode)

		srv.exitCode = exitCode
		srv.running = false
	}()

	return srv, nil
}
