package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func WrapWithSignature(inputStruct interface{}, apiSecret string, timestamp int64) interface{} {
	v := reflect.ValueOf(inputStruct)
	t := reflect.TypeOf(inputStruct)

	pairs := make([]string, 0, v.NumField()+1)

	result := make(map[string]interface{}, v.NumField()+1)

	for i := 0; i < v.NumField(); i++ {
		key := t.Field(i).Name

		// Prefer JSON tag if available, fall back to field name
		if jsonTag := t.Field(i).Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			key = jsonTag
		}

		val := v.Field(i).Interface()

		pairs = append(pairs, fmt.Sprintf("%s=%v", key, val))

		result[key] = val
	}

	ts := strconv.FormatInt(timestamp, 10)
	pairs = append(pairs, "timestamp="+ts)
	result["timestamp"] = timestamp

	raw := strings.Join(pairs, "&")

	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(raw))
	sig := hex.EncodeToString(mac.Sum(nil))

	result["timestamp"] = timestamp
	result["signature"] = sig

	return result
}
