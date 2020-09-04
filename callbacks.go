package validations

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/moisespsena-go/aorm"
)

// NewError generate a new error for a model's field
func Failed(resource interface{}, column, err string, path ...string) error {
	if len(path) == 0 {
		path = []string{""}
	}
	return &ValidationFailed{Resource: resource, Column: column, Message: err, Path: path[0]}
}

// NewError generate a new error for a model's field
func FieldFailed(resource interface{}, fieldName, err string, path ...string) error {
	if len(path) == 0 {
		path = []string{""}
	}
	return &ValidationFailed{Resource: resource, FieldName: fieldName, Message: err, Path: path[0]}
}

// Error is a validation error struct that hold model, column and error message
type ValidationFailed struct {
	Resource  interface{}
	Column    string
	FieldName string
	Message   string
	Path      string
}

// Label is a label including model type, primary key and column name
func (err ValidationFailed) Label() string {
	struc := aorm.StructOf(err.Resource)
	return fmt.Sprintf("%v_%s_%v", struc.Type.Name(), struc.GetID(err.Resource), err.Column)
}

// Error show error message
func (err ValidationFailed) Error() string {
	return err.Message
}

var skipValidations = "validations:skip_validations"

func validate(scope *aorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		if result, ok := scope.DB().Get(skipValidations); !(ok && result.(bool)) {
			if !scope.HasError() {
				scope.CallMethod("Validate")
				if scope.Value != nil {
					resource := scope.IndirectValue().Interface()
					_, validatorErrors := govalidator.ValidateStruct(resource)
					if validatorErrors != nil {
						if errors, ok := validatorErrors.(govalidator.Errors); ok {
							for _, err := range flatValidatorErrors(errors) {
								scope.DB().AddError(formattedError(err, resource))
							}
						} else {
							scope.DB().AddError(validatorErrors)
						}
					}
				}
			}
		}
	}
}

func flatValidatorErrors(validatorErrors govalidator.Errors) []govalidator.Error {
	resultErrors := []govalidator.Error{}
	for _, validatorError := range validatorErrors.Errors() {
		if errors, ok := validatorError.(govalidator.Errors); ok {
			for _, e := range errors {
				resultErrors = append(resultErrors, e.(govalidator.Error))
			}
		}
		if e, ok := validatorError.(govalidator.Error); ok {
			resultErrors = append(resultErrors, e)
		}
	}
	return resultErrors
}

func formattedError(err govalidator.Error, resource interface{}) error {
	message := err.Error()
	attrName := err.Name
	if strings.Index(message, "non zero value required") >= 0 {
		message = fmt.Sprintf("%v can't be blank", attrName)
	} else if strings.Index(message, "as length") >= 0 {
		reg, _ := regexp.Compile(`\(([0-9]+)\|([0-9]+)\)`)
		submatch := reg.FindSubmatch([]byte(err.Error()))
		message = fmt.Sprintf("%v is the wrong length (should be %v~%v characters)", attrName, string(submatch[1]), string(submatch[2]))
	} else if strings.Index(message, "as numeric") >= 0 {
		message = fmt.Sprintf("%v is not a number", attrName)
	} else if strings.Index(message, "as email") >= 0 {
		message = fmt.Sprintf("%v is not a valid email address", attrName)
	}
	return FieldFailed(resource, attrName, message)

}

var VALIDATE_CALLBACK = PREFIX + ":validate"

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *aorm.DB) *aorm.DB {
	db.Callback().Create().Before("gorm:before_create").Register(VALIDATE_CALLBACK, validate)
	db.Callback().Update().Before("gorm:before_update").Register(VALIDATE_CALLBACK, validate)
	return db.Set(VALIDATE_CALLBACK, true)
}

func RegisteredCallbacks(db *aorm.DB) bool {
	if _, ok := db.Get(VALIDATE_CALLBACK); ok {
		return true
	}
	return false
}

func RegisteredCallbacksOrError(db *aorm.DB) {
	if !RegisteredCallbacks(db) {
		panic(fmt.Errorf("%v: callbacks does not registered.", VALIDATE_CALLBACK))
	}
}
