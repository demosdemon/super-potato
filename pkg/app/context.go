package app

import (
	"context"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func CancelOnSignal(ctx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, len(signals))
	signal.Notify(ch, signals...)

	go func() {
		select {
		case sig := <-ch:
			logrus.WithField("signal", sig).Debug("received signal")
		case <-ctx.Done():
			logrus.WithField("err", ctx.Err()).Debug("context done")
		}

		signal.Stop(ch)
		close(ch)
		cancel()
	}()

	return ctx, cancel
}
