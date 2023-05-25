package common

import "fmt"

const GrpcPort = 9999

func GetGrpcAddress() string {
	return fmt.Sprintf("127.0.0.1:%d", GrpcPort)
}
