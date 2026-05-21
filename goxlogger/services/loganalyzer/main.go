package loganalyzer

import (
	"fmt"
)

func main() {
	server := NewGrpcServer(9000)
	err := server.Start()
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
