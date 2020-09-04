package validations

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/moisespsena-go/aorm"
	error_utils "github.com/unapu-go/error-utils"
)

// NewError generate a new error for a model's field
func NewError(resource interface{}, column, err string) error {
	return &Error{Resource: resource, Column: column, Message: err}
}

// Error is a validation error struct that hold model, column and error message
type Error struct {
	Resource interface{}
	Column   string
	Message  string
}

// Label is a label including model type, primary key and column name
func (err Error) Label() string {
	struc := aorm.StructOf(err.Resource)
	return fmt.Sprintf("%v_%s_%v", struc.Type.Name(), IDToString(struc.GetID(err.Resource)), err.Column)
}

// Error show error message
func (err Error) Error() string {
	return fmt.Sprintf("%v", err.Message)
}

// IsError returns if any error is Error type
func IsError(err ...error) bool {
	return error_utils.IsErrorTyp(reflect.TypeOf(Error{}), err...)
}

func IDToString(id aorm.ID) string {
	if id == nil {
		return ""
	}
	return hex.EncodeToString(id.Bytes())
}
