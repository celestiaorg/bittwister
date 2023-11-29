package latency

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/celestiaorg/bittwister/xdp"
	"go.uber.org/zap"
)

type Latency struct {
	NetworkInterface *net.Interface
	Latency          time.Duration
	Jitter           time.Duration
	TcBinPath        string // default: tc
	ready            bool
}

var _ xdp.XdpLoader = (*Latency)(nil)

// Latency uses TC under the hood to impose latency and jitter on packets.
// This is a temporary solution until we have a better way to do this; probably in XDP.
func (l *Latency) Start(ctx context.Context, logger *zap.Logger) {
	if l.TcBinPath == "" {
		l.TcBinPath = "tc"
	}

	if !l.isTcInstalled() {
		logger.Fatal("tc command not found")
	}

	if err := l.deleteTc(); err != nil {
		logger.Fatal("failed to delete tc rule", zap.Error(err))
	}

	if err := l.addTc(); err != nil {
		logger.Fatal("failed to set latency/jitter using tc", zap.Error(err))
	}

	logger.Info(
		fmt.Sprintf("Latency/Jitter started with %d milliseconds latency and %d jitter on device %q",
			l.Latency.Milliseconds(),
			l.Jitter.Milliseconds(),
			l.NetworkInterface.Name,
		),
	)

	l.ready = true
	<-ctx.Done()

	// Cleanup
	if err := l.deleteTc(); err != nil {
		logger.Fatal("failed to delete tc rule", zap.Error(err))
	}

	l.ready = false
	logger.Info(fmt.Sprintf("Latency/Jitter stopped on device %q", l.NetworkInterface.Name))
}

func (l *Latency) Ready() bool {
	return l.ready
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
	return exec.Command(l.TcBinPath, "qdisc", "del", "dev", l.NetworkInterface.Name, "root").Run()
}

func (l *Latency) addTc() error {
	latencyStr := fmt.Sprintf("%dms", l.Latency.Milliseconds())
	jitterStr := fmt.Sprintf("%dms", l.Jitter.Milliseconds())
	return exec.Command(l.TcBinPath, "qdisc", "add", "dev", l.NetworkInterface.Name, "root", "netem", "delay", latencyStr, jitterStr).Run()
}

func (l *Latency) isThereTcNetEmRule() bool {
	out, err := exec.Command(l.TcBinPath, "qdisc", "show", "dev", l.NetworkInterface.Name).Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "netem")
}
