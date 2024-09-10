package reflect

import (
	"reflect"
)

type FuncInfo struct {
	Name        string
	InputTypes  []reflect.Type
	OutputTypes []reflect.Type
	Result      any
}

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	typ := reflect.TypeOf(entity)

	numMethod := typ.NumMethod()
	res := make(map[string]FuncInfo, numMethod)
	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)

		numIn := method.Type.NumIn()
		inputTypes := make([]reflect.Type, 0, numIn)
		inputTypes = append(inputTypes, reflect.TypeOf(entity))
		inputValues := make([]reflect.Value, 0, numIn)
		inputValues = append(inputValues, reflect.ValueOf(entity))
		for j := 1; j < numIn; j++ {
			//inputType := method.Func.Type().In(j)
			inputType := method.Type.In(j)
			inputTypes = append(inputTypes, inputType)
			inputValues = append(inputValues, reflect.Zero(inputType))
		}

		numOut := method.Type.NumOut()
		outputTypes := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			outputTypes = append(outputTypes, method.Type.Out(j))
		}

		rts := method.Func.Call(inputValues)
		results := make([]any, 0, len(rts))
		for _, rt := range rts {
			results = append(results, rt.Interface())
		}

		res[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  inputTypes,
			OutputTypes: outputTypes,
			Result:      results,
		}

	}

	return res, nil
}
