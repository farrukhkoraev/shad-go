//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	events := map[int]int{}
	
	for _, g := range guests {
		events[g.CheckInDate] += 1
		events[g.CheckOutDate] -= 1
	}
	
	dates := []int{}
	for d := range events {
		dates = append(dates, d)
	}
	sort.Ints(dates)

	loads := []Load{}
	current := 0
	for _, d := range dates {
		if events[d] != 0 {
			current += events[d]
		    loads = append(loads, Load{d, current})
		}
	}
	return loads
}
