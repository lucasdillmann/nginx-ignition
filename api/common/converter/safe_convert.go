package converter

import (
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
)

var badRequestErr = coreerror.New("One or more attributes were sent in an invalid format", true)

func Wrap[I any, O any](converter func(I) O, input I) O {
	return safeExecute(func() O {
		return converter(input)
	})
}

func Wrap2[I1 any, I2 any, O any](converter func(I1, I2) O, input1 I1, input2 I2) O {
	return safeExecute(func() O {
		return converter(input1, input2)
	})
}

func safeExecute[T any](action func() T) T {
	type result struct {
		value T
		err   any
	}
	ch := make(chan result, 1)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- result{err: err}
			}
		}()

		ch <- result{value: action()}
	}()

	res := <-ch
	if res.err != nil {
		panic(badRequestErr)
	}

	return res.value
}
