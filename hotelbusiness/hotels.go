// +build !solution

package hotelbusiness

import (
	"sort"
)

// Guest ...
type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

// Load ...
type Load struct {
	StartDate  int
	GuestCount int
}

// ComputeLoad ...
func ComputeLoad(guests []Guest) []Load {
	info := make(map[int]int)
	for _, g := range guests {
		info[g.CheckInDate]++
		info[g.CheckOutDate]--
	}
	type data struct {
		d, c int
	}
	scanLine := make([]data, 0, len(info))
	resCap := 0
	for key, val := range info {
		scanLine = append(scanLine, data{
			d: key,
			c: val,
		})
		if val != 0 {
			resCap++
		}
	}
	sort.Slice(scanLine, func(i, j int) bool {
		return scanLine[i].d < scanLine[j].d
	})
	cur := 0
	res := make([]Load, 0, resCap)
	for _, d := range scanLine {
		if d.c != 0 {
			cur += d.c
			res = append(res, Load{
				StartDate:  d.d,
				GuestCount: cur,
			})
		}
	}
	return res
}
