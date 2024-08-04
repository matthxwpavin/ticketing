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
	Properties []*NamedProperty
}

func (s *Schema) MongoSchema() map[string]any {
	schema := make(map[string]any)
	var required []string
	type properties map[string]map[string]any
	if len(s.Properties) > 0 {
		schema["properties"] = make(properties)
	}
	for _, prop := range s.Properties {
		if prop.IsRequired {
			required = append(required, prop.Name)
		}
		p := prop.Property
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
		if len(p.Enum) > 0 {
			propsMap["enum"] = p.Enum
		}
		schema["properties"].(properties)[prop.Name] = propsMap
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

const (
	BSONTypeString = "string"
	BSONTypeDouble = "double"
	BSONTypeDate   = "date"
	BSONTypeInt    = "int"
)

type Property struct {
	BSONType    *string
	Description *string
	MinLength   *uint
	Enum        []string
}

type NamedProperty struct {
	Name       string
	Property   *Property
	IsRequired bool
}
