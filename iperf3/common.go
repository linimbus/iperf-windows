package iperf3

func CloseWindows() {
	ClientClose()
	ServerClose()
	NotifyExit()
}
