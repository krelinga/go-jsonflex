// Package jsonflex provides utilities for working with JSON data in a flexible,
// type-safe manner. It offers type conversion functions and utilities for
// extracting values from JSON objects and arrays with proper error handling.
package jsonflex

import (
	"fmt"
	"math"
)

// Object represents a JSON object as a map with string keys and any values.
// This type alias provides semantic clarity when working with JSON object data.
type Object = map[string]any

// Array represents a JSON array as a slice of any values.
// This type alias provides semantic clarity when working with JSON array data.
type Array = []any

// Number represents a JSON number as a float64.
// JSON numbers are represented as float64 in Go's standard library.
type Number = float64

// Converter is a generic function type that converts an interface{} value
// to a specific type T, returning an error if the conversion fails.
// Converters are used throughout this package to provide type-safe value extraction.
type Converter[T any] func(any) (T, error)

// AsFloat64 returns a Converter that converts a value to float64.
// It accepts float64 values and returns an error for nil or other types.
// This is the primary converter for JSON numbers.
func AsFloat64() Converter[float64] {
	return func(v any) (float64, error) {
		if v == nil {
			return 0, fmt.Errorf("cannot convert nil to float64")
		}
		if f, ok := v.(float64); ok {
			return f, nil
		}
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

// AsString returns a Converter that converts a value to string.
// It accepts string values and returns an error for nil or other types.
// This converter is used for extracting JSON string values.
func AsString() Converter[string] {
	return func(v any) (string, error) {
		if v == nil {
			return "", fmt.Errorf("cannot convert nil to string")
		}
		if s, ok := v.(string); ok {
			return s, nil
		}
		return "", fmt.Errorf("cannot convert %T to string", v)
	}
}

// AsBool returns a Converter that converts a value to bool.
// It accepts bool values and returns an error for nil or other types.
// This converter is used for extracting JSON boolean values.
func AsBool() Converter[bool] {
	return func(v any) (bool, error) {
		if v == nil {
			return false, fmt.Errorf("cannot convert nil to bool")
		}
		if b, ok := v.(bool); ok {
			return b, nil
		}
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

// AsInt32 returns a Converter that converts a value to int32.
// It first converts the value to float64 using AsFloat64, then checks if the
// result can be safely converted to int32 without loss of precision.
// The value must be within the int32 range and be a whole number.
func AsInt32() Converter[int32] {
	return func(v any) (int32, error) {
		f, err := AsFloat64()(v)
		if err != nil {
			return 0, err
		}
		if f >= float64(math.MinInt32) && f <= float64(math.MaxInt32) && f == float64(int32(f)) {
			return int32(f), nil
		}
		return 0, fmt.Errorf("cannot convert %T to int32", v)
	}
}

// AsObject returns a Converter that converts a value to a type T that is based on Object.
// The type constraint ~Object allows for custom types that have Object as their underlying type.
// It accepts Object values (map[string]any) and returns an error for nil or other types.
// This is useful for creating custom object types with specific methods while maintaining type safety.
func AsObject[T ~Object]() Converter[T] {
	return func(v any) (T, error) {
		if v == nil {
			return T{}, fmt.Errorf("cannot convert nil to Object")
		}
		obj, ok := v.(Object)
		if !ok {
			return T{}, fmt.Errorf("cannot convert %T to Object", v)
		}
		return T(obj), nil
	}
}

// AsArray returns a Converter that converts a value to a slice of type T.
// It takes a valueConv Converter[T] to convert each element of the array.
// The function accepts Array values ([]any) and applies the valueConv to each element,
// returning an error if the input is not an array or if any element conversion fails.
// This enables type-safe extraction of homogeneous arrays from JSON data.
func AsArray[T any](valueConv Converter[T]) Converter[[]T] {
	return func(v any) ([]T, error) {
		if v == nil {
			return nil, fmt.Errorf("cannot convert nil to Array")
		}
		arr, ok := v.([]any)
		if !ok {
			return nil, fmt.Errorf("cannot convert %T to Array", v)
		}
		result := make([]T, len(arr))
		for i, item := range arr {
			converted, err := valueConv(item)
			if err != nil {
				return nil, fmt.Errorf("error converting item %d: %w", i, err)
			}
			result[i] = converted
		}
		return result, nil
	}
}

// AsAny returns a Converter that accepts any value without conversion.
// This converter always succeeds and returns the input value unchanged.
// It's useful when you want to extract a value but don't need type conversion,
// or as a fallback in generic scenarios.
func AsAny() Converter[any] {
	return func(v any) (any, error) {
		return v, nil
	}
}

// GetField extracts a field from an Object and converts it to type T using the provided Converter.
// It takes an Object, a field key, and a Converter[T] to apply to the field value.
// Returns an error if the object is nil, the field doesn't exist, or the conversion fails.
// This is the primary function for type-safe field extraction from JSON objects.
func GetField[T any](obj Object, key string, conv Converter[T]) (T, error) {
	if obj == nil {
		var zero T
		return zero, fmt.Errorf("cannot access field %q on nil object", key)
	}
	value, exists := obj[key]
	if !exists {
		var zero T
		return zero, fmt.Errorf("field %q does not exist in object", key)
	}
	return conv(value)
}

// FromArray converts an Array to a slice of type T using the provided Converter.
// This is a convenience function that wraps AsArray for direct array conversion.
// It takes an Array and a Converter[T], returning a slice of T or an error.
// Use this when you have an Array value and want to convert it directly without
// going through the GetField function.
func FromArray[T any](arr Array, conv Converter[T]) ([]T, error) {
	return AsArray(conv)(arr)
}
