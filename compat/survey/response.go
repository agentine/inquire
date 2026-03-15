package survey

import (
	"fmt"
	"reflect"
)

func setResponse(response interface{}, answer interface{}) error {
	rv := reflect.ValueOf(response)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("survey: response must be a pointer")
	}

	rv = rv.Elem()
	av := reflect.ValueOf(answer)

	if !av.Type().AssignableTo(rv.Type()) {
		// Try common conversions
		if av.Type().ConvertibleTo(rv.Type()) {
			rv.Set(av.Convert(rv.Type()))
			return nil
		}
		return fmt.Errorf("survey: cannot assign %T to %s", answer, rv.Type())
	}

	rv.Set(av)
	return nil
}
