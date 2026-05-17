// Copied from https://github.com/golang/go/blob/master/src/flag/flag.go
package flag

import (
	"io"
	"os"
	"fmt"
	"time"
	"errors"
	"slices"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"encoding"
)

var ErrHelp = errors.New("flag: help requested")

var errParse = errors.New("parse error")

var errRange = errors.New("value out of range")

func numError(err error) error {
	ne, ok := err.(*strconv.NumError)

	if !ok {
		return err
	}

	if ne.Err == strconv.ErrSyntax {
		return errParse
	}

	if ne.Err == strconv.ErrRange {
		return errRange
	}

	return err
}

type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val

	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)

	if err != nil {
		err = errParse
	}

	*b = boolValue(v)

	return err
}

func (b *boolValue) Get() any { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }

type boolFlag interface {
	Value
	IsBoolFlag() bool
}

type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val

	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)

	if err != nil {
		err = numError(err)
	}

	*i = intValue(v)

	return err
}

func (i *intValue) Get() any { return int(*i) }

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val

	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)

	if err != nil {
		err = numError(err)
	}

	*i = int64Value(v)

	return err
}

func (i *int64Value) Get() any { return int64(*i) }

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val

	return (*uintValue)(p)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)

	if err != nil {
		err = numError(err)
	}

	*i = uintValue(v)

	return err
}

func (i *uintValue) Get() any { return uint(*i) }

func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }

type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val

	return (*uint64Value)(p)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)

	if err != nil {
		err = numError(err)
	}

	*i = uint64Value(v)

	return err
}

func (i *uint64Value) Get() any { return uint64(*i) }

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val

	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)

	return nil
}

func (s *stringValue) Get() any { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val

	return (*float64Value)(p)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)

	if err != nil {
		err = numError(err)
	}

	*f = float64Value(v)

	return err
}

func (f *float64Value) Get() any { return float64(*f) }

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val

	return (*durationValue)(p)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)

	if err != nil {
		err = errParse
	}

	*d = durationValue(v)

	return err
}

func (d *durationValue) Get() any { return time.Duration(*d) }

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

type textValue struct{ p encoding.TextUnmarshaler }

func newTextValue(val encoding.TextMarshaler, p encoding.TextUnmarshaler) textValue {
	ptrVal := reflect.ValueOf(p)

	if ptrVal.Kind() != reflect.Ptr {
		panic("variable value type must be a pointer")
	}

	defVal := reflect.ValueOf(val)

	if defVal.Kind() == reflect.Ptr {
		defVal = defVal.Elem()
	}

	if defVal.Type() != ptrVal.Type().Elem() {
		panic(fmt.Sprintf("default type does not match variable type: %v != %v", defVal.Type(), ptrVal.Type().Elem()))
	}

	ptrVal.Elem().Set(defVal)

	return textValue{p}
}

func (v textValue) Set(s string) error {
	return v.p.UnmarshalText([]byte(s))
}

func (v textValue) Get() any {
	return v.p
}

func (v textValue) String() string {
	if m, ok := v.p.(encoding.TextMarshaler); ok {
		if b, err := m.MarshalText(); err == nil {
			return string(b)
		}
	}

	return ""
}

type funcValue func(string) error

func (f funcValue) Set(s string) error { return f(s) }

func (f funcValue) String() string { return "" }

type boolFuncValue func(string) error

func (f boolFuncValue) Set(s string) error { return f(s) }

func (f boolFuncValue) String() string { return "" }

func (f boolFuncValue) IsBoolFlag() bool { return true }

type Value interface {
	String() string
	Set(string) error
}
type Getter interface {
	Value
	Get() any
}

type ErrorHandling int

const (
	ContinueOnError ErrorHandling = iota
	ExitOnError
	PanicOnError
)

type FlagSet struct {
	Usage func()

	name string
	parsed bool
	actual map[string]*Flag
	formal map[string]*Flag
	args []string
	errorHandling ErrorHandling
	output io.Writer
	undef map[string]string
}

type Flag struct {
	Name string
	Usage string
	Value Value
	DefValue string
}

func sortFlags(flags map[string]*Flag) []*Flag {
	result := make([]*Flag, len(flags))

	i := 0

	for _, f := range flags {
		result[i] = f

		i++
	}

	slices.SortFunc(result, func(a, b *Flag) int {
		return strings.Compare(a.Name, b.Name)
	})

	return result
}

func (f *FlagSet) Output() io.Writer {
	if f.output == nil {
		return os.Stderr
	}

	return f.output
}

