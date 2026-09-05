package simulation

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	ModeHeader  = "X-Demo-Simulation"
	DelayHeader = "X-Demo-Simulation-Delay"

	MinDelay = 250 * time.Millisecond
	MaxDelay = 2 * time.Second
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Torus must always be able to evaluate backend health
		// independently of request simulation.
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		mode, err := parseRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch mode {
		case ModeNormal:
			next.ServeHTTP(w, r)

		case ModeSlow:
			delay, err := parseDelay(r.Header.Get(DelayHeader))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			time.Sleep(delay)
			next.ServeHTTP(w, r)

		case ModeError:
			http.Error(
				w,
				"simulated backend error",
				http.StatusServiceUnavailable,
			)

		default:
			http.Error(
				w,
				"unsupported simulation mode",
				http.StatusBadRequest,
			)
		}
	})
}

func parseRequest(r *http.Request) (Mode, error) {
	value := r.Header.Get(ModeHeader)

	if value == "" {
		return ModeNormal, nil
	}

	return ParseMode(value)
}

func parseDelay(value string) (time.Duration, error) {
	if value == "" {
		return 750 * time.Millisecond, nil
	}

	delayMS, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid simulation delay %q", value)
	}

	delay := time.Duration(delayMS) * time.Millisecond

	if delay < MinDelay || delay > MaxDelay {
		return 0, fmt.Errorf(
			"simulation delay must be between %dms and %dms",
			MinDelay/time.Millisecond,
			MaxDelay/time.Millisecond,
		)
	}

	return delay, nil
}
