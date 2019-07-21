package main

var supportOperations = []string{"*", "/", "+", "-"}

func contains(arraySource []string, object string) bool {
	for _, v := range arraySource {
		if v == object {
			return true
		}
	}
	return false
}

func getFirstIndex(arraySource []string, object string) int {
	for i, v := range arraySource {
		if v == object {
			return i
		}
	}
	return -1
}
