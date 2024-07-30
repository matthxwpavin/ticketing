package mongo

import (
	"encoding/json"
	"testing"

	"github.com/matthxwpavin/ticketing/ptr"
)

func TestValidatorOption(t *testing.T) {
	v := &Validator{
		Schema: &Schema{
			Properties: []*NamedProperty{
				{
					Name: "title",
					Property: &Property{
						BSONType:    ptr.Of(BSONTypeString),
						Description: ptr.Of("must be a string and is required"),
					},
					IsRequired: true,
				},
				{
					Name: "price",
					Property: &Property{
						BSONType:    ptr.Of(BSONTypeDouble),
						Description: ptr.Of("must be a double and is required"),
					},
					IsRequired: true,
				},
			},
		},
	}

	var expected = map[string]any{
		"$jsonSchema": map[string]any{
			"required": []string{"title", "price"},
			"title": map[string]any{
				"bsonType":    "string",
				"description": "must be a string and is required",
			},
			"price": map[string]any{
				"bsonType":    "double",
				"description": "must be a double and is required",
			},
		},
	}

	jsonScheme, err := json.Marshal(v.MongoSchema())
	if err != nil {
		t.Fatalf("could not marshal: %v", err)
	}
	target, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("could not marshal: %v", err)
	}
	if string(jsonScheme) != string(target) {
		t.Fatalf("schema is not expected")
	}
}
