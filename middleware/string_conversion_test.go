package middleware

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
	"github.com/casualjim/go-swagger/swag"
	"github.com/stretchr/testify/assert"
)

var evaluatesAsTrue = []string{"true", "1", "yes", "ok", "y", "on", "selected", "checked", "t", "enabled"}

type unmarshallerSlice []string

func (u *unmarshallerSlice) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return errors.New("an error")
	}
	*u = strings.Split(string(data), ",")
	return nil
}

type SomeOperationParams struct {
	Name        string
	ID          int64
	Confirmed   bool
	Age         int
	Visits      int32
	Count       int16
	Seq         int8
	UID         uint64
	UAge        uint
	UVisits     uint32
	UCount      uint16
	USeq        uint8
	Score       float32
	Rate        float64
	Timestamp   strfmt.DateTime
	Birthdate   strfmt.Date
	LastFailure *strfmt.DateTime
	Unsupported struct{}
	Tags        []string
	Prefs       []int32
	Categories  unmarshallerSlice
}

func FloatParamTest(t *testing.T, fName, pName, format string, val reflect.Value, defVal, expectedDef interface{}, actual func() interface{}) {
	fld := val.FieldByName(pName)
	binder := &untypedParamBinder{
		parameter: spec.QueryParam(pName).Typed("number", "double").WithDefault(defVal),
		Name:      pName,
	}

	err := binder.setFieldValue(fld, defVal, "5")
	assert.NoError(t, err)
	assert.EqualValues(t, 5, actual())

	err = binder.setFieldValue(fld, defVal, "")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedDef, actual())

	err = binder.setFieldValue(fld, defVal, "yada")
	assert.Error(t, err)
}

func IntParamTest(t *testing.T, pName string, val reflect.Value, defVal, expectedDef interface{}, actual func() interface{}) {
	fld := val.FieldByName(pName)

	binder := &untypedParamBinder{
		parameter: spec.QueryParam(pName).Typed("integer", "int64").WithDefault(defVal),
		Name:      pName,
	}
	err := binder.setFieldValue(fld, defVal, "5")
	assert.NoError(t, err)
	assert.EqualValues(t, 5, actual())

	err = binder.setFieldValue(fld, defVal, "")
	assert.NoError(t, err)
	assert.EqualValues(t, expectedDef, actual())

	err = binder.setFieldValue(fld, defVal, "yada")
	assert.Error(t, err)
}

