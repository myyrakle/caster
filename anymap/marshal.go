package anymap

import (
	"fmt"
	"reflect"

	"github.com/myyrakle/caster/utils"
)

type AnyMap map[string]any

func isScalarType(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		return true
	default:
		return false
	}
}

func Unmarshal(source map[string]any, destination any) error {
	typeInfo := reflect.TypeOf(destination)

	// destination이 포인터가 아니면 에러
	if typeInfo.Kind() != reflect.Pointer {
		return fmt.Errorf("destination must be a pointer")
	}

	destinationPtr := destination.(*any)

	// 포인터 역참조
	dereferecedValueInfo := reflect.ValueOf(destinationPtr).Elem()
	dereferecedTypeInfo := reflect.TypeOf(*destinationPtr)

	// swap을 위한 임시 저장공간
	tmp := reflect.New(dereferecedValueInfo.Elem().Type()).Elem()
	tmp.Set(dereferecedValueInfo.Elem())

	for i := 0; i < dereferecedTypeInfo.NumField(); i++ {
		fieldInfo := dereferecedTypeInfo.Field(i)
		fieldTypeInfo := fieldInfo.Type
		fieldName := fieldInfo.Name
		fieldValue := tmp.FieldByName(fieldName)

		// json 태그 가져오기
		jsonTag := fieldInfo.Tag.Get("json")

		// bson 태그 가져오기
		bsonTag := fieldInfo.Tag.Get("bson")

		// 가져올 이름: 우선순위는 json > bson > 필드 이름 그대로
		fieldNameToUnMarshal := ""

		if jsonTag != "" {
			fieldNameToUnMarshal = jsonTag
		}

		if fieldNameToUnMarshal == "" && bsonTag != "" {
			fieldNameToUnMarshal = bsonTag
		}

		if fieldNameToUnMarshal == "" {
			fieldNameToUnMarshal = fieldName
		}

		anyValue, ok := source[fieldNameToUnMarshal]

		if !ok {
			continue
		}

		if anyValue == nil {
			continue
		}

		// anyValue 값이 nil이면 벗겨냄
		if fieldInfo.Type.Kind() == reflect.Ptr {
			fieldTypeInfo = fieldTypeInfo.Elem()
		}

		// anyValue값이 스칼라 타입이면 바로 할당
		if isScalarType(fieldTypeInfo.Kind()) {
			setScalarValue(fieldTypeInfo, &fieldValue, anyValue)

			continue
		}

		// anyValue값이 array 타입이면 재귀 호출
		if reflect.TypeOf(anyValue).Kind() == reflect.Slice {
			panic("not implemented yet")
		}

		// anyValue값이 map 타입이면 재귀 호출
		if reflect.TypeOf(anyValue).Kind() == reflect.Map {
			panic("not implemented yet")
		}
	}

	dereferecedValueInfo.Set(tmp)

	return nil
}

func setScalarValue(typeInfo reflect.Type, dest *reflect.Value, value any) {
	if dest.Kind() == reflect.Ptr {
		switch typeInfo.Kind() {
		case reflect.Int:
			dest.Set(reflect.ValueOf(utils.ToPointer(int(reflect.ValueOf(value).Int()))))
		case reflect.Int8:
			dest.Set(reflect.ValueOf(utils.ToPointer(int8(reflect.ValueOf(value).Int()))))
		case reflect.Int16:
			dest.Set(reflect.ValueOf(utils.ToPointer(int16(reflect.ValueOf(value).Int()))))
		case reflect.Int32:
			dest.Set(reflect.ValueOf(utils.ToPointer(int32(reflect.ValueOf(value).Int()))))
		case reflect.Int64:
			dest.Set(reflect.ValueOf(utils.ToPointer(int64(reflect.ValueOf(value).Int()))))
		case reflect.Uint:
			dest.Set(reflect.ValueOf(utils.ToPointer(uint(reflect.ValueOf(value).Uint()))))
		case reflect.Uint8:
			dest.Set(reflect.ValueOf(utils.ToPointer(uint8(reflect.ValueOf(value).Uint()))))
		case reflect.Uint16:
			dest.Set(reflect.ValueOf(utils.ToPointer(uint16(reflect.ValueOf(value).Uint()))))
		case reflect.Uint32:
			dest.Set(reflect.ValueOf(utils.ToPointer(uint32(reflect.ValueOf(value).Uint()))))
		case reflect.Uint64:
			dest.Set(reflect.ValueOf(utils.ToPointer(uint64(reflect.ValueOf(value).Uint()))))
		case reflect.Float32:
			dest.Set(reflect.ValueOf(utils.ToPointer(float32(reflect.ValueOf(value).Float()))))
		case reflect.Float64:
			dest.Set(reflect.ValueOf(utils.ToPointer(float64(reflect.ValueOf(value).Float()))))
		case reflect.String:
			dest.Set(reflect.ValueOf(utils.ToPointer(reflect.ValueOf(value).String())))
		case reflect.Bool:
			dest.Set(reflect.ValueOf(utils.ToPointer(reflect.ValueOf(value).Bool())))
		}
	} else {
		switch typeInfo.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dest.SetInt(int64(reflect.ValueOf(value).Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			dest.SetUint(uint64(reflect.ValueOf(value).Uint()))
		case reflect.Float32, reflect.Float64:
			dest.SetFloat(float64(reflect.ValueOf(value).Float()))
		case reflect.String:
			dest.SetString(reflect.ValueOf(value).String())
		case reflect.Bool:
			dest.SetBool(reflect.ValueOf(value).Bool())
		}
	}
}
