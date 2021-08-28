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
