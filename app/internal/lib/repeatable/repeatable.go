// Package repeatable provides retry utilities for unreliable or flaky operations.
package repeatable

import "time"

// DoWithTries executes fn up to attempts times with a fixed delay between failures.
// It returns nil on first success, or the last encountered error if all attempts fail.
func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return
}
