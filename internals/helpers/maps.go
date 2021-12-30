package helpers

import "strconv"

type Convertions struct{}

func GetMapUintKeys(m map[string]uint) []uint {
	temp := []uint{}

	for key := range m {

		value, _ := strconv.Atoi(key)
		temp = append(temp, uint(value))

	}
	return temp
}

func GetMapStringstKeys(m map[string]string) []string {
	temp := []string{}

	for key := range m {

		temp = append(temp, key)

	}
	return temp
}

// func MapValues(m map[interface{}]interface{}) []interface{} {
// 	temp := make([]interface{}, 0)

// 	for _, value := range m {

// 		temp = append(temp, value.(reflect.Kind))
// 	}

// }
