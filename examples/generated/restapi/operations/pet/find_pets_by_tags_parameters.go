package pet

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/middleware"
	"github.com/casualjim/go-swagger/strfmt"
	"github.com/casualjim/go-swagger/swag"
)

// FindPetsByTagsParams contains all the bound params for the find pets by tags operation
// typically these are obtained from a http.Request
type FindPetsByTagsParams struct {
	// Tags to filter by
	Tags []string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *FindPetsByTagsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	qs := r.URL.Query()

	if err := o.bindTags(swag.SplitByFormat(qs.Get("tags"), "multi"), route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *FindPetsByTagsParams) bindTags(raw []string, formats strfmt.Registry) error {
	size := len(raw)

	if size == 0 {
		return nil
	}

	ic := raw
	isz := size
	var ir []string
	iValidateElement := func(i int, value string) *errors.Validation {

		return nil
	}

	for i := 0; i < isz; i++ {

		if err := iValidateElement(i, ic[i]); err != nil {
			return err
		}
		ir = append(ir, ic[i])
	}

	o.Tags = ir

	return nil
}
