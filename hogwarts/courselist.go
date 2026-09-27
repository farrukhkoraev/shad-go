//go:build !solution

package hogwarts

// {A : [B C], B: [D], C: []}
// [A, B, D,  ]

type VisitFlag int
const (
	NotVisited VisitFlag = iota
	Partially
	Fully
)

func contains(list []string, target string) bool {
	found := false
	for _, v := range list {
		if v == target {
			found = true
			break
		}
	}
	return found
}

func GetCourseList(prereqs map[string][]string) []string {
	result := []string{}
    visited := map[string]VisitFlag{}

	var visit func(string, []string) []string
	visit = func(course string, result []string) []string {

		switch visited[course] {
		case Fully:
			return result
		case Partially:
			panic("cycle detected")
		case NotVisited:
			visited[course] = Partially
		    for _, c := range prereqs[course] {
				result = visit(c, result)
		    }
			visited[course] = Fully
			return append(result, course)
		default:
			return result
		}

	}
	
	for k, _ := range prereqs {
		visited[k] = NotVisited
	}
	
	for k, _ := range prereqs {
		result = visit(k, result)
	}
	
	return result
}
