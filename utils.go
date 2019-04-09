package main

// in check if a slice has a string
func in(str string, list []string) bool {
	for _, e := range list {
		if str == e {
			return true
		}
	}
	return false
}
