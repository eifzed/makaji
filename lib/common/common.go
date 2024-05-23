package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/oklog/ulid"
	"gopkg.in/mgo.v2/bson"
)

func SafelyCloseFile(f io.Closer) {
	if err := f.Close(); err != nil {
		log.Warnf("Failed to close file: %s\n", err)
	}
}

func IsDevelopment() bool {
	isLocal := os.Getenv("ISLOCAL")
	return isLocal == "1"
}

func GenerateUUIDV7() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)

	ulid := ulid.MustNew(ulid.Timestamp(t), entropy)
	return ulid.String()
}

func ComputeSHA256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func ToBsonM(s interface{}) bson.M {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	result := bson.M{}
	structToBsonM(result, val, "")
	return bson.M{"$set": result}
}

func structToBsonM(result bson.M, val reflect.Value, parentKey string) {
	valType := val.Type()

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		fieldValue := val.Field(i)
		fieldKey := field.Name

		// Use JSON tag if available
		if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
			fieldKey = strings.Split(tag, ",")[0]
		}

		if parentKey != "" {
			fieldKey = parentKey + "." + fieldKey
		}

		if !fieldValue.CanInterface() {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String:
			if fieldValue.String() != "" {
				result[fieldKey] = fieldValue.String()
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fieldValue.Int() != 0 {
				result[fieldKey] = fieldValue.Int()
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if fieldValue.Uint() != 0 {
				result[fieldKey] = fieldValue.Uint()
			}
		case reflect.Float32, reflect.Float64:
			if fieldValue.Float() != 0 {
				result[fieldKey] = fieldValue.Float()
			}
		case reflect.Bool:
			if fieldValue.Bool() {
				result[fieldKey] = fieldValue.Bool()
			}
		case reflect.Slice, reflect.Array:
			if fieldValue.Len() > 0 {
				result[fieldKey] = fieldValue.Interface()
			}
		case reflect.Map:
			if fieldValue.Len() > 0 {
				result[fieldKey] = fieldValue.Interface()
			}
		case reflect.Struct:
			nestedResult := bson.M{}
			structToBsonM(nestedResult, fieldValue, "")
			for k, v := range nestedResult {
				result[fieldKey+"."+k] = v
			}
		case reflect.Ptr:
			if !fieldValue.IsNil() {
				structToBsonM(result, fieldValue.Elem(), fieldKey)
			}
		default:
			if !fieldValue.IsZero() {
				result[fieldKey] = fieldValue.Interface()
			}
		}
	}
}
