package util

import "fmt"

func DoubleToString(data float64) string {
	s := fmt.Sprintf("%f", data) 
	return s
}

func IntArrayToInterface(data []int) []interface {} {
	num := len(data)
	arrayInterface := make([]interface{}, num)
	for i := 0; i < num; i++ {
		arrayInterface[i] = IntToInterface(data[i])
	}
	return arrayInterface
}

func FloatArrayToInterface(data []float64) []interface {} {
	num := len(data)
	arrayInterface := make([]interface{}, num)
	for i := 0; i < num; i++ {
		arrayInterface[i] = FloatToInterface(data[i])
	}
	return arrayInterface
}

func StringArrayToInterface(data []string) []interface {}{
	num := len(data)
	arrayInterface := make([]interface{}, num)
	for i := 0; i < num; i++ {
		arrayInterface[i] = StringToInterface(data[i])
	}
	return arrayInterface
}

func StringToInterface(data string) interface {} {
	var dataInterface interface {} = ""
	dataInterface = data
	return dataInterface
}

func FloatToInterface(data float64) interface {} {
	var dataInterface interface {} = 0.0
	dataInterface = data
	return dataInterface
}

func IntToInterface(data int) interface {} {
	var dataInterface interface {} = 0
	dataInterface = data
	return dataInterface
}

func findIndexOfStringArray(findArray []string, findData string) int {
	var position int = 0
	for i := 0; i < len(findArray); i++ {
		if findArray[i] == findData {
			position = i
			break
		}
	}
	return position
}

func MaxLength(ArrayString []string) int {
	var len_init int = len(ArrayString)
	var max int = 0
	for i := 0; i < len_init; i++ {
		if (len(ArrayString[i]) > max) {
			max = len(ArrayString[i])
		}
	}
	return max
}