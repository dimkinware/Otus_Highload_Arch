package main

import (
	"fmt"
	"time"
)

func mainPlayground() {

	var tm, err = time.Parse(time.DateOnly, "2017-02-01")
	if err != nil {
		panic(err)
	}
	var ms = tm.UnixMilli()
	fmt.Println(ms)

	var newTm = time.UnixMilli(ms)
	var newTmFormatted = newTm.Format(time.DateOnly)
	fmt.Println(newTmFormatted)

}
