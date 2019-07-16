package main

import (
	"fmt"
	"github.com/pkg/errors"
)


type MyErr int32

func (myerr *MyErr) Error() string {
	return "No"
}

// empty interface
func summ2(w interface{}) (string, error) {
	switch w.(type) {
	case int:
		return "int", nil
	case string:
		return "string",nil
	default:
		return "", errors.New("unknown")
	}
}



func main(){
	fmt.Println("Heaven")
}