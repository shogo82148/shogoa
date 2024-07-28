package design

import "testing"

func TestDataType_IsObject(t *testing.T) {
	tests := []struct {
		name     string
		dataType DataType
		want     bool
	}{
		{
			name:     "primitive",
			dataType: String,
			want:     false,
		},
		{
			name:     "array",
			dataType: &Array{ElemType: &AttributeDefinition{Type: String}},
			want:     false,
		},
		{
			name: "hash",
			dataType: &Hash{
				KeyType:  &AttributeDefinition{Type: String},
				ElemType: &AttributeDefinition{Type: String},
			},
			want: false,
		},
		{
			name: "nil user type",
			dataType: &UserTypeDefinition{
				AttributeDefinition: &AttributeDefinition{Type: nil},
			},
		},
		{
			name:     "object",
			dataType: &Object{},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dataType.IsObject(); got != tt.want {
				t.Errorf("DataType.IsObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