func TestParamBinding(t *testing.T) {

	actual := new(SomeOperationParams)
	val := reflect.ValueOf(actual).Elem()
	pName := "Name"
	fld := val.FieldByName(pName)

	binder := &untypedParamBinder{
		parameter: spec.QueryParam(pName).Typed("string", "").WithDefault("some-name"),
		Name:      pName,
	}

	err := binder.setFieldValue(fld, "some-name", "the name value")
	assert.NoError(t, err)
	assert.Equal(t, "the name value", actual.Name)

	err = binder.setFieldValue(fld, "some-name", "")
	assert.NoError(t, err)
	assert.Equal(t, "some-name", actual.Name)

	IntParamTest(t, "ID", val, 1, 1, func() interface{} { return actual.ID })
	IntParamTest(t, "ID", val, nil, 0, func() interface{} { return actual.ID })
	IntParamTest(t, "Age", val, 1, 1, func() interface{} { return actual.Age })
	IntParamTest(t, "Age", val, nil, 0, func() interface{} { return actual.Age })
	IntParamTest(t, "Visits", val, 1, 1, func() interface{} { return actual.Visits })
	IntParamTest(t, "Visits", val, nil, 0, func() interface{} { return actual.Visits })
	IntParamTest(t, "Count", val, 1, 1, func() interface{} { return actual.Count })
	IntParamTest(t, "Count", val, nil, 0, func() interface{} { return actual.Count })
	IntParamTest(t, "Seq", val, 1, 1, func() interface{} { return actual.Seq })
	IntParamTest(t, "Seq", val, nil, 0, func() interface{} { return actual.Seq })
	IntParamTest(t, "UID", val, uint64(1), 1, func() interface{} { return actual.UID })
	IntParamTest(t, "UID", val, uint64(0), 0, func() interface{} { return actual.UID })
	IntParamTest(t, "UAge", val, uint(1), 1, func() interface{} { return actual.UAge })
	IntParamTest(t, "UAge", val, nil, 0, func() interface{} { return actual.UAge })
	IntParamTest(t, "UVisits", val, uint32(1), 1, func() interface{} { return actual.UVisits })
	IntParamTest(t, "UVisits", val, nil, 0, func() interface{} { return actual.UVisits })
	IntParamTest(t, "UCount", val, uint16(1), 1, func() interface{} { return actual.UCount })
	IntParamTest(t, "UCount", val, nil, 0, func() interface{} { return actual.UCount })
	IntParamTest(t, "USeq", val, uint8(1), 1, func() interface{} { return actual.USeq })
	IntParamTest(t, "USeq", val, nil, 0, func() interface{} { return actual.USeq })

	FloatParamTest(t, "score", "Score", "float", val, 1.0, 1, func() interface{} { return actual.Score })
	FloatParamTest(t, "score", "Score", "float", val, nil, 0, func() interface{} { return actual.Score })
	FloatParamTest(t, "rate", "Rate", "double", val, 1.0, 1, func() interface{} { return actual.Rate })
	FloatParamTest(t, "rate", "Rate", "double", val, nil, 0, func() interface{} { return actual.Rate })

	pName = "Confirmed"
	confirmedField := val.FieldByName(pName)
	binder = &untypedParamBinder{
		parameter: spec.QueryParam(pName).Typed("boolean", "").WithDefault(true),
		Name:      pName,
	}

	for _, tv := range evaluatesAsTrue {
		err = binder.setFieldValue(confirmedField, true, tv)
		assert.NoError(t, err)
		assert.True(t, actual.Confirmed)
	}

	err = binder.setFieldValue(confirmedField, true, "")
	assert.NoError(t, err)
	assert.True(t, actual.Confirmed)

	err = binder.setFieldValue(confirmedField, true, "0")
	assert.NoError(t, err)
	assert.False(t, actual.Confirmed)

	pName = "Timestamp"
	timeField := val.FieldByName(pName)
	dt := strfmt.DateTime{Time: time.Date(2014, 3, 19, 2, 9, 0, 0, time.UTC)}
	binder = &untypedParamBinder{
		parameter: spec.QueryParam(pName).Typed("string", "date-time").WithDefault(dt),
		Name:      pName,
	}
	exp := strfmt.DateTime{Time: time.Date(2014, 5, 14, 2, 9, 0, 0, time.UTC)}

	err = binder.setFieldValue(timeField, dt, exp.String())
	assert.NoError(t, err)
	assert.Equal(t, exp, actual.Timestamp)

	err = binder.setFieldValue(timeField, dt, "")
	assert.NoError(t, err)
	assert.Equal(t, dt, actual.Timestamp)

	err = binder.setFieldValue(timeField, dt, "yada")
	assert.Error(t, err)

	ddt := strfmt.Date{Time: time.Date(2014, 3, 19, 0, 0, 0, 0, time.UTC)}
	pName = "Birthdate"
	dateField := val.FieldByName(pName)
	binder = &untypedParamBinder{
		parameter: spec.QueryParam(pName).Typed("string", "date").WithDefault(ddt),
		Name:      pName,
	}
	expd := strfmt.Date{Time: time.Date(2014, 5, 14, 0, 0, 0, 0, time.UTC)}

	err = binder.setFieldValue(dateField, ddt, expd.String())
	assert.NoError(t, err)
	assert.Equal(t, expd, actual.Birthdate)

	err = binder.setFieldValue(dateField, ddt, "")
	assert.NoError(t, err)
	assert.Equal(t, ddt, actual.Birthdate)

	err = binder.setFieldValue(dateField, ddt, "yada")
	assert.Error(t, err)

	fdt := &strfmt.DateTime{Time: time.Date(2014, 3, 19, 2, 9, 0, 0, time.UTC)}
	pName = "LastFailure"
	ftimeField := val.FieldByName(pName)
	binder = &untypedParamBinder{
		parameter: spec.QueryParam(pName).Typed("string", "date").WithDefault(fdt),
		Name:      pName,
	}
	fexp := &strfmt.DateTime{Time: time.Date(2014, 5, 14, 2, 9, 0, 0, time.UTC)}

	err = binder.setFieldValue(ftimeField, fdt, fexp.String())
	assert.NoError(t, err)
	assert.Equal(t, fexp, actual.LastFailure)

	err = binder.setFieldValue(ftimeField, fdt, "")
	assert.NoError(t, err)
	assert.Equal(t, fdt, actual.LastFailure)

	err = binder.setFieldValue(ftimeField, fdt, "")
	assert.NoError(t, err)
	assert.Equal(t, fdt, actual.LastFailure)

	actual.LastFailure = nil
	err = binder.setFieldValue(ftimeField, fdt, "yada")
	assert.Error(t, err)
	assert.Nil(t, actual.LastFailure)

	pName = "Unsupported"
	unsupportedField := val.FieldByName(pName)
	binder = &untypedParamBinder{
		parameter: spec.QueryParam(pName).Typed("string", ""),
		Name:      pName,
	}
	err = binder.setFieldValue(unsupportedField, nil, "")
	assert.Error(t, err)
}

