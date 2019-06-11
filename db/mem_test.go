package db

import (
	"reflect"
	"testing"
)

func Test_inMemDb_Get(t *testing.T) {
	tests := []struct {
		name       string
		existing   map[string]string
		input      string
		wantResult string
		wantExists bool
	}{
		{
			name:       "Non existent field",
			existing:   map[string]string{},
			input:      "somethingRandom",
			wantExists: false,
		},
		{
			name: "Non existent field with other fields",
			existing: map[string]string{
				"existing_field": "existing_value",
			},
			input:      "somethingRandom",
			wantExists: false,
		},
		{
			name: "Select Only Field",
			existing: map[string]string{
				"existing_field": "existing_value",
			},
			input:      "existing_field",
			wantResult: "existing_value",
			wantExists: true,
		},
		{
			name: "Select Only Of Many",
			existing: map[string]string{
				"existing_field":   "existing_value",
				"existing_field_1": "existing_value_1",
				"existing_field_2": "existing_value_2",
				"existing_field_3": "existing_value_3",
			},
			input:      "existing_field",
			wantResult: "existing_value",
			wantExists: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &inMemDb{
				m: tt.existing,
			}
			gotResult, gotExists := db.Get(tt.input)
			if gotResult != tt.wantResult {
				t.Errorf("inMemDb.Get() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotExists != tt.wantExists {
				t.Errorf("inMemDb.Get() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func Test_inMemDb_AddMapping(t *testing.T) {
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name      string
		existing  map[string]string
		args      args
		wantErr   bool
		wantState map[string]string
	}{
		{
			name:     "first add",
			existing: map[string]string{},
			args:     args{"a", "b"},
			wantErr:  false,
			wantState: map[string]string{
				"a": "b",
			},
		},
		{
			name: "unique add",
			existing: map[string]string{
				"existing_field": "existing_value",
			},
			args:    args{"a", "b"},
			wantErr: false,
			wantState: map[string]string{
				"existing_field": "existing_value",
				"a":              "b",
			},
		},
		{
			name: "collision add",
			existing: map[string]string{
				"existing_field": "existing_value",
			},
			args:    args{"existing_field", "b"},
			wantErr: true,
			wantState: map[string]string{
				"existing_field": "existing_value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &inMemDb{
				m: tt.existing,
			}
			if err := db.AddMapping(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("inMemDb.AddMapping() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(db.m, tt.wantState) {
				t.Errorf("inMemDb.AddMapping() state = %v, wantState %v", db.m, tt.wantState)
			}
		})
	}
}

func Test_inMemDb_RemoveEntry(t *testing.T) {
	tests := []struct {
		name      string
		existing  map[string]string
		input     string
		wantErr   bool
		wantState map[string]string
	}{
		{
			name:      "remove from empty",
			existing:  map[string]string{},
			input:     "toRemove",
			wantErr:   false,
			wantState: map[string]string{},
		},
		{
			name: "remove only field",
			existing: map[string]string{
				"toRemove": "sol",
			},
			input:     "toRemove",
			wantErr:   false,
			wantState: map[string]string{},
		},
		{
			name: "remove one of many",
			existing: map[string]string{
				"toRemove":         "sol",
				"existing_field":   "existing_value",
				"existing_field_1": "existing_value_1",
				"existing_field_2": "existing_value_2",
				"existing_field_3": "existing_value_3",
			},
			input:   "toRemove",
			wantErr: false,
			wantState: map[string]string{
				"existing_field":   "existing_value",
				"existing_field_1": "existing_value_1",
				"existing_field_2": "existing_value_2",
				"existing_field_3": "existing_value_3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &inMemDb{
				m: tt.existing,
			}
			if err := db.RemoveEntry(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("inMemDb.RemoveEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(db.m, tt.wantState) {
				t.Errorf("inMemDb.AddMapping() state = %v, wantState %v", db.m, tt.wantState)
			}
		})
	}
}
