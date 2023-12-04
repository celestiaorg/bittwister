package latency

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/celestiaorg/bittwister/xdp"
)

type Latency struct {
	NetworkInterface *net.Interface
	Latency          time.Duration
	Jitter           time.Duration
	TcBinPath        string // default: tc
}

var _ xdp.XdpLoader = (*Latency)(nil)

// Latency uses TC under the hood to impose latency and jitter on packets.
// This is a temporary solution until we have a better way to do this; probably in XDP.
func (l *Latency) Start() (xdp.CancelFunc, error) {
	if l.TcBinPath == "" {
		l.TcBinPath = "tc"
	}

	if !l.isTcInstalled() {
		return nil, fmt.Errorf("tc command not found")
	}

	if err := l.deleteTc(); err != nil {
		return nil, err
	}

	if err := l.addTc(); err != nil {
		return nil, fmt.Errorf("set latency/jitter using tc: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
	}()

	cancelFunc := xdp.CancelFunc(func() error {
		if err := l.deleteTc(); err != nil {
			return err
		}
		cancel()
		return nil
	})

	return cancelFunc, nil
}

// Check if the tc command is installed.
func (l *Latency) isTcInstalled() bool {
	_, err := exec.LookPath(l.TcBinPath)
	return err == nil
}

func (l *Latency) deleteTc() error {
	if !l.isThereTcNetEmRule() {
		return nil
	}
	out, err := exec.Command(l.TcBinPath, "qdisc", "del", "dev", l.NetworkInterface.Name, "root").CombinedOutput()
	if err != nil {
		return fmt.Errorf("delete tc rule: %w, output: `%s`", err, string(out))
	}
	return nil
}

func (l *Latency) addTc() error {
	latencyStr := fmt.Sprintf("%dms", l.Latency.Milliseconds())
	jitterStr := fmt.Sprintf("%dms", l.Jitter.Milliseconds())
	out, err := exec.Command(l.TcBinPath, "qdisc", "add", "dev", l.NetworkInterface.Name, "root", "netem", "delay", latencyStr, jitterStr).CombinedOutput()
	if err != nil {
		return fmt.Errorf("add tc rule: %w, output: `%s`", err, string(out))
	}
	return nil
}

func (l *Latency) isThereTcNetEmRule() bool {
	out, err := exec.Command(l.TcBinPath, "qdisc", "show", "dev", l.NetworkInterface.Name).CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "netem")
}
