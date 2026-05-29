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

var errParse = errors.New("parse error")

var errRange = errors.New("value out of range")

var ErrHelp = errors.New("flag: help requested")

func numError(err error) error {
	numError, ok := err.(*strconv.NumError)

	if !ok {
		return err
	}

	if numError.Err == strconv.ErrSyntax {
		return errParse
	}

	if numError.Err == strconv.ErrRange {
		return errRange
	}

	return err
}

type boolValue bool

func newBoolValue(value bool, pointer *bool) *boolValue {
	*pointer = value

	return (*boolValue) (pointer)
}

func (self *boolValue) Set(str string) error {
	value, err := strconv.ParseBool(str)

	if err != nil {
		err = errParse
	}

	*self = boolValue(value)

	return err
}

func (self *boolValue) Get() any {
	return bool(*self)
}

func (self *boolValue) String() string {
	return strconv.FormatBool(bool(*self))
}

func (self *boolValue) IsBoolFlag() bool {
	return true
}

type boolFlag interface {
	Value
	IsBoolFlag() bool
}

type intValue int

func newIntValue(value int, pointer *int) *intValue {
	*pointer = value

	return (*intValue) (pointer)
}

func (self *intValue) Set(str string) error {
	value, err := strconv.ParseInt(str, 0, strconv.IntSize)

	if err != nil {
		err = numError(err)
	}

	*self = intValue(value)

	return err
}

func (self *intValue) Get() any {
	return int(*self)
}

func (self *intValue) String() string {
	return strconv.Itoa(int(*self))
}

type int64Value int64

func newInt64Value(value int64, pointer *int64) *int64Value {
	*pointer = value

	return (*int64Value) (pointer)
}

func (self *int64Value) Set(str string) error {
	value, err := strconv.ParseInt(str, 0, 64)

	if err != nil {
		err = numError(err)
	}

	*self = int64Value(value)

	return err
}

func (self *int64Value) Get() any {
	return int64(*self)
}

func (self *int64Value) String() string {
	return strconv.FormatInt(int64(*self), 10)
}

type uintValue uint

func newUintValue(value uint, pointer *uint) *uintValue {
	*pointer = value

	return (*uintValue) (pointer)
}

func (self *uintValue) Set(str string) error {
	value, err := strconv.ParseUint(str, 0, strconv.IntSize)

	if err != nil {
		err = numError(err)
	}

	*self = uintValue(value)

	return err
}

func (self *uintValue) Get() any {
	return uint(*self)
}

func (self *uintValue) String() string {
	return strconv.FormatUint(uint64(*self), 10)
}

type uint64Value uint64

func newUint64Value(value uint64, pointer *uint64) *uint64Value {
	*pointer = value

	return (*uint64Value) (pointer)
}

func (self *uint64Value) Set(str string) error {
	value, err := strconv.ParseUint(str, 0, 64)

	if err != nil {
		err = numError(err)
	}

	*self = uint64Value(value)

	return err
}

func (self *uint64Value) Get() any {
	return uint64(*self)
}

func (self *uint64Value) String() string {
	return strconv.FormatUint(uint64(*self), 10)
}

type stringValue string

func newStringValue(value string, pointer *string) *stringValue {
	*pointer = value

	return (*stringValue) (pointer)
}

func (self *stringValue) Set(value string) error {
	*self = stringValue(value)

	return nil
}

func (self *stringValue) Get() any {
	return string(*self)
}

func (self *stringValue) String() string {
	return string(*self)
}

type float64Value float64

func newFloat64Value(value float64, pointer *float64) *float64Value {
	*pointer = value

	return (*float64Value) (pointer)
}

func (self *float64Value) Set(str string) error {
	value, err := strconv.ParseFloat(str, 64)

	if err != nil {
		err = numError(err)
	}

	*self = float64Value(value)

	return err
}

func (self *float64Value) Get() any {
	return float64(*self)
}

func (self *float64Value) String() string {
	return strconv.FormatFloat(float64(*self), 'g', -1, 64)
}

type durationValue time.Duration

func newDurationValue(value time.Duration, pointer *time.Duration) *durationValue {
	*pointer = value

	return (*durationValue) (pointer)
}

