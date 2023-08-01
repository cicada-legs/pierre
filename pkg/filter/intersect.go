package filter

/*

 */

func Intersection(a, b []string) (result []string) { //get all of the common status codes, adapt later for other filtering too
	string_boo := make(map[string]bool)

	for _, thing := range a { //index = _, current element = a
		string_boo[thing] = true
	}

	for _, thing := range b {
		if _, exists := string_boo[thing]; exists {
			result = append(result, thing)
		}
	}
	return result
}