func (f *FlagSet) Name() string {
	return f.name
}

func (f *FlagSet) ErrorHandling() ErrorHandling {
	return f.errorHandling
}

func (f *FlagSet) SetOutput(output io.Writer) {
	f.output = output
}

func (f *FlagSet) VisitAll(fn func(*Flag)) {
	for _, flag := range sortFlags(f.formal) {
		fn(flag)
	}
}

func VisitAll(fn func(*Flag)) {
	CommandLine.VisitAll(fn)
}

func (f *FlagSet) Visit(fn func(*Flag)) {
	for _, flag := range sortFlags(f.actual) {
		fn(flag)
	}
}

func Visit(fn func(*Flag)) {
	CommandLine.Visit(fn)
}

func (f *FlagSet) Lookup(name string) *Flag {
	return f.formal[name]
}

func Lookup(name string) *Flag {
	return CommandLine.formal[name]
}

func (f *FlagSet) Set(name, value string) error {
	return f.set(name, value)
}

func (f *FlagSet) set(name, value string) error {
	flag, ok := f.formal[name]

	if !ok {
		_, file, line, ok := runtime.Caller(2)

		if !ok {
			file = "?"

			line = 0
		}

		if f.undef == nil {
			f.undef = map[string]string{}
		}

		f.undef[name] = fmt.Sprintf("%s:%d", file, line)

		return fmt.Errorf("no such flag -%v", name)
	}

	err := flag.Value.Set(value)

	if err != nil {
		return err
	}

	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}

	f.actual[name] = flag

	return nil
}

func Set(name, value string) error {
	return CommandLine.set(name, value)
}

func isZeroValue(flag *Flag, value string) (ok bool, err error) {
	typ := reflect.TypeOf(flag.Value)

	var z reflect.Value

	if typ.Kind() == reflect.Pointer {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}

	defer func() {
		if e := recover(); e != nil {
			if typ.Kind() == reflect.Pointer {
				typ = typ.Elem()
			}

			err = fmt.Errorf("panic calling String method on zero %v for flag %s: %v", typ, flag.Name, e)
		}
	}()

	return value == z.Interface().(Value).String(), nil
}

func UnquoteUsage(flag *Flag) (name string, usage string) {
	usage = flag.Usage

	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]

					usage = usage[:i] + name + usage[j+1:]

					return name, usage
				}
			}

			break
		}
	}

	name = "value"

	switch fv := flag.Value.(type) {
	case boolFlag:
		if fv.IsBoolFlag() {
			name = ""
		}
	case *durationValue:
		name = "duration"
	case *float64Value:
		name = "float"
	case *intValue, *int64Value:
		name = "int"
	case *stringValue:
		name = "string"
	case *uintValue, *uint64Value:
		name = "uint"
	}

	return
}

func (f *FlagSet) PrintDefaults() {
	var isZeroValueErrs []error

	fmt.Println("\nOptions:")

	f.VisitAll(func(flag *Flag) {
		var b strings.Builder

		name, usage := UnquoteUsage(flag)

		fmt.Fprintf(&b, "  -%s %v", flag.Name, name)

		maxNumber := 15 - len(fmt.Sprintf("-%s %v", flag.Name, name))

		for number := 0; number < maxNumber; number++ {
			b.WriteString(" ")
		}

		fmt.Fprintf(&b, "%s", usage)

		if isZero, err := isZeroValue(flag, flag.DefValue); err != nil {
			isZeroValueErrs = append(isZeroValueErrs, err)
		} else if !isZero {
			if _, ok := flag.Value.(*stringValue); ok {
				fmt.Fprintf(&b, " (default is %q)", flag.DefValue)
			} else {
				fmt.Fprintf(&b, " (default is %v)", flag.DefValue)
			}
		}

		fmt.Fprint(f.Output(), b.String(), "\n")
	})

	if errs := isZeroValueErrs; len(errs) > 0 {
		fmt.Fprintln(f.Output())

		for _, err := range errs {
			fmt.Fprintln(f.Output(), err)
		}
	}
}

func PrintDefaults() {
	CommandLine.PrintDefaults()
}

func (f *FlagSet) defaultUsage() {
	if f.name == "" {
		fmt.Fprintf(f.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(f.Output(), "Usage: amt %s [options]\n", f.name)
	}

	f.PrintDefaults()
}

var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])

	PrintDefaults()
}

func (f *FlagSet) NFlag() int { return len(f.actual) }