func (self *durationValue) Set(str string) error {
	value, err := time.ParseDuration(str)

	if err != nil {
		err = errParse
	}

	*self = durationValue(value)

	return err
}

func (self *durationValue) Get() any {
	return time.Duration(*self)
}

func (self *durationValue) String() string {
	return (*time.Duration) (self).String()
}

type textValue struct {
	pointer encoding.TextUnmarshaler
}

func newTextValue(value encoding.TextMarshaler, pointer encoding.TextUnmarshaler) textValue {
	pointerValue := reflect.ValueOf(pointer)

	if pointerValue.Kind() != reflect.Ptr {
		panic("variable value type must be a pointer")
	}

	defineValue := reflect.ValueOf(value)

	if defineValue.Kind() == reflect.Ptr {
		defineValue = defineValue.Elem()
	}

	if defineValue.Type() != pointerValue.Type().Elem() {
		panic(fmt.Sprintf("default type does not match variable type: %v != %v", defineValue.Type(), pointerValue.Type().Elem()))
	}

	pointerValue.Elem().Set(defineValue)

	return textValue {pointer}
}

func (self textValue) Set(str string) error {
	return self.pointer.UnmarshalText([]byte(str))
}

func (self textValue) Get() any {
	return self.pointer
}

func (self textValue) String() string {
	if marshaler, ok := self.pointer.(encoding.TextMarshaler); ok {
		if encodedText, err := marshaler.MarshalText(); err == nil {
			return string(encodedText)
		}
	}

	return ""
}

type funcValue func(string) error

func (self funcValue) Set(str string) error {
	return self(str)
}

func (self funcValue) String() string {
	return ""
}

type boolFuncValue func(string) error

func (self boolFuncValue) Set(str string) error {
	return self(str)
}

func (self boolFuncValue) String() string {
	return ""
}

func (self boolFuncValue) IsBoolFlag() bool {
	return true
}

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
	undefined map[string]string
}

type Flag struct {
	Name string
	Usage string
	Value Value
	DefineValue string
}

func sortFlags(flags map[string]*Flag) []*Flag {
	result := make([]*Flag, len(flags))

	index := 0

	for _, flag := range flags {
		result[index] = flag

		index++
	}

	slices.SortFunc(result, func(a, b *Flag) int { // TODO: descobrir o que "a" e "b" significa
		return strings.Compare(a.Name, b.Name)
	})

	return result
}

func (self *FlagSet) Output() io.Writer {
	if self.output == nil {
		return os.Stderr
	}

	return self.output
}

func (self *FlagSet) Name() string {
	return self.name
}

func (self *FlagSet) ErrorHandling() ErrorHandling {
	return self.errorHandling
}

func (self *FlagSet) SetOutput(output io.Writer) {
	self.output = output
}

func (self *FlagSet) VisitAll(function func(*Flag)) {
	for _, flag := range sortFlags(self.formal) {
		function(flag)
	}
}

func VisitAll(function func(*Flag)) {
	CommandLine.VisitAll(function)
}

func (self *FlagSet) Visit(function func(*Flag)) {
	for _, flag := range sortFlags(self.actual) {
		function(flag)
	}
}

func Visit(function func(*Flag)) {
	CommandLine.Visit(function)
}

func (self *FlagSet) Lookup(name string) *Flag {
	return self.formal[name]
}

func Lookup(name string) *Flag {
	return CommandLine.formal[name]
}

func (self *FlagSet) Set(name, value string) error {
	return self.set(name, value)
}

func (self *FlagSet) set(name, value string) error {
	flag, ok := self.formal[name]

	if !ok {
		_, file, line, ok := runtime.Caller(2)

		if !ok {
			file = "?"

			line = 0
		}

		if self.undefined == nil {
			self.undefined = map[string]string{}
		}

		self.undefined[name] = fmt.Sprintf("%s:%d", file, line)

		return fmt.Errorf("no such flag -%v", name)
	}

	err := flag.Value.Set(value)

	if err != nil {
		return err
	}

	if self.actual == nil {
		self.actual = make(map[string]*Flag)
	}

	self.actual[name] = flag

	return nil
}

func Set(name, value string) error {
	return CommandLine.set(name, value)
}

