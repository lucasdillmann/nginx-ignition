package converter

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/i18n"
)

func Wrap[I, O any](ctx context.Context, converter func(I) O, input I) O {
	return safeExecute(ctx, func() O {
		return converter(input)
	})
}

func Wrap2[I1, I2, O any](ctx context.Context, converter func(I1, I2) O, input1 I1, input2 I2) O {
	return safeExecute(ctx, func() O {
		return converter(input1, input2)
	})
}

func safeExecute[T any](ctx context.Context, action func() T) T {
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
		panic(coreerror.New(
			i18n.M(ctx, i18n.K.CommonErrorInvalidFormat),
			true,
		))
	}

	return res.value
}
