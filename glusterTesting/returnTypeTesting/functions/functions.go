package main

func ReturnMap(inputMap map[string]int) map[string]int {

	output := make(map[string]int)
	output["test"] = inputMap["test"]
	return output
}