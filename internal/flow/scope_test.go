package flow

import (
	"reflect"
	"testing"
)

func TestFlattenScope(t *testing.T) {
	tests := []struct {
		name   string
		global map[string]AnswerNode
		stack  []LoopFrame
		want   map[string]interface{}
	}{
		{
			name:   "empty",
			global: map[string]AnswerNode{},
			stack:  nil,
			want:   map[string]interface{}{},
		},
		{
			name:   "global only",
			global: map[string]AnswerNode{"project_name": "x", "linter": "biome"},
			stack:  nil,
			want:   map[string]interface{}{"project_name": "x", "linter": "biome"},
		},
		{
			name:   "single loop frame overrides global",
			global: map[string]AnswerNode{"project_name": "x", "name": "global"},
			stack: []LoopFrame{
				{QuestionID: "services", Index: 0, Answers: map[string]AnswerNode{"name": "auth"}},
			},
			want: map[string]interface{}{"project_name": "x", "name": "auth"},
		},
		{
			name:   "deeper frame overrides shallower",
			global: map[string]AnswerNode{"project_name": "x"},
			stack: []LoopFrame{
				{QuestionID: "services", Index: 0, Answers: map[string]AnswerNode{"service_name": "auth", "name": "auth"}},
				{QuestionID: "tables", Index: 0, Answers: map[string]AnswerNode{"name": "users"}},
				{QuestionID: "columns", Index: 0, Answers: map[string]AnswerNode{"name": "id", "type": "uuid"}},
			},
			want: map[string]interface{}{
				"project_name": "x",
				"service_name": "auth",
				"name":         "id",
				"type":         "uuid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenScope(tt.global, tt.stack)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FlattenScope() = %v, want %v", got, tt.want)
			}
		})
	}
}