func NFlag() int { return len(CommandLine.actual) }

func (f *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}

	return f.args[i]
}

func Arg(i int) string {
	return CommandLine.Arg(i)
}

func (f *FlagSet) NArg() int { return len(f.args) }

func NArg() int { return len(CommandLine.args) }

func (f *FlagSet) Args() []string { return f.args }

func Args() []string { return CommandLine.args }

func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	f.Var(newBoolValue(value, p), name, usage)
}

func BoolVar(p *bool, name string, value bool, usage string) {
	CommandLine.Var(newBoolValue(value, p), name, usage)
}

func (f *FlagSet) Bool(name string, value bool, usage string) *bool {
	p := new(bool)

	f.BoolVar(p, name, value, usage)

	return p
}

func Bool(name string, value bool, usage string) *bool {
	return CommandLine.Bool(name, value, usage)
}

func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.Var(newIntValue(value, p), name, usage)
}

func IntVar(p *int, name string, value int, usage string) {
	CommandLine.Var(newIntValue(value, p), name, usage)
}

func (f *FlagSet) Int(name string, value int, usage string) *int {
	p := new(int)

	f.IntVar(p, name, value, usage)

	return p
}

func Int(name string, value int, usage string) *int {
	return CommandLine.Int(name, value, usage)
}

func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	f.Var(newInt64Value(value, p), name, usage)
}

func Int64Var(p *int64, name string, value int64, usage string) {
	CommandLine.Var(newInt64Value(value, p), name, usage)
}

func (f *FlagSet) Int64(name string, value int64, usage string) *int64 {
	p := new(int64)

	f.Int64Var(p, name, value, usage)

	return p
}

func Int64(name string, value int64, usage string) *int64 {
	return CommandLine.Int64(name, value, usage)
}

func (f *FlagSet) UintVar(p *uint, name string, value uint, usage string) {
	f.Var(newUintValue(value, p), name, usage)
}

func UintVar(p *uint, name string, value uint, usage string) {
	CommandLine.Var(newUintValue(value, p), name, usage)
}

func (f *FlagSet) Uint(name string, value uint, usage string) *uint {
	p := new(uint)

	f.UintVar(p, name, value, usage)

	return p
}

func Uint(name string, value uint, usage string) *uint {
	return CommandLine.Uint(name, value, usage)
}

func (f *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) {
	f.Var(newUint64Value(value, p), name, usage)
}

func Uint64Var(p *uint64, name string, value uint64, usage string) {
	CommandLine.Var(newUint64Value(value, p), name, usage)
}

func (f *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	p := new(uint64)

	f.Uint64Var(p, name, value, usage)

	return p
}

func Uint64(name string, value uint64, usage string) *uint64 {
	return CommandLine.Uint64(name, value, usage)
}

func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
	f.Var(newStringValue(value, p), name, usage)
}

func StringVar(p *string, name string, value string, usage string) {
	CommandLine.Var(newStringValue(value, p), name, usage)
}

func (f *FlagSet) String(name string, value string, usage string) *string {
	p := new(string)

	f.StringVar(p, name, value, usage)

	return p
}

func String(name string, value string, usage string) *string {
	return CommandLine.String(name, value, usage)
}

func (f *FlagSet) Float64Var(p *float64, name string, value float64, usage string) {
	f.Var(newFloat64Value(value, p), name, usage)
}

func Float64Var(p *float64, name string, value float64, usage string) {
	CommandLine.Var(newFloat64Value(value, p), name, usage)
}

func (f *FlagSet) Float64(name string, value float64, usage string) *float64 {
	p := new(float64)

	f.Float64Var(p, name, value, usage)

	return p
}

func Float64(name string, value float64, usage string) *float64 {
	return CommandLine.Float64(name, value, usage)
}

func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	f.Var(newDurationValue(value, p), name, usage)
}

func DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	CommandLine.Var(newDurationValue(value, p), name, usage)
}

func (f *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)

	f.DurationVar(p, name, value, usage)

	return p
}

func Duration(name string, value time.Duration, usage string) *time.Duration {
	return CommandLine.Duration(name, value, usage)
}

func (f *FlagSet) TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string) {
	f.Var(newTextValue(value, p), name, usage)
}

func TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string) {
	CommandLine.Var(newTextValue(value, p), name, usage)
}

func (f *FlagSet) Func(name, usage string, fn func(string) error) {
	f.Var(funcValue(fn), name, usage)
}

