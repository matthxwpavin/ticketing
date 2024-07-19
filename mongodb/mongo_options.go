package mongodb

type Validator struct {
	Schema Schema `structs:"$jsonSchema"`
}

type Schema struct {
	Required   []string            `structs:"required"`
	Properties map[string]Property `structs:"properties"`
}

const (
	BSONTypeString = "string"
	BSONTypeDouble = "double"
)

type Property struct {
	BSONType    string `structs:"bsonType"`
	Description string `structs:"description"`
	MinLength   uint   `structs:"minLength"`
}

type NamedProperty struct {
	Name     string
	Property Property
}

func NewProperties(prop NamedProperty, props ...NamedProperty) map[string]Property {
	ret := make(map[string]Property)
	for _, prop := range append([]NamedProperty{prop}, props...) {
		ret[prop.Name] = prop.Property
	}
	return ret
}