func isZeroValue(flag *Flag, value string) (ok bool, err error) {
	typeOf := reflect.TypeOf(flag.Value)

	var reflectedValue reflect.Value

	if typeOf.Kind() == reflect.Pointer {
		reflectedValue = reflect.New(typeOf.Elem())
	} else {
		reflectedValue = reflect.Zero(typeOf)
	}

	defer func() {
		if e := recover(); e != nil { // TODO: nomear corretamente a variável "e"
			if typeOf.Kind() == reflect.Pointer {
				typeOf = typeOf.Elem()
			}

			// TOOD: descobrir para que serve a variável "err"
			err = fmt.Errorf("panic calling String method on zero %v for flag %s: %v", typeOf, flag.Name, e)
		}
	}()

	return value == reflectedValue.Interface().(Value).String(), nil
}

func UnquoteUsage(flag *Flag) (name string, usage string) {
	usage = flag.Usage

	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i + 1 : j]

					usage = usage[:i] + name + usage[j + 1:]

					return name, usage
				}
			}

			break
		}
	}

	name = "value"

	switch flagValue := flag.Value.(type) {
	case boolFlag:
		if flagValue.IsBoolFlag() {
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

func (self *FlagSet) PrintDefaults() {
	var isZeroValueErrs []error

	fmt.Println("Options:")

	self.VisitAll(func(flag *Flag) {
		var builder strings.Builder

		name, usage := UnquoteUsage(flag)

		fmt.Fprintf(&builder, "  -%s %v", flag.Name, name)

		maxNumber := 14 - len(fmt.Sprintf("-%s %v", flag.Name, name))

		for number := 0; number < maxNumber; number++ {
			builder.WriteString(" ")
		}

		fmt.Fprintf(&builder, "%s", usage)

		if isZero, err := isZeroValue(flag, flag.DefineValue); err != nil {
			isZeroValueErrs = append(isZeroValueErrs, err)
		} else if !isZero {
			if _, ok := flag.Value.(*stringValue); ok {
				fmt.Fprintf(&builder, " (default is %q)", flag.DefineValue)
			} else {
				fmt.Fprintf(&builder, " (default is %v)", flag.DefineValue)
			}
		}

		fmt.Fprint(self.Output(), builder.String(), "\n")
	})

	if errs := isZeroValueErrs; len(errs) > 0 {
		fmt.Fprintln(self.Output())

		for _, err := range errs {
			fmt.Fprintln(self.Output(), err)
		}
	}
}

func PrintDefaults() {
	CommandLine.PrintDefaults()
}

func (self *FlagSet) defaultUsage() {
	if self.name == "" {
		fmt.Fprintf(self.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(self.Output(), "Usage: amt %s [options]\n\n", self.name)
	}

	self.PrintDefaults()
}

var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])

	PrintDefaults()
}

func (self *FlagSet) NFlag() int {
	return len(self.actual)
}

func NFlag() int {
	return len(CommandLine.actual)
}

func (self *FlagSet) Arg(index int) string {
	if index < 0 || index >= len(self.args) {
		return ""
	}

	return self.args[index]
}

func Arg(index int) string {
	return CommandLine.Arg(index)
}

func (self *FlagSet) NArg() int {
	return len(self.args)
}

func NArg() int {
	return len(CommandLine.args)
}

func (self *FlagSet) Args() []string {
	return self.args
}

func Args() []string {
	return CommandLine.args
}

func (self *FlagSet) BoolVar(pointer *bool, name string, value bool, usage string) {
	self.Var(newBoolValue(value, pointer), name, usage)
}

func BoolVar(pointer *bool, name string, value bool, usage string) {
	CommandLine.Var(newBoolValue(value, pointer), name, usage)
}

func (self *FlagSet) Bool(name string, value bool, usage string) *bool {
	pointer := new(bool)

	self.BoolVar(pointer, name, value, usage)

	return pointer
}

func Bool(name string, value bool, usage string) *bool {
	return CommandLine.Bool(name, value, usage)
}

func (self *FlagSet) IntVar(pointer *int, name string, value int, usage string) {
	self.Var(newIntValue(value, pointer), name, usage)
}

func IntVar(pointer *int, name string, value int, usage string) {
	CommandLine.Var(newIntValue(value, pointer), name, usage)
}

