package shutdown

import (
	"context"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Group interface {
	Add(name string, timeout time.Duration) (
		ctx context.Context, done chan struct{})
	size() int
	shutdown(ctx context.Context, logger logging.Logger) (incomplete int)
}

type group struct {
	prefix   string
	routines []routine
}

func NewGroup(prefix string) Group {
	return &group{
		prefix: prefix,
	}
}

func (g *group) Add(name string, timeout time.Duration) (ctx context.Context, done chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	done = make(chan struct{})
	routine := routine{
		name:    name,
		cancel:  cancel,
		done:    done,
		timeout: timeout,
	}
	g.routines = append(g.routines, routine)
	return ctx, done
}

func (g *group) size() int { return len(g.routines) }

func (g *group) shutdown(ctx context.Context, logger logging.Logger) (incomplete int) {
	completed := make(chan bool)

	for _, r := range g.routines {
		go func(r routine) {
			if err := r.shutdown(ctx); err != nil {
				logger.Warn(g.prefix + err.Error() + " ⚠️")
				completed <- false
			} else {
				logger.Info(g.prefix + r.name + " terminated ✔️")
				completed <- err == nil
			}
		}(r)
	}

	for range g.routines {
		c := <-completed
		if !c {
			incomplete++
		}
	}

	return incomplete
}
