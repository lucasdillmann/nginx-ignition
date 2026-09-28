package provider

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode"
)

type (
	nginxDirective string
	nginxRawConfig string
)

func encodeNginxArguments(args []any) []any {
	encodedArgs := make([]any, len(args))

	for index, arg := range args {
		encodedArgs[index] = encodeNginxArgument(arg)
	}

	return encodedArgs
}

func nginxSprintf(format string, args ...any) string {
	return fmt.Sprintf(format, encodeNginxArguments(args)...)
}

func nginxFprintf(writer io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(writer, format, encodeNginxArguments(args)...)
}

func nginxSprintfArgs[T any](format string, args ...T) string {
	builder := strings.Builder{}
	nginxFprintfArgs(&builder, format, args...)

	return builder.String()
}

func nginxFprintfArgs[T any](writer io.Writer, format string, args ...T) {
	if len(args) == 0 {
		nginxFprintf(writer, format)
		return
	}

	forwardArgs := make([]any, len(args))
	for index, resolver := range args {
		forwardArgs[index] = resolver
	}

	placeholders := strings.TrimSpace(strings.Repeat("%s ", len(args)))
	format = strings.Replace(format, "%s", placeholders, 1)

	nginxFprintf(writer, format, forwardArgs...)
}

func encodeNginxArgument(arg any) any {
	switch value := arg.(type) {
	case nginxDirective:
		return string(value)
	case nginxRawConfig:
		return string(value)
	case string:
		return quoteNginxArgument(value)
	case bool,
		int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		uintptr,
		float32,
		float64:
		return value
	case nil:
		return value
	}

	reflected := reflect.ValueOf(arg)
	if reflected.IsValid() && reflected.Kind() == reflect.String {
		return quoteNginxArgument(reflected.String())
	}

	return quoteNginxArgument(fmt.Sprint(arg))
}

func quoteNginxArgument(value string) string {
	builder := strings.Builder{}
	builder.Grow(len(value) + 2)
	_ = builder.WriteByte('"')

	for _, character := range value {
		switch character {
		case '\x00':
			_ = builder.WriteByte(' ')
		case '\t':
			_, _ = builder.WriteString(`\t`)
		case '\n':
			_, _ = builder.WriteString(`\n`)
		case '\r':
			_, _ = builder.WriteString(`\r`)
		case '\\':
			_, _ = builder.WriteString(`\\`)
		case '"':
			_, _ = builder.WriteString(`\"`)
		default:
			if unicode.IsControl(character) {
				_ = builder.WriteByte(' ')
				continue
			}

			_, _ = builder.WriteRune(character)
		}
	}

	_ = builder.WriteByte('"')
	return builder.String()
}

func directiveFragment(contents string) nginxDirective {
	return nginxDirective(contents)
}

func rawConfigFragment(contents string) nginxRawConfig {
	return nginxRawConfig(contents)
}

func sanitizeHeaderValue(value string) string {
	return strings.
		NewReplacer(
			"\r\n", " ",
			"\r", " ",
			"\n", " ",
			"\x00", " ",
		).
		Replace(value)
}

func sanitizeHtpasswdUsername(value string) string {
	return strings.Map(
		func(character rune) rune {
			if character == ':' || unicode.IsControl(character) {
				return '_'
			}

			return character
		},
		value,
	)
}
