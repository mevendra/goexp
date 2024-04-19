package internal

import (
	"fmt"
)

type Report struct {
	TotalTime        int64
	NumberOfRequests int
	NumberOfSuccess  int
	Others           map[int]int

	numberOfInternalErrors int
}

func (r Report) String() string {
	var othersTotal int
	var othersReport string
	for key, value := range r.Others {
		othersTotal += value
		othersReport += fmt.Sprintf("\n\t\t\tCode %d: %d", key, value)
	}
	return fmt.Sprintf("Report: \n\tTotalTime: %dms;\n\tRequests: %d;\n\t\tSuccessfull: %d;\n\t\tOthers: %d; %s",
		r.TotalTime, r.NumberOfRequests, r.NumberOfSuccess, othersTotal, othersReport)
}
