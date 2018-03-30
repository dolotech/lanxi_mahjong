package utils

import (
	"errors"
	"reflect"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {

		v = v.Elem()
		obj = v.Interface()
		v = reflect.ValueOf(obj)
	}
	t := reflect.TypeOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return errors.New("No such field: %s in obj" + name)
	}

	if !structFieldValue.CanSet() {
		return errors.New("Cannot set %s field value" + name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

type RPCConfig struct {
	ip   string
	port string
}

func (s *RPCConfig) FillStruct(m map[string]string) error {
	for k, v := range m {
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
