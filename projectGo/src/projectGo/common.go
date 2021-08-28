package projectGo

import (
	"fmt"
)


// CheckErr checks the error and prints it out if not nil
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func mapInt2Bool(n int) bool {
	if n == 0 {
		return true
	}
	return false
}

func mapBool2Int(b bool) int {
	if b == true {
		return 0
	}
	return 1
}
