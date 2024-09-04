package reflect

import (
	"reflect"
)

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	/*
		拿到实体上的方法
		拿到输入类型， 输入值
		拿到输出类型
		fn.call(输入值) 拿到结果
	*/
	typ := reflect.TypeOf(entity)
	numMethod := typ.NumMethod()
	methods := make(map[string]FuncInfo, numMethod)

	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		//fn := method.Func.Type()

		//numIn := fn.NumIn()
		numIn := method.Type.NumIn()
		input := make([]reflect.Type, 0, numIn) // input类型
		input = append(input, reflect.TypeOf(entity))
		inputValue := make([]reflect.Value, 0, numIn) // input值
		inputValue = append(inputValue, reflect.ValueOf(entity))
		for j := 1; j < numIn; j++ {
			//inputType := fn.In(j)
			inputType := method.Type.In(j)
			input = append(input, inputType)
			inputValue = append(inputValue, reflect.Zero(inputType))
		}

		//numOut := fn.NumOut()
		numOut := method.Type.NumOut()
		output := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			//output = append(output, fn.Out(j))
			output = append(output, method.Type.Out(j))
		}

		//resValues := method.Func.Call(inputValue)
		resValues := method.Func.Call(inputValue)
		res := make([]any, 0, len(resValues))
		for _, resValue := range resValues {
			res = append(res, resValue.Interface())
		}

		methods[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  input,
			OutputTypes: output,
			Result:      res,
		}
	}
	return methods, nil
}

type FuncInfo struct {
	Name        string
	InputTypes  []reflect.Type
	OutputTypes []reflect.Type
	Result      []any
}