func Func(name, usage string, fn func(string) error) {
	CommandLine.Func(name, usage, fn)
}

func (f *FlagSet) BoolFunc(name, usage string, fn func(string) error) {
	f.Var(boolFuncValue(fn), name, usage)
}

func BoolFunc(name, usage string, fn func(string) error) {
	CommandLine.BoolFunc(name, usage, fn)
}

func (f *FlagSet) Var(value Value, name string, usage string) {
	if strings.HasPrefix(name, "-") {
		panic(f.sprintf("flag %q begins with -", name))
	} else if strings.Contains(name, "=") {
		panic(f.sprintf("flag %q contains =", name))
	}

	flag := &Flag{name, usage, value, value.String()}

	_, alreadythere := f.formal[name]

	if alreadythere {
		var msg string

		if f.name == "" {
			msg = f.sprintf("flag redefined: %s", name)
		} else {
			msg = f.sprintf("%s flag redefined: %s", f.name, name)
		}

		panic(msg)
	}

	if pos := f.undef[name]; pos != "" {
		panic(fmt.Sprintf("flag %s set at %s before being defined", name, pos))
	}

	if f.formal == nil {
		f.formal = make(map[string]*Flag)
	}

	f.formal[name] = flag
}

func Var(value Value, name string, usage string) {
	CommandLine.Var(value, name, usage)
}

func (f *FlagSet) sprintf(format string, a ...any) string {
	msg := fmt.Sprintf(format, a...)

	fmt.Fprintln(f.Output(), msg)

	return msg
}

func (f *FlagSet) failf(format string, a ...any) error {
	msg := f.sprintf(format, a...)

	f.usage()

	return errors.New(msg)
}

func (f *FlagSet) usage() {
	if f.Usage == nil {
		f.defaultUsage()
	} else {
		f.Usage()
	}
}

func (f *FlagSet) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}

	s := f.args[0]

	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}

	numMinuses := 1

	if s[1] == '-' {
		numMinuses++

		if len(s) == 2 {
			f.args = f.args[1:]

			return false, nil
		}
	}

	name := s[numMinuses:]

	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, f.failf("bad flag syntax: %s", s)
	}

	f.args = f.args[1:]

	hasValue := false

	value := ""

	for i := 1; i < len(name); i++ {
		if name[i] == '=' {
			value = name[i+1:]

			hasValue = true

			name = name[0:i]

			break
		}
	}

	flag, ok := f.formal[name]

	if !ok {
		if name == "help" || name == "h" {
			f.usage()

			return false, ErrHelp
		}

		return false, f.failf("flag provided but not defined: -%s", name)
	}

	if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() {
		if hasValue {
			if err := fv.Set(value); err != nil {
				return false, f.failf("invalid boolean value %q for -%s: %v", value, name, err)
			}
		} else {
			if err := fv.Set("true"); err != nil {
				return false, f.failf("invalid boolean flag %s: %v", name, err)
			}
		}
	} else {
		if !hasValue && len(f.args) > 0 {
			hasValue = true

			value, f.args = f.args[0], f.args[1:]
		}

		if !hasValue {
			return false, f.failf("flag needs an argument: -%s", name)
		}

		if err := flag.Value.Set(value); err != nil {
			return false, f.failf("invalid value %q for flag -%s: %v", value, name, err)
		}
	}

	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}

	f.actual[name] = flag

	return true, nil
}

func (f *FlagSet) Parse(arguments []string) error {
	f.parsed = true

	f.args = arguments

	for {
		seen, err := f.parseOne()

		if seen {
			continue
		}

		if err == nil {
			break
		}

		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			if err == ErrHelp {
				os.Exit(0)
			}

			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}

	return nil
}

func (f *FlagSet) Parsed() bool {
	return f.parsed
}

func Parse() {
	CommandLine.Parse(os.Args[1:])
}

func Parsed() bool {
	return CommandLine.Parsed()
}

var CommandLine *FlagSet

func init() {
	if len(os.Args) == 0 {
		CommandLine = NewFlagSet("", ExitOnError)
	} else {
		CommandLine = NewFlagSet(os.Args[0], ExitOnError)
	}

	CommandLine.Usage = commandLineUsage
}

func commandLineUsage() {
	Usage()
}

func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name: name,
		errorHandling: errorHandling,
	}

	f.Usage = f.defaultUsage

	return f
}

func (f *FlagSet) Init(name string, errorHandling ErrorHandling) {
	f.name = name

	f.errorHandling = errorHandling
}