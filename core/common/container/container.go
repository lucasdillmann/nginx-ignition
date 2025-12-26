package container

import (
	"context"

	"go.uber.org/dig"
)

var delegate *dig.Container

func Init(ctx context.Context) {
	delegate = dig.New()

	_ = Singleton(ctx)
}

func Provide(providers ...any) error {
	for _, provider := range providers {
		if err := delegate.Provide(provider); err != nil {
			return err
		}
	}

	return nil
}

func Singleton[T any](value T) error {
	return Provide(func() T {
		return value
	})
}

func Run(fns ...any) error {
	for _, fn := range fns {
		if err := delegate.Invoke(fn); err != nil {
			return err
		}
	}

	return nil
}

func Get[T any]() T {
	var output T
	_ = delegate.Invoke(func(value T) {
		output = value
	})

	return output
}
