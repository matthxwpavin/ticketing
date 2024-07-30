package mongo

type MigrationOptions struct {
	CollectionName string
	Validator      *Validator
}

type Validator struct {
	Schema *Schema
}

func (v *Validator) MongoSchema() map[string]any {
	return map[string]any{
		"$jsonSchema": v.Schema.MongoSchema(),
	}
}

type Schema struct {
	// Required []string `structs:"required"`
	// Properties map[string]Property `structs:"properties"`
	Properties []*NamedProperty
}

func (s *Schema) MongoSchema() map[string]any {
	schema := make(map[string]any)
	var required []string
	for _, prop := range s.Properties {
		if prop.IsRequired {
			required = append(required, prop.Name)
		}
		p := prop.Property
		// "bsonType":    p.BSONType,
		// 	"description": p.Description,
		// 	"minLength":   p.MinLength,
		propsMap := make(map[string]any)
		if p.BSONType != nil {
			propsMap["bsonType"] = *p.BSONType
		}
		if p.Description != nil {
			propsMap["description"] = *p.Description
		}
		if p.MinLength != nil {
			propsMap["minLength"] = *p.MinLength
		}
		schema[prop.Name] = propsMap
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

const (
	BSONTypeString = "string"
	BSONTypeDouble = "double"
)

type Property struct {
	BSONType    *string
	Description *string
	MinLength   *uint
}

type NamedProperty struct {
	Name       string
	Property   *Property
	IsRequired bool
}
