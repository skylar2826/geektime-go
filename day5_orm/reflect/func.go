package reflect

import "reflect"

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	typ := reflect.TypeOf(entity)
	numMethod := typ.NumMethod()

	res := make(map[string]FuncInfo, numMethod)

	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		fn := method.Func

		numIn := fn.Type().NumIn()
		input := make([]reflect.Type, 0, numIn) // 输入类型
		input = append(input, reflect.TypeOf(entity))
		inputValues := make([]reflect.Value, 0, numIn) // 输入值
		inputValues = append(inputValues, reflect.ValueOf(entity))
		for j := 1; j < numIn; j++ {
			fnInType := fn.Type().In(j)
			input = append(input, fnInType)
			inputValues = append(inputValues, reflect.Zero(fnInType))
		}
		numOut := fn.Type().NumOut()
		output := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			output = append(output, fn.Type().Out(j))
		}

		resValues := fn.Call(inputValues)
		result := make([]any, 0, len(resValues))
		for _, resValue := range resValues {
			result = append(result, resValue.Interface())
		}
		res[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  input,
			OutputTypes: output,
			Result:      result,
		}
	}
	return res, nil
}

type FuncInfo struct {
	Name        string
	InputTypes  []reflect.Type
	OutputTypes []reflect.Type
	Result      []any
}
