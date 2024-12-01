package godefault

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func SetDefaults(data interface{}) error {
	if data == nil {
		return fmt.Errorf("input data is nil")
	}

	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("input data must be a non-nil pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input data must point to a struct")
	}

	return setDefaults(v)
}

func setDefaults(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Handle nested structs
		if field.Kind() == reflect.Struct {
			if err := setDefaults(field); err != nil {
				return err
			}
			continue
		}

		// Check if field is zero value
		if !isZeroValue(field) {
			continue
		}

		// Check for 'omitempty' in tags
		if hasOmitempty(fieldType.Tag.Get("json")) ||
			hasOmitempty(fieldType.Tag.Get("bson")) ||
			hasOmitempty(fieldType.Tag.Get("xml")) {
			continue
		}

		// Get 'default' tag value
		defaultValue := fieldType.Tag.Get("default")
		if defaultValue == "" {
			continue
		}

		// Extract custom delimiter if specified
		delimiter := ","
		if delimTag := fieldType.Tag.Get("default_delim"); delimTag != "" {
			delimiter = delimTag
		}

		// Set the field to the default value
		if err := setFieldValue(field, defaultValue, delimiter); err != nil {
			return fmt.Errorf("failed to set default value for field '%s': %v", fieldType.Name, err)
		}
	}
	return nil
}

func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func hasOmitempty(tag string) bool {
	// Tags are comma-separated
	parts := strings.Split(tag, ",")
	for _, part := range parts[1:] { // Skip the first part as it's the tag name
		if part == "omitempty" {
			return true
		}
	}
	return false
}

func setFieldValue(field reflect.Value, value, delimiter string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Type().PkgPath() == "time" && field.Type().Name() == "Duration" {
			d, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(d))
		} else {
			i, err := strconv.ParseInt(value, 0, field.Type().Bits())
			if err != nil {
				return err
			}
			field.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u, err := strconv.ParseUint(value, 0, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case reflect.Ptr:
		if field.IsNil() {
			elemType := field.Type().Elem()
			elem := reflect.New(elemType).Elem()
			if err := setFieldValue(elem, value, delimiter); err != nil {
				return err
			}
			field.Set(elem.Addr())
		}
	case reflect.Slice:
		elemType := field.Type().Elem()
		elements := strings.Split(value, delimiter)
		slice := reflect.MakeSlice(field.Type(), 0, len(elements))
		for _, elemValueStr := range elements {
			elemValueStr = strings.TrimSpace(elemValueStr)
			elemValue := reflect.New(elemType).Elem()
			err := setFieldValue(elemValue, elemValueStr, delimiter)
			if err != nil {
				return fmt.Errorf("failed to set slice element value '%s': %v", elemValueStr, err)
			}
			slice = reflect.Append(slice, elemValue)
		}
		field.Set(slice)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Type().String())
	}
	return nil
}
