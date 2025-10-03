package utils

import "fmt"

func HandleError(err error, context string, dontPanic bool) {
	if dontPanic {
		if err != nil {
			fmt.Printf("error (%s)\nactual error: %v\n", context, err)
			return
		}
	}
	if err != nil {
		panic(fmt.Sprintf("error (%s)\nactual error: %v\n", context, err))
	}
}
