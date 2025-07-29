package jsonflex

import (
	"fmt"
	"math"
)

type Object = map[string]any

type Array = []any

type Number = float64

type Converter[T any] func(any) (T, error)

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

func AsAny() Converter[any] {
	return func(v any) (any, error) {
		return v, nil
	}
}

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

func FromArray[T any](arr Array, conv Converter[T]) ([]T, error) {
	return AsArray(conv)(arr)
}