func (self *FlagSet) Int(name string, value int, usage string) *int {
	pointer := new(int)

	self.IntVar(pointer, name, value, usage)

	return pointer
}

func Int(name string, value int, usage string) *int {
	return CommandLine.Int(name, value, usage)
}

func (self *FlagSet) Int64Var(pointer *int64, name string, value int64, usage string) {
	self.Var(newInt64Value(value, pointer), name, usage)
}

func Int64Var(pointer *int64, name string, value int64, usage string) {
	CommandLine.Var(newInt64Value(value, pointer), name, usage)
}

func (self *FlagSet) Int64(name string, value int64, usage string) *int64 {
	pointer := new(int64)

	self.Int64Var(pointer, name, value, usage)

	return pointer
}

func Int64(name string, value int64, usage string) *int64 {
	return CommandLine.Int64(name, value, usage)
}

func (self *FlagSet) UintVar(pointer *uint, name string, value uint, usage string) {
	self.Var(newUintValue(value, pointer), name, usage)
}

func UintVar(pointer *uint, name string, value uint, usage string) {
	CommandLine.Var(newUintValue(value, pointer), name, usage)
}

func (self *FlagSet) Uint(name string, value uint, usage string) *uint {
	pointer := new(uint)

	self.UintVar(pointer, name, value, usage)

	return pointer
}

func Uint(name string, value uint, usage string) *uint {
	return CommandLine.Uint(name, value, usage)
}

func (self *FlagSet) Uint64Var(pointer *uint64, name string, value uint64, usage string) {
	self.Var(newUint64Value(value, pointer), name, usage)
}

func Uint64Var(pointer *uint64, name string, value uint64, usage string) {
	CommandLine.Var(newUint64Value(value, pointer), name, usage)
}

func (self *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	pointer := new(uint64)

	self.Uint64Var(pointer, name, value, usage)

	return pointer
}

func Uint64(name string, value uint64, usage string) *uint64 {
	return CommandLine.Uint64(name, value, usage)
}

func (self *FlagSet) StringVar(pointer *string, name string, value string, usage string) {
	self.Var(newStringValue(value, pointer), name, usage)
}

func StringVar(pointer *string, name string, value string, usage string) {
	CommandLine.Var(newStringValue(value, pointer), name, usage)
}

func (self *FlagSet) String(name string, value string, usage string) *string {
	pointer := new(string)

	self.StringVar(pointer, name, value, usage)

	return pointer
}

func String(name string, value string, usage string) *string {
	return CommandLine.String(name, value, usage)
}

func (self *FlagSet) Float64Var(pointer *float64, name string, value float64, usage string) {
	self.Var(newFloat64Value(value, pointer), name, usage)
}

func Float64Var(pointer *float64, name string, value float64, usage string) {
	CommandLine.Var(newFloat64Value(value, pointer), name, usage)
}

func (self *FlagSet) Float64(name string, value float64, usage string) *float64 {
	pointer := new(float64)

	self.Float64Var(pointer, name, value, usage)

	return pointer
}

func Float64(name string, value float64, usage string) *float64 {
	return CommandLine.Float64(name, value, usage)
}

func (self *FlagSet) DurationVar(pointer *time.Duration, name string, value time.Duration, usage string) {
	self.Var(newDurationValue(value, pointer), name, usage)
}

func DurationVar(pointer *time.Duration, name string, value time.Duration, usage string) {
	CommandLine.Var(newDurationValue(value, pointer), name, usage)
}

func (self *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	pointer := new(time.Duration)

	self.DurationVar(pointer, name, value, usage)

	return pointer
}

func Duration(name string, value time.Duration, usage string) *time.Duration {
	return CommandLine.Duration(name, value, usage)
}

func (self *FlagSet) TextVar(pointer encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string) {
	self.Var(newTextValue(value, pointer), name, usage)
}

func TextVar(pointer encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string) {
	CommandLine.Var(newTextValue(value, pointer), name, usage)
}

func (self *FlagSet) Func(name, usage string, function func(string) error) {
	self.Var(funcValue(function), name, usage)
}

func Func(name, usage string, function func(string) error) {
	CommandLine.Func(name, usage, function)
}

