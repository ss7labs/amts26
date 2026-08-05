package common

import (
	"os"

	"github.com/pkg/errors"
)

func FileExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

func Must(err error) {
	if err != nil {
		panic(errors.WithStack(err))
	}
}
func Must2(v interface{}, err error) interface{} {
	Must(err)
	return v
}
