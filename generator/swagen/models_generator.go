package swagen

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gopher-fleece/gleece/definitions"
)

var objectType = &openapi3.Types{"object"}
var arrayType = &openapi3.Types{"array"}

func generateModelSpec(openapi *openapi3.T, model definitions.ModelMetadata) {
	schema := &openapi3.Schema{
		Title:       model.Name,
		Description: model.Description,
		Type:        objectType,
		Properties:  openapi3.Schemas{},
		Deprecated:  IsDeprecated(&model.Deprecation),
	}

	requiredFields := []string{}

	for _, field := range model.Fields {
		fieldSchemaRef := InterfaceToSchemaRef(openapi, field.Type)

		fieldSchemaRef.Value.Description = field.Description

		BuildSchemaValidation(fieldSchemaRef, field.Validator, field.Type)

		// If the schema marked as deprecated, the field / property should be marked as deprecated as well
		// Setting it as not deprecated (even if the field itself is not marked deprecated) will override the model deprecation
		if !fieldSchemaRef.Value.Deprecated {
			fieldSchemaRef.Value.Deprecated = IsDeprecated(field.Deprecation)
		}

		// Add field to schema properties
		schema.Properties[field.Name] = fieldSchemaRef

		// If the field should be required, add its name to the requiredFields slice
		if IsFieldRequired(field.Validator) {
			requiredFields = append(requiredFields, field.Name)
		}
	}

	// Add required fields to the schema
	if len(requiredFields) > 0 {
		schema.Required = requiredFields
	}

	// Add schema to components
	openapi.Components.Schemas[model.Name] = &openapi3.SchemaRef{
		Value: schema,
	}
}

// Fill the schema references in the components
func fillSchemaRef(openapi *openapi3.T) {
	// Iterate over all EmptyRefSchemas
	for _, schemaRefMap := range schemaRefMap {
		// Get the schema from the components
		schema := openapi.Components.Schemas[schemaRefMap.InterfaceType].Value
		// Set the schema reference to the schema
		schemaRefMap.SchemaRef.Value = schema
	}
}

func GenerateModelsSpec(openapi *openapi3.T, models []definitions.ModelMetadata) error {
	for _, model := range models {
		generateModelSpec(openapi, model)
	}
	fillSchemaRef(openapi)
	return nil
}
