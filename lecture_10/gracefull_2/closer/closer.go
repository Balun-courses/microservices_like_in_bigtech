package closer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
)

var globalCloser = New()

func Add(f ...CloseFunc) {
	globalCloser.Add(f...)
}

func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

type closer struct {
	mu    sync.Mutex
	once  sync.Once
	funcs []CloseFunc
}

// CloseFunc - smth that close obj
type CloseFunc func(ctx context.Context) error

// New - returns closer
func New(sig ...os.Signal) *closer {
	c := &closer{}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll(context.Background())
		}()
	}
	return c
}

func (c *closer) Add(f ...CloseFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f...)
}

func (c *closer) CloseAll(ctx context.Context) error {
	var closeErr error
	c.once.Do(func() {
		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		wg := sync.WaitGroup{}
		for _, f := range funcs {
			wg.Add(1)
			go func(f CloseFunc) {
				defer wg.Done()
				errs <- f(ctx)
			}(f)
		}

		go func() {
			wg.Done()
			close(errs)
		}()

		msgs := make([]string, 0, len(funcs))
	Loop:
		for {
			select {
			case err, ok := <-errs:
				if !ok {
					break Loop
				}
				if err != nil {
					msgs = append(msgs, fmt.Sprintf("[!] %v", err))
				}
			case <-ctx.Done():
				break Loop
			}
		}

		if len(msgs) > 0 {
			closeErr = fmt.Errorf(
				"shutdown finished with error(s): \n%s",
				strings.Join(msgs, "\n"),
			)
		}
	})

	return closeErr
}
