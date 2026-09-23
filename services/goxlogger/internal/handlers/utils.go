package handlers

import (
	"errors"
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

const (
	DirectionAfter = "after"
	DirectionBefore = "before"
)

func queryDate(dateValue string, logs []tickerapp.LogLine, direction string) ([]tickerapp.LogLine, error) {
	var filteredLogs []tickerapp.LogLine

	date, err := time.Parse(time.RFC3339, dateValue)
	if err != nil {
		return nil, err
	}

	for _, log := range logs {
		logDateTime, err := time.Parse(time.RFC1123, log.DateTime)
		if err != nil {
			return nil, err
		}

		if direction == DirectionAfter {
			if logDateTime.After(date) || logDateTime.Equal(date) {
				filteredLogs = append(filteredLogs, log)
			}
		} else if direction == DirectionBefore {
			if logDateTime.Before(date) || logDateTime.Equal(date) {
				filteredLogs = append(filteredLogs, log)
			}
		}
	}
	return filteredLogs, nil
}

func queryMethods(methods string, logs []tickerapp.LogLine) []tickerapp.LogLine {
	if methods == "" {
		return logs
	}

	methodsList := strings.Split(methods, ",")
	var filteredLogs []tickerapp.LogLine
	for _, log := range logs {
		for _, method := range methodsList {
			if log.Method == method {
				filteredLogs = append(filteredLogs, log)
				break
			}
		}
	}
	return filteredLogs
}

func querySearch(search string, logs []tickerapp.LogLine) []tickerapp.LogLine {
	if search == "" {
		return logs
	}

	var filteredLogs []tickerapp.LogLine
	for _, log := range logs {
		if strings.Contains(log.Path, search) {
			filteredLogs = append(filteredLogs, log)
		}
	}
	return filteredLogs
}

func queryStatus(status string, logs []tickerapp.LogLine) []tickerapp.LogLine {
	if status == "" {
		return logs
	}

	var filteredLogs []tickerapp.LogLine
	for _, log := range logs {
		if statusCode, err := strconv.Atoi(status); err == nil && log.StatusCode == statusCode {
			filteredLogs = append(filteredLogs, log)
		}
	}
	return filteredLogs
}

func querySuccessful(successful string, logs []tickerapp.LogLine) []tickerapp.LogLine {
	if successful == "" {
		return logs
	}

	var filteredLogs []tickerapp.LogLine
	success, err := strconv.ParseBool(successful)
	if err != nil {
		return logs
	}
	for _, log := range logs {
		if (log.StatusCode >= 200 && log.StatusCode < 300) == success {
			filteredLogs = append(filteredLogs, log)
		}
	}
	return filteredLogs
}

func resolveQuery(r *http.Request, logs []tickerapp.LogLine) ([]tickerapp.LogLine, error) {
	filteredLogs := logs

	startDate := r.URL.Query().Get("startDate")
	if startDate != "" {
		var err error
		filteredLogs, err = queryDate(startDate, filteredLogs, DirectionAfter)
		if err != nil {
			return []tickerapp.LogLine{}, err
		}
	}

	endDate := r.URL.Query().Get("endDate")
	if endDate != "" {
		var err error
		filteredLogs, err = queryDate(endDate, filteredLogs, DirectionBefore)
		if err != nil {
			return []tickerapp.LogLine{}, err
		}
	}

	methods := r.URL.Query().Get("methods")
	if methods != "" {
		filteredLogs = queryMethods(methods, filteredLogs)
	}
	
	search := r.URL.Query().Get("search")
	if search != "" {
		filteredLogs = querySearch(search, filteredLogs)
	}

	
	status := r.URL.Query().Get("status")
	if status != "" {
		filteredLogs = queryStatus(status, filteredLogs)
	}

	successfull := r.URL.Query().Get("successful")
	if successfull != "" {
		filteredLogs = querySuccessful(successfull, filteredLogs)
	}

	return filteredLogs, nil
}
