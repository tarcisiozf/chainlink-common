package loop_test

import (
	"os/exec"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	errorlogtest "github.com/smartcontractkit/chainlink-common/pkg/loop/core/services/errorlog/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	reportingplugintest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/reportingplugin/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/relayer/pluginprovider/ext/median/test"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
)

func TestMedianService(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	median := loop.NewMedianService(lggr, loop.GRPCOpts{}, func() *exec.Cmd {
		return NewHelperProcessCommand(loop.PluginMedianName, false, 0)
	}, median_test.MedianProvider(lggr), median_test.MedianContractID, median_test.DataSource, median_test.JuelsPerFeeCoinDataSource, median_test.GasPriceSubunitsDataSource, errorlogtest.ErrorLog, nil)
	hook := median.PluginService.XXXTestHook()
	servicetest.Run(t, median)

	t.Run("control", func(t *testing.T) {
		reportingplugintest.RunFactory(t, median)
	})

	t.Run("Kill", func(t *testing.T) {
		hook.Kill()

		// wait for relaunch
		time.Sleep(2 * goplugin.KeepAliveTickDuration)

		reportingplugintest.RunFactory(t, median)
	})

	t.Run("Reset", func(t *testing.T) {
		hook.Reset()

		// wait for relaunch
		time.Sleep(2 * goplugin.KeepAliveTickDuration)

		reportingplugintest.RunFactory(t, median)
	})
}

func TestMedianService_recovery(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	var limit atomic.Int32
	median := loop.NewMedianService(lggr, loop.GRPCOpts{}, func() *exec.Cmd {
		return HelperProcessCommand{
			Command: loop.PluginMedianName,
			Limit:   int(limit.Add(1)),
		}.New()
	}, median_test.MedianProvider(lggr), median_test.MedianContractID, median_test.DataSource, median_test.JuelsPerFeeCoinDataSource, median_test.GasPriceSubunitsDataSource, errorlogtest.ErrorLog, nil)
	servicetest.Run(t, median)

	reportingplugintest.RunFactory(t, median)
}
