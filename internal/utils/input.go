package utils

import (
	"fmt"
	"strings"
)

func AskConfirm(message string) bool {
	var answer string

	fmt.Printf("%s [Y/n]: ", message)

	//nolint:errcheck
	fmt.Scanf("%s", &answer)

	if strings.ToUpper(answer) == "Y" {
		return true
	}

	fmt.Println("Aborted.")
	return false
}
