package converter

import (
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
)

var badRequestErr = coreerror.New("One or more attributes were sent in an invalid format", true)

func Wrap[I any, O any](converter func(I) *O, input I) *O {
	if err := recover(); err != nil {
		panic(badRequestErr)
	}

	return converter(input)
}

func Wrap2[I1 any, I2 any, O any](converter func(I1, I2) *O, input1 I1, input2 I2) *O {
	if err := recover(); err != nil {
		panic(badRequestErr)
	}

	return converter(input1, input2)
}
