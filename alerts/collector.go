package alerts

import (
	"fmt"
	"reflect"
)

type Collector struct {
	alerts []Alert
}

func NewCollector() Collector {
	return Collector{
		alerts: make([]Alert, 0),
	}
}

func (c *Collector) NewAlert(alert Alert, args ...any) Alert {
	alertValue := reflect.ValueOf(alert).Elem()
	alertType := reflect.TypeOf(alert).Elem()

	fieldsSet := 0
	panicMessage := "Attempt to construct %s{} field `%s` of type `%s`%s, with `%s` at %d"

	for i, arg := range args {
		field := alertValue.Field(i)
		fieldType := field.Type()
		argValue := reflect.ValueOf(arg)
		argType := argValue.Type()

		if field.Kind() == reflect.Interface {
			if !argType.Implements(fieldType) {
				panic(fmt.Sprintf(panicMessage, alertValue.Type().Name(), alertType.Field(i).Name, fieldType, " (interfce)", argType, i+1))
			}
			field.Set(argValue)
		} else {
			if argType == fieldType {
				field.Set(argValue)
			} else {
				stringType := reflect.TypeFor[string]()
				if method, ok := argType.MethodByName("String"); ok {
					argType = stringType
					argValue = method.Func.Call([]reflect.Value{argValue})[0]
				} else if argValue.CanConvert(stringType) {
					argType = stringType
					argValue = argValue.Convert(argType)
				} else {
					panic(fmt.Sprintf(panicMessage, alertValue.Type().Name(), alertType.Field(i).Name, fieldType, "", argType, i+1))
				}
				field.Set(argValue)
			}
		}

		fieldsSet++
	}

	for i := range alertType.NumField() {
		field := alertValue.Field(i)
		if !field.IsZero() {
			continue
		}

		if defaultValue, ok := alertType.Field(i).Tag.Lookup("default"); ok {
			argValue := reflect.ValueOf(defaultValue)
			field.Set(argValue)
			fieldsSet++
		}
	}

	if alertValue.NumField() != fieldsSet {
		panicMessage := "Attempt to construct %s{} with invalid amount of arguments: expected %d, but got: %d"
		panic(fmt.Sprintf(panicMessage, alertValue.Type().Name(), alertValue.NumField(), fieldsSet))
	}

	return alertValue.Addr().Interface().(Alert)
}

func (c *Collector) Alert_(alertType Alert, args ...any) {
	c.alerts = append(c.alerts, c.NewAlert(alertType, args...))
}

func (c *Collector) AlertI_(alertType Alert) {
	c.alerts = append(c.alerts, alertType)
}

func (c *Collector) GetAlerts() []Alert {
	return c.alerts
}
