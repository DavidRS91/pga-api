package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

// borrowed from https://stackoverflow.com/questions/12753805/type-converting-slices-of-interfaces
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func InterfaceMap(i interface{}) map[string]interface{} {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Map {
		panic("InterfaceToMap() given a non-map type")
	}
	m := make(map[string]interface{})
	for _, k := range v.MapKeys() {
		strct := v.MapIndex(k)
		key := k.Interface().(string)
		m[string(key)] = strct.Interface()
	}
	return m
}

func InterfaceInt(i interface{}) (int, error) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.String {
		s := i.(string)
		return strconv.Atoi(s)
	}
	if v.Kind() == reflect.Float64 {
		f := i.(float64)
		return int(f), nil
	}
	return 0, fmt.Errorf("unknown interface type")
}
