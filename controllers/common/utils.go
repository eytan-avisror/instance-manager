/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ContainsString returns true if a given slice 'slice' contains string 's', otherwise return false
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// ContainsEqualFold returns true if a given slice 'slice' contains string 's' under unicode case-folding
func ContainsEqualFold(slice []string, s string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, s) {
			return true
		}
	}
	return false
}

func StringMD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func StringMapSliceContains(m []map[string]string, contains map[string]string) bool {
	for _, obj := range m {
		if reflect.DeepEqual(obj, contains) {
			return true
		}
	}
	return false
}

func FieldValue(path string, obj map[string]interface{}) interface{} {
	var resourceField interface{}

	p := FieldPath(path)

	if field, ok, _ := unstructured.NestedFieldCopy(obj, p...); ok {
		resourceField = field
	}

	return resourceField
}

func SetFieldValue(path string, obj map[string]interface{}, value interface{}) error {
	p := FieldPath(path)
	if err := unstructured.SetNestedField(obj, value, p...); err != nil {
		return err
	}

	return nil
}

func FieldPathString(path ...string) string {
	return strings.Join(path, ".")
}

func FieldPath(path string) []string {
	return strings.Split(path, ".")
}

func AppendUnique(slice []interface{}, i interface{}) []interface{} {
	for _, ele := range slice {
		if reflect.DeepEqual(ele, i) {
			return slice
		}
	}
	return append(slice, i)
}

func AppendUniqueIndex(slice []interface{}, i interface{}, idx string, override bool) []interface{} {
	var fieldStr, fieldStr2 string
	appendIdxVal := reflect.ValueOf(i)

	switch appendIdxVal.Kind() {
	case reflect.Map:
		for _, e := range appendIdxVal.MapKeys() {
			if strings.EqualFold(e.String(), idx) {
				fieldStr = appendIdxVal.MapIndex(e).Elem().String()
			}
		}

		for ix, ele := range slice {
			compareIdxVal := reflect.ValueOf(ele)
			for _, e := range compareIdxVal.MapKeys() {
				if strings.EqualFold(e.String(), idx) {
					fieldStr2 = compareIdxVal.MapIndex(e).Elem().String()
				}
			}
			if strings.EqualFold(fieldStr, fieldStr2) {
				if override {
					slice[ix] = i
				}
				return slice
			}
		}
	default:
		return slice
	}

	return append(slice, i)
}

func MergeSliceByUnique(sl1, sl2 []interface{}) []interface{} {
	for _, ele := range sl2 {
		sl1 = AppendUnique(sl1, ele)
	}
	return sl1
}

func MergeSliceByIndex(sl1, sl2 []interface{}, idx string, override bool) []interface{} {
	for _, ele := range sl2 {
		sl1 = AppendUniqueIndex(sl1, ele, idx, override)
	}
	return sl1
}

func StringSliceEqualFold(x []string, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for _, element := range x {
		if !ContainsEqualFold(y, element) {
			return false
		}
	}
	return true
}

func SliceEmpty(slice []string) bool {
	return len(slice) == 0
}

func MapEmpty(m map[string]string) bool {
	return len(m) == 0
}

func StringEmpty(str string) bool {
	return str == ""
}

func StringSliceEquals(x, y []string) bool {
	sort.Strings(x)
	sort.Strings(y)
	return reflect.DeepEqual(x, y)
}

func StringSliceContains(x, y []string) bool {
	for _, s := range x {
		if !ContainsString(y, s) {
			return false
		}
	}
	return true
}

func GetLastElementBy(s, sep string) string {
	sp := strings.Split(s, sep)
	return sp[len(sp)-1]
}

// ConcatenateList joins lists to strings delimited with `delimiter`
func ConcatenateList(list []string, delimiter string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(list)), delimiter), "[]")
}

func ReadFile(path string) ([]byte, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GetTimeString() string {
	n := time.Now().UTC()
	return n.Format("20060102150405")
}

// Set Difference: A - B
func Difference(a, b []string) []string {
	var diff []string
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}
