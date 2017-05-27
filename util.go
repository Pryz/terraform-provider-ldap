package main

import (
	"encoding/json"
)

func marshalObjectClasses(set map[string]struct{}) string {
	list := make([]string, 0, len(set))
	for item := range set {
		list = append(list, item)
	}
	data, err := json.Marshal(list)
	if err != nil {
		return "[]"
	}
	s := string(data)
	return s
}

func unmarshalObjectClasses(s string) map[string]struct{} {
	var list []string
	err := json.Unmarshal([]byte(s), &list)
	set := make(map[string]struct{}, len(list))
	if err != nil {
		return set
	}
	for _, item := range list {
		set[item] = struct{}{}
	}
	return set
}
