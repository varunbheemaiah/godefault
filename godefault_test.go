package godefault

import (
	"reflect"
	"testing"
	"time"
)

type TestStruct struct {
	StringField           string          `default:"default string"`
	IntField              int             `default:"42"`
	BoolField             bool            `default:"true"`
	FloatField            float64         `default:"3.14"`
	DurationField         time.Duration   `default:"1h2m3s"`
	OmitField             string          `default:"should be omitted" json:"omitField,omitempty"`
	SliceOfStringsField   []string        `default:"a,b,c"`
	SliceOfIntsField      []int           `default:"1,2,3"`
	SliceOfFloaatsField   []float64       `default:"1.1,2.2,3.3"`
	SliceOfBoolsField     []bool          `default:"true,false,true"`
	SliceOfDurationsField []time.Duration `default:"1h,2h,3h"`
	NestedField           NestedStruct
}

type NestedStruct struct {
	NestedString string `default:"nested default"`
	NestedInt    int    `default:"7"`
}

func TestSetDefaults(t *testing.T) {
	test := TestStruct{}
	err := SetDefaults(&test)
	if err != nil {
		t.Errorf("SetDefaults returned an error: %v", err)
	}

	expected := TestStruct{
		StringField:           "default string",
		IntField:              42,
		BoolField:             true,
		FloatField:            3.14,
		DurationField:         time.Hour + 2*time.Minute + 3*time.Second,
		OmitField:             "",
		SliceOfStringsField:   []string{"a", "b", "c"},
		SliceOfIntsField:      []int{1, 2, 3},
		SliceOfFloaatsField:   []float64{1.1, 2.2, 3.3},
		SliceOfBoolsField:     []bool{true, false, true},
		SliceOfDurationsField: []time.Duration{time.Hour, 2 * time.Hour, 3 * time.Hour},
		NestedField: NestedStruct{
			NestedString: "nested default",
			NestedInt:    7,
		},
	}

	if !reflect.DeepEqual(test, expected) {
		t.Errorf("Structs are not equal.\nGot: %+v\nWant: %+v", test, expected)
	}
}
