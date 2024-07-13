package util

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
)

func StructToStruct(src any, dest any) error {
	jsonBytes, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, dest)
	if err != nil {
		return fmt.Errorf("failed to convert json %s: err=%s", string(jsonBytes), err)
	}

	return nil
}

func CompareMapsIgnoringNulls(a map[string]interface{}, b map[string]interface{}) bool {
	aNoNulls := make(map[string]interface{})
	maps.Copy(aNoNulls, a)

	bNoNulls := make(map[string]interface{})
	maps.Copy(bNoNulls, b)

	removeNullValuesFromMap(aNoNulls)
	removeNullValuesFromMap(bNoNulls)

	aAsJson, _ := json.Marshal(aNoNulls)
	bAsJson, _ := json.Marshal(bNoNulls)

	fmt.Printf("deletethis: aAsJson=%v, bAsJson=%v\n", string(aAsJson), string(bAsJson))
	return reflect.DeepEqual(aNoNulls, bNoNulls)
}

func removeNullValuesFromMap(m map[string]interface{}) {
	for k, v := range m {
		fmt.Printf("deletethis: testing k=%q,v=%v\n", k, v)
		if v == nil {
			fmt.Printf("deletethis: deleting k=%v\n", k)
			delete(m, k)
		} else if nestedMap, ok := v.(map[string]interface{}); ok {
			if len(nestedMap) == 0 {
				delete(m, k)
			}
			removeNullValuesFromMap(nestedMap)
		} else if nestedSliceOfMaps, ok := v.([]map[string]interface{}); ok {
			for _, m := range nestedSliceOfMaps {
				removeNullValuesFromMap(m)
			}

		} else if nestedSlice, ok := v.([]interface{}); ok {
			for _, sliceElem := range nestedSlice {
				if asMap, ok := sliceElem.(map[string]interface{}); ok {
					removeNullValuesFromMap(asMap)
				}
			}

		} else {
			fmt.Printf("deletethis: not nil, not a map, v:%v\n", reflect.TypeOf(v))
		}
	}
}