func TestSliceConversion(t *testing.T) {

	actual := new(SomeOperationParams)
	val := reflect.ValueOf(actual).Elem()

	// prefsField := val.FieldByName("Prefs")
	// cData := "yada,2,3"
	// _, _, err := readFormattedSliceFieldValue("Prefs", prefsField, cData, "csv", nil)
	// assert.Error(t, err)

	sliced := []string{"some", "string", "values"}
	seps := map[string]string{"ssv": " ", "tsv": "\t", "pipes": "|", "csv": ",", "": ","}

	tagsField := val.FieldByName("Tags")
	for k, sep := range seps {
		binder := &untypedParamBinder{
			Name:      "Tags",
			parameter: spec.QueryParam("tags").CollectionOf(stringItems, k),
		}

		actual.Tags = nil
		cData := strings.Join(sliced, sep)
		tags, _, err := binder.readFormattedSliceFieldValue(cData, tagsField)
		assert.NoError(t, err)
		assert.Equal(t, sliced, tags)
		cData = strings.Join(sliced, " "+sep+" ")
		tags, _, err = binder.readFormattedSliceFieldValue(cData, tagsField)
		assert.NoError(t, err)
		assert.Equal(t, sliced, tags)
		tags, _, err = binder.readFormattedSliceFieldValue("", tagsField)
		assert.NoError(t, err)
		assert.Empty(t, tags)
	}

	assert.Nil(t, swag.SplitByFormat("yada", "multi"))
	assert.Nil(t, swag.SplitByFormat("", ""))

	categoriesField := val.FieldByName("Categories")
	binder := &untypedParamBinder{
		Name:      "Categories",
		parameter: spec.QueryParam("categories").CollectionOf(stringItems, "csv"),
	}
	cData := strings.Join(sliced, ",")
	categories, custom, err := binder.readFormattedSliceFieldValue(cData, categoriesField)
	assert.NoError(t, err)
	assert.EqualValues(t, sliced, actual.Categories)
	assert.True(t, custom)
	assert.Empty(t, categories)
	categories, custom, err = binder.readFormattedSliceFieldValue("", categoriesField)
	assert.Error(t, err)
	assert.True(t, custom)
	assert.Empty(t, categories)
}
