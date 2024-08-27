package utils

import (
	"fmt"
	"log"
)

func LogAndPrintError(message string, err error) {
	fmt.Printf(message, err)
	log.Fatalf(message, err)
}