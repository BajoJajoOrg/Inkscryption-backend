package filter

const (
	DataTypeDate = "date"

	OperatorBetween = "be"
	OperatorEq      = "eq"
)

type options struct {
	isToApply bool
	fields    []Field
}

func NewOptions(isToApply bool, fields []Field) *options {
	return &options{isToApply: isToApply, fields: fields}
}

type Field struct {
	Name     string
	Value    string
	Operator string
	Type     string
}

type Options interface {
	IsToApply() bool
	AddField(name, operator, value, dtype string)
	Fields() []Field
}

func (o *options) IsToApply() bool {
	return o.isToApply
}

func (o *options) AddField(name, operator, value, dtype string) {
	o.fields = append(o.fields, Field{
		Name:     name,
		Value:    value,
		Operator: operator,
		Type:     dtype,
	})
}

func (o *options) Fields() []Field {
	return o.fields
}
