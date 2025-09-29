package utils

import "fmt"

func HandleError(err error, context string) {
	if err != nil {
		panic(fmt.Sprintf("error (%s)\nactual error: %v\n", context, err))
	}
}