func (self *FlagSet) BoolFunc(name, usage string, function func(string) error) {
	self.Var(boolFuncValue(function), name, usage)
}

func BoolFunc(name, usage string, function func(string) error) {
	CommandLine.BoolFunc(name, usage, function)
}

func (self *FlagSet) Var(value Value, name string, usage string) {
	if strings.HasPrefix(name, "-") {
		panic(self.sprintf("flag %q begins with -", name))
	} else if strings.Contains(name, "=") {
		panic(self.sprintf("flag %q contains =", name))
	}

	flag := &Flag{
		name,
		usage,
		value,
		value.String(),
	}

	_, alreadythere := self.formal[name]

	if alreadythere {
		var message string

		if self.name == "" {
			message = self.sprintf("flag redefined: %s", name)
		} else {
			message = self.sprintf("%s flag redefined: %s", self.name, name)
		}

		panic(message)
	}

	if position := self.undefined[name]; position != "" {
		panic(fmt.Sprintf("flag %s set at %s before being defined", name, position))
	}

	if self.formal == nil {
		self.formal = make(map[string]*Flag)
	}

	self.formal[name] = flag
}

func Var(value Value, name string, usage string) {
	CommandLine.Var(value, name, usage)
}

func (self *FlagSet) sprintf(format string, any ...any) string {
	message := fmt.Sprintf(format, any...)

	fmt.Fprintln(self.Output(), message)

	return message
}

func (self *FlagSet) failf(format string, any ...any) error {
	message := self.sprintf(format, any...)

	self.usage()

	return errors.New(message)
}

func (self *FlagSet) usage() {
	if self.Usage == nil {
		self.defaultUsage()
	} else {
		self.Usage()
	}
}

func (self *FlagSet) parseOne() (bool, error) {
	if len(self.args) == 0 {
		return false, nil
	}

	str := self.args[0]

	if len(str) < 2 || str[0] != '-' {
		return false, nil
	}

	numMinuses := 1

	if str[1] == '-' {
		numMinuses++

		if len(str) == 2 {
			self.args = self.args[1:]

			return false, nil
		}
	}

	name := str[numMinuses:]

	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, self.failf("bad flag syntax: %s", str)
	}

	self.args = self.args[1:]

	hasValue := false

	value := ""

	for index := 1; index < len(name); index++ {
		if name[index] == '=' {
			value = name[index + 1:]

			hasValue = true

			name = name[0:index]

			break
		}
	}

	flag, ok := self.formal[name]

	if !ok {
		if name == "help" || name == "h" {
			self.usage()

			return false, ErrHelp
		}

		return false, self.failf("flag provided but not defined: -%s", name)
	}

	if flagValue, ok := flag.Value.(boolFlag); ok && flagValue.IsBoolFlag() {
		if hasValue {
			if err := flagValue.Set(value); err != nil {
				return false, self.failf("invalid boolean value %q for -%s: %v", value, name, err)
			}
		} else {
			if err := flagValue.Set("true"); err != nil {
				return false, self.failf("invalid boolean flag %s: %v", name, err)
			}
		}
	} else {
		if !hasValue && len(self.args) > 0 {
			hasValue = true

			value, self.args = self.args[0], self.args[1:]
		}

		if !hasValue {
			return false, self.failf("flag needs an argument: -%s", name)
		}

		if err := flag.Value.Set(value); err != nil {
			return false, self.failf("invalid value %q for flag -%s: %v", value, name, err)
		}
	}

	if self.actual == nil {
		self.actual = make(map[string]*Flag)
	}

	self.actual[name] = flag

	return true, nil
}

func (self *FlagSet) Parse(arguments []string) error {
	self.parsed = true

	self.args = arguments

	for {
		seen, err := self.parseOne()

		if seen {
			continue
		}

		if err == nil {
			break
		}

		switch self.errorHandling {
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

func (self *FlagSet) Parsed() bool {
	return self.parsed
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
	flagSet := &FlagSet{
		name: name,
		errorHandling: errorHandling,
	}

	flagSet.Usage = flagSet.defaultUsage

	return flagSet
}

func (self *FlagSet) Init(name string, errorHandling ErrorHandling) {
	self.name = name

	self.errorHandling = errorHandling
}