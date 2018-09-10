package main

func ReturnMap(inputMap map[string]int) map[string]int {

	var output map[string]int
	output["test"] = inputMap["test"]
	return output
}