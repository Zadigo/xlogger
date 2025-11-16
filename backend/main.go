package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type LogLine struct {
	// The IP address of the client
	remoteAddress string
	// The authenticated user
	remoteUser string
	// Date and time at which the request was made + TZ
	dateTime string
	// Method used for the request
	method string
	// The path of the request
	path string
	// The HTTP protocole used e.g HTTP/2.0
	protocole string
	// The request status code
	statusCode int
	// Number of bytes sent to the client (body only, not headers)
	bodyBytesSent int
	// The page from which the user came
	referrer string
	// The client's user agent
	userAgent string

	// The date part of the date time
	remoteDate string
	// the time part of the date time
	remoteTime string

	// Whether the request was successfull
	isSuccess bool
}

// Checks the value of the status code and returns
// if it was successful or not
func analyzeStatusCode(code int) bool {
	return code >= 200 && code <= 226
}

func parseLine(textValue string) (*LogLine, error) {
	logLineRegex := regexp.MustCompile(`^(\S+) - (\S+) \[([^\]]+)\] "(GET|POST|PUT|DELETE|HEAD|OPTIONS|PATCH) ([^"]+) (HTTP\/[0-9\.]+)" (\d{3}) (\d+) "([^"]*)" "([^"]*)"$`)

	var matched []string = logLineRegex.FindStringSubmatch(textValue)
	if matched == nil {
		return nil, fmt.Errorf("line is not valid")
	}

	status, _ := strconv.Atoi(matched[7])
	bytes, _ := strconv.Atoi(matched[8])

	// Parse the actual date and time from
	// the datetime string that we got

	var remoteDate string
	var remoteTime string

	dateLayout := "02/Jan/2006:15:04:05 -0700"
	parsedDate, err := time.Parse(dateLayout, matched[3])

	if err != nil {
		remoteDate = parsedDate.Format("2006-01-02")
		remoteTime = parsedDate.Format("15:04:05")
		fmt.Println(remoteDate)
	}

	fmt.Println(remoteDate)

	statusCodeAnalysis := analyzeStatusCode(status)

	return &LogLine{
		remoteAddress: matched[1],
		remoteUser:    matched[2],
		dateTime:      matched[3],
		method:        matched[4],
		path:          matched[5],
		protocole:     matched[6],
		statusCode:    status,
		bodyBytesSent: bytes,
		referrer:      matched[9],
		userAgent:     matched[10],

		remoteDate: remoteDate,
		remoteTime: remoteTime,

		isSuccess: statusCodeAnalysis,
	}, nil
}

// Reads the log file and returns an array of strings
// representing each line in the file
func readFile(path string) ([]string, error) {
	file, err := os.Open(path)

	var logs []string = make([]string, 0)
	if err != nil {
		log.Fatal("Could not open file")
		return logs, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		logs = append(logs, line)
	}

	return logs, nil
}

func main() {
	logs, err := readFile("example2.log")

	fmt.Println(logs)

	if err != nil {
		log.Fatal("Could not read file")
	}

	for i, value := range logs {
		result, err := parseLine(value)

		if err != nil {
			log.Printf("Could not parse line %v: %v", i, err.Error())
			continue
		}
		fmt.Println(result)
	}

	// var line string = `172.21.0.2 - - [08/May/2025:13:53:52 +0000] "GET /techs/societeinfo.png HTTP/1.1" 200 33239 "https://johnpm.gency313.fr/growth-marketing/tech" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"`
	// result, err := parseLine(line)

	// if err != nil {
	// 	log.Fatal("Could not parse line")
	// }

	// fmt.Println(result)
}
