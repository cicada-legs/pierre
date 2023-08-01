package filter

/*
	Used for common operations for filtering and dealing with slices
*/

/*
	Intersection takes two slices of any type and returns a slice of objects that are in both slices.
	This is used for detecting and properly handling common elements in two slices.
*/
func Intersection[T comparable](a, b []T) (result []T) { //get all of the common status codes, adapt later for other filtering too
	thing_bool_map := make(map[T]bool)

	for _, thing := range a { //index = _, current element = a
		thing_bool_map[thing] = true
	}

	for _, thing := range b {
		if _, exists := thing_bool_map[thing]; exists {
			result = append(result, thing)
		}
	}
	return result
}

func Contains[T comparable](slice []T, obj T) bool {
	for _, thing := range slice {
		if thing == obj {
			return true
		}
	}
	return false
}
