package templates

import (
	"fmt"
	"strconv"
	"strings"
)

// TemplateError represents a validation error in a template.
type TemplateError struct {
	Field   string
	Message string
}

func (e *TemplateError) Error() string {
	return fmt.Sprintf("template %s: %s", e.Field, e.Message)
}

// TemplateErrors is a collection of template validation errors.
type TemplateErrors []*TemplateError

// validConfigFieldTypes defines allowed config field types.
var validConfigFieldTypes = map[string]bool{
	"text": true, "number": true, "password": true, "select": true,
}

func (errs TemplateErrors) Error() string {
	if len(errs) == 0 {
		return ""
	}
	msgs := make([]string, len(errs))
	for i, e := range errs {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

func (errs TemplateErrors) HasErrors() bool {
	return len(errs) > 0
}

// Validate checks a template for structural correctness.
func Validate(tmpl *Template) TemplateErrors {
	var errs TemplateErrors

	if tmpl == nil {
		errs = append(errs, &TemplateError{Field: "template", Message: "template is nil"})
		return errs
	}

	if tmpl.Name == "" {
		errs = append(errs, &TemplateError{Field: "name", Message: "name is required"})
	}

	if tmpl.Description == "" {
		errs = append(errs, &TemplateError{Field: "description", Message: "description is required"})
	}

	validCategories := map[string]bool{
		"database": true,
		"cache":    true,
		"storage":  true,
		"search":   true,
	}
	if !validCategories[tmpl.Category] {
		errs = append(errs, &TemplateError{
			Field:   "category",
			Message: fmt.Sprintf("invalid category %q: must be database, cache, storage, or search", tmpl.Category),
		})
	}

	if len(tmpl.Versions) == 0 {
		errs = append(errs, &TemplateError{Field: "versions", Message: "at least one version is required"})
	}

	seenPorts := map[int]bool{}
	for _, v := range tmpl.Versions {
		if v.Key == "" {
			errs = append(errs, &TemplateError{Field: "versions.key", Message: "version key is required"})
		}
		if v.Image == "" {
			errs = append(errs, &TemplateError{
				Field:   fmt.Sprintf("versions.%s.image", v.Key),
				Message: "image is required",
			})
		}
		if v.DefaultPort < 1024 || v.DefaultPort > 65535 {
			errs = append(errs, &TemplateError{
				Field:   fmt.Sprintf("versions.%s.default_port", v.Key),
				Message: fmt.Sprintf("port %d must be between 1024 and 65535", v.DefaultPort),
			})
		}
		if seenPorts[v.DefaultPort] {
			errs = append(errs, &TemplateError{
				Field:   fmt.Sprintf("versions.%s.default_port", v.Key),
				Message: fmt.Sprintf("port %d is duplicated", v.DefaultPort),
			})
		}
		seenPorts[v.DefaultPort] = true
	}

	for _, f := range tmpl.ConfigFields {
		if f.Key == "" {
			errs = append(errs, &TemplateError{Field: "config_fields.key", Message: "field key is required"})
		}
		if !validConfigFieldTypes[f.Type] {
			errs = append(errs, &TemplateError{
				Field:   fmt.Sprintf("config_fields.%s.type", f.Key),
				Message: fmt.Sprintf("invalid type %q: must be text, number, password, or select", f.Type),
			})
		}
	}

	return errs
}

// ValidateConfigField validates a user input against a config field definition.
func ValidateConfigField(field ConfigField, value string) error {
	if field.Required && value == "" {
		return fmt.Errorf("%s is required", field.Label)
	}

	if value == "" {
		return nil
	}

	switch field.Type {
	case "number":
		n, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("%s must be a number", field.Label)
		}
		if n < 1024 || n > 65535 {
			return fmt.Errorf("%s must be between 1024 and 65535", field.Label)
		}
	case "text", "password":
		if len(value) > 256 {
			return fmt.Errorf("%s is too long (max 256 characters)", field.Label)
		}
	case "select":
		if len(field.Options) > 0 {
			valid := false
			for _, opt := range field.Options {
				if opt == value {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("%s must be one of: %s", field.Label, strings.Join(field.Options, ", "))
			}
		}
	}

	return nil
}
