package filter

import "fmt"

const (
	DataTypeStr  = "string"
	DataTypeDate = "date"

	OperatorEq      = "eq"
	OperatorNotEq   = "neq"
	OperatorBetween = "between"
	OperatorLike    = "like"
)

type options struct {
	isToApply bool
	fields    []Field
}

func NewOptions() *options {
	return &options{}
}

type Field struct {
	Name     string
	Value    string
	Operator string
	Type     string
}

type Options interface {
	IsToApply() bool
	AddField(name, operator, value, dtype string) error
	Fields() []Field
	GetField(name string) *Field
}

func (o *options) IsToApply() bool {
	return o.isToApply
}

func (o *options) AddField(name, operator, value, dtype string) error {
	err := validateOperator(operator)
	if err != nil {
		return err
	}

	o.isToApply = true
	o.fields = append(o.fields, Field{
		Name:     name,
		Value:    value,
		Operator: operator,
		Type:     dtype,
	})
	return nil
}

func (o *options) Fields() []Field {
	return o.fields
}

func (o *options) GetField(name string) *Field {
	for _, field := range o.Fields() {
		if field.Name == name {
			return &field
		}
	}
	return nil
}

func validateOperator(operator string) error {
	switch operator {
	case OperatorEq:
	case OperatorNotEq:
	case OperatorBetween:
	case OperatorLike:
	default:
		return fmt.Errorf("bad operator")
	}
	return nil
}
