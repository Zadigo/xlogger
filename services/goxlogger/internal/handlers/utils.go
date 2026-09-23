package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/gorilla/websocket"
)

func IsWebsocketClose(err error) bool {
	if websocket.IsCloseError(err,
		websocket.CloseNormalClosure,   // 1000
		websocket.CloseGoingAway,       // 1001
		websocket.CloseAbnormalClosure, // 1006
	) {
		return true
	}

	// Also catches abrupt disconnects
	// (io.EOF, reset by peer, etc.)
	if errors.Is(err, io.EOF) {
		return true
	}

	return false
}

// logFilter holds parsed, validated query parameters for filtering logs.
type logFilter struct {
	startDate  *time.Time
	endDate    *time.Time
	methods    map[string]struct{}
	search     string
	status     map[int]struct{}
	successful bool
	needsDate  bool
}
 
// newLogFilter parses and validates all filter-related query parameters once,
// up front, so bad input produces a real error instead of a silently empty result.
func newLogFilter(r *http.Request) (*logFilter, error) {
	q := r.URL.Query()
	f := &logFilter{
		search:     q.Get("search"),
		successful: q.Get("successful") == "1",
	}
 
	if v := q.Get("startDate"); v != "" {
		d, err := time.Parse("2006-01-02", v)
		if err != nil {
			return nil, fmt.Errorf("invalid startDate %q: %w", v, err)
		}
		f.startDate = &d
	}
 
	if v := q.Get("endDate"); v != "" {
		d, err := time.Parse("2006-01-02", v)
		if err != nil {
			return nil, fmt.Errorf("invalid endDate %q: %w", v, err)
		}
		// Make the end date inclusive of the whole day, not just midnight.
		d = d.Add(24*time.Hour - time.Nanosecond)
		f.endDate = &d
	}
	f.needsDate = f.startDate != nil || f.endDate != nil
 
	if v := q.Get("methods"); v != "" {
		f.methods = make(map[string]struct{})
		for value := range strings.SplitSeq(v, ",") {
			f.methods[strings.ToUpper(strings.TrimSpace(value))] = struct{}{}
		}
	}
 
	if v := q.Get("status"); v != "" {
		status := strings.Split(v, ",")
		f.status = make(map[int]struct{})
		for _, s := range status {
			code, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil {
				return nil, fmt.Errorf("invalid status %q: %w", s, err)
			}
			f.status[code] = struct{}{}
		}
	}
 
	return f, nil
}
 
// matches reports whether a log line satisfies every active filter.
// logDateTime is only used when a date filter is active; pass the zero
// value when f.needsDate is false.
func (f *logFilter) matches(log tickerapp.LogLine, logDateTime time.Time) bool {
	if f.startDate != nil && logDateTime.Before(*f.startDate) {
		return false
	}
	if f.endDate != nil && logDateTime.After(*f.endDate) {
		return false
	}
	if f.methods != nil {
		if _, ok := f.methods[strings.ToUpper(log.Method)]; !ok {
			return false
		}
	}
	if f.search != "" && !strings.Contains(log.Path, f.search) {
		return false
	}
	if f.status != nil {
		if _, ok := f.status[log.StatusCode]; !ok {
			return false
		}
	}
	if f.successful && (log.StatusCode < 200 || log.StatusCode >= 300) {
		return false
	}
	return true
}
 
// resolveQuery filters logs against every query parameter in a single pass,
// parsing each log's DateTime at most once and only when actually needed.
func resolveQuery(r *http.Request, logLines []tickerapp.LogLine) ([]tickerapp.LogLine, error) {
	f, err := newLogFilter(r)
	if err != nil {
		return nil, err
	}
 
	filtered := make([]tickerapp.LogLine, 0, len(logLines))
	for _, log := range logLines {
		var logDateTime time.Time
		if f.needsDate {
			logDateTime, err = time.Parse(tickerapp.DateLayout, log.DateTime)
			if err != nil {
				return nil, fmt.Errorf("invalid log DateTime %q: %w", log.DateTime, err)
			}
		}
		if f.matches(log, logDateTime) {
			filtered = append(filtered, log)
		}
	}
	return filtered, nil
}
