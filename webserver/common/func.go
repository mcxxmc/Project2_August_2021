package common

import "go.uber.org/zap"

// CheckErr checks the error and prints it out if not nil
func CheckErr(err error) {
	if err != nil {
		zap.S().Error(err)
	}
}

// PanicErr raises a panic if the error is not nil.
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}
