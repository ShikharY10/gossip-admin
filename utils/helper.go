package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/fatih/color"
)

// Generate fixed size byte array
func GenerateUniqueId(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func Encode(data []byte) string {
	hb := base64.StdEncoding.EncodeToString([]byte(data))
	return hb
}

// Decoding the base string to array of bytes
func Decode(data string) []byte {
	hb, _ := base64.StdEncoding.DecodeString(data)
	return hb
}

func ShowError(heading string, err error) {
	red := color.New(color.FgRed).PrintfFunc()
	white := color.New(color.FgWhite).PrintfFunc()
	red(heading + " : ")
	white(err.Error())
	fmt.Println("")
}

func ShowSucces(message string, major bool) {
	var messagePrinter func(format string, a ...interface{})
	if major {
		messagePrinter = color.New(color.FgBlue, color.Bold, color.Underline).PrintfFunc()
	} else {
		messagePrinter = color.New(color.FgWhite).PrintfFunc()
	}
	head := color.New(color.FgGreen).PrintfFunc()
	head("[SUCCESS] -> ")
	messagePrinter(message)
	fmt.Println("")
}
