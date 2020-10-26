// +build !solution

package hogwarts

// GetCourseList ...
func GetCourseList(prereqs map[string][]string) []string {
	res := make([]string, 0, len(prereqs))
	used := make(map[string]bool)
	done := make(map[string]bool)
	var dfs func(cur string)
	dfs = func(cur string) {
		used[cur] = true
		for _, next := range prereqs[cur] {
			_, cicle := used[next]
			_, done := done[next]
			if done {
				continue
			}

			if cicle {
				panic("Cicle!")
			}

			dfs(next)
		}
		done[cur] = true
		res = append(res, cur)
	}

	for key := range prereqs {
		if _, ok := used[key]; !ok {
			dfs(key)
		}
	}
	return res
}
