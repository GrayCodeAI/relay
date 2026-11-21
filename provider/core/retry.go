package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/GrayCodeAI/relay/types"
)

// RetryConfig controls retry behavior for HTTP clients.
// It embeds types.RetryConfig for the core fields and adds RetryOn for
// HTTP-status-code–driven retry decisions.
type RetryConfig struct {
	types.RetryConfig
	RetryOn []int // HTTP status codes to retry on
}

// NewRetryConfig constructs a RetryConfig from core fields and optional
// HTTP status codes to retry on.
func NewRetryConfig(maxRetries int, baseDelay, maxDelay time.Duration, retryOn ...int) RetryConfig {
	return RetryConfig{
		RetryConfig: types.RetryConfig{MaxRetries: maxRetries, BaseDelay: baseDelay, MaxDelay: maxDelay},
		RetryOn:     retryOn,
	}
}

// DefaultRetryConfig returns sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return NewRetryConfig(
		DefaultMaxRetries, DefaultBaseDelay, DefaultMaxDelay,
		429, 500, 502, 503, 529,
	)
}

// ShouldRetry checks if a status code is retryable.
func (rc RetryConfig) ShouldRetry(statusCode int) bool {
	for _, code := range rc.RetryOn {
		if code == statusCode {
			return true
		}
	}
	return false
}

// backoffDelay calculates delay with exponential backoff + jitter.
// Respects Retry-After headers from HTTP responses when present.
func (rc RetryConfig) backoffDelay(attempt int, resp *http.Response) time.Duration {
	if rc.MaxDelay <= 0 {
		return 0
	}
	// Respect Retry-After header if present
	if resp != nil {
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.ParseUint(ra, 10, 64); err == nil || errors.Is(err, strconv.ErrRange) {
				if secs > uint64(rc.MaxDelay/time.Second) {
					return rc.MaxDelay
				}
				return time.Duration(secs) * time.Second
			}
			if t, err := http.ParseTime(ra); err == nil {
				d := time.Until(t)
				if d > rc.MaxDelay {
					d = rc.MaxDelay
				}
				if d > 0 {
					return d
				}
			}
		}
	}

	// Use shared exponential backoff with full jitter
	return types.BackoffDelay(attempt, rc.BaseDelay, rc.MaxDelay)
}

var retryDelayRe = regexp.MustCompile(`(?i)(?:retry|try again)\s+(?:in|after)\s+(\d+(?:\.\d+)?)\s*(ms|milliseconds?|s|seconds?)`)

// parseRetryDelay extracts a delay hint bounded by the configured backoff limit.
func parseRetryDelay(errMsg string, maxDelay time.Duration) time.Duration {
	m := retryDelayRe.FindStringSubmatch(errMsg)
	if m == nil || maxDelay <= 0 {
		return 0
	}
	val, err := strconv.ParseFloat(m[1], 64)
	unit := time.Second
	if strings.EqualFold(m[2][:1], "m") {
		unit = time.Millisecond
	}
	if val >= float64(maxDelay)/float64(unit) {
		return maxDelay
	}
	if err != nil {
		return 0
	}
	return time.Duration(val * float64(unit))
}

// DoWithRetry executes an HTTP request with retry logic.
//
// Note: DoWithRetry operates at the transport layer, before
// formatAPIError constructs *RelayError. Structured-error awareness
// lives in the fallback chain (fallback.go:240-244) where
// *RelayError.IsRetriable() / IsAuthError() drive provider
// rotation. DoWithRetry only needs the raw transport status code
// and the underlying network error to decide whether to retry the
// same request.
func DoWithRetry(ctx context.Context, httpClient *http.Client, req *http.Request, rc RetryConfig, logger *slog.Logger) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response

	for attempt := 0; attempt <= rc.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := rc.backoffDelay(attempt-1, lastResp)
			if lastErr != nil {
				delay = max(delay, parseRetryDelay(lastErr.Error(), rc.MaxDelay))
			}
			logger.Debug(
				"retrying request",
				"attempt", attempt, "max", rc.MaxRetries,
				"delay", delay, "host", req.URL.Host,
			)
			// Use time.NewTimer + Stop instead of time.After to avoid leaking
			// the timer in the runtime when ctx is cancelled before the delay
			// elapses. time.After allocates a timer that lives until it fires,
			// even if the caller has already moved on.
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
		}

		// Clone request body for retry (body may have been consumed).
		// If a Body is set but GetBody is missing, retries would silently
		// send a drained body and likely fail with a confusing provider
		// error; surface that immediately instead.
		retryReq := req.Clone(ctx)
		if req.Body != nil && req.GetBody == nil {
			return nil, fmt.Errorf("retry requires GetBody set on *http.Request (setBody=%T)", req.Body)
		}
		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("failed to clone request body: %w", err)
			}
			retryReq.Body = body
		}

		resp, err := httpClient.Do(retryReq)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			// net/http includes the request URL in transport errors. Keep the
			// error type and cause, but omit private paths and credentials.
			var urlErr *url.Error
			if errors.As(err, &urlErr) {
				safe := *urlErr
				safe.URL = req.URL.Scheme + "://" + req.URL.Host
				err = &safe
			}
			lastResp = nil
			lastErr = err
			logger.Warn("request failed", "attempt", attempt, "error", err)
			continue
		}

		if !rc.ShouldRetry(resp.StatusCode) {
			return resp, nil
		}

		lastResp = resp
		lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
		logger.Warn("retryable status", "attempt", attempt, "status", resp.StatusCode)
		_ = resp.Body.Close()
	}

	return nil, fmt.Errorf("relay: %s: max retries (%d) exceeded: %w", req.URL.Host, rc.MaxRetries, lastErr)
}
