package cmd

import (
	"fmt"
	"runtime/debug"
	"sort"
	"time"
)

func printPanicStack(a ...any) {
	if x := recover(); x != nil {
		fmt.Println(a...)
		debug.PrintStack()
	}
}

type TimeSorter []time.Time

var _ sort.Interface = (*TimeSorter)(nil)

func (s TimeSorter) Len() int {
	return len(s)
}

func (s TimeSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s TimeSorter) Less(i, j int) bool {
	return s[i].Before(s[j])
}

/*



 */
