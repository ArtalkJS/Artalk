package env_provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NumberKeysMapToSlice(t *testing.T) {
	result := NumberKeysMapToSlice(map[string]interface{}{
		"NormalMap": map[string]interface{}{
			"0": "1",
			"1": "2",
			"2": "3",
		},
		"NormalSlice": []interface{}{"1", "2", "3"},
		"Normal":      "b",
		"Map_nest_Map": map[string]interface{}{
			"c": map[string]interface{}{
				"d": "d",
			},
		},
		"Slice_nest_Slice": []interface{}{
			[]interface{}{"a", "b", "c"},
		},
		"Slice_nest_Map": []interface{}{
			map[string]interface{}{
				"1": "a",
				"2": "b",
				"3": "c",
			},
		},
		"Map_nest_Slice": map[string]interface{}{
			"c": map[string]interface{}{
				"1": "a",
				"2": "b",
				"3": "c",
			},
		},
		"MapSlice_nest_Slice": []map[string]interface{}{
			{"a": []interface{}{"1", "2", "3"}},
		},
		"Slice_nest_MapSlice": []interface{}{
			[]map[string]interface{}{
				{"a": "1"},
				{"b": "2"},
				{"c": "3"},
			},
		},
		"DeepMixed": map[string]interface{}{
			"0": map[string]interface{}{
				"1": []interface{}{
					"2",
					map[string]interface{}{
						"4": "5",
					},
				},
				"2": "3",
				"4": map[string]interface{}{
					"5": "6",
					"6": []interface{}{
						"7",
					},
					"8": map[string]interface{}{
						"2": "10",
					},
				},
			},
			"a": "b",
			"1": map[string]interface{}{
				"0": nil,
				"1": "2",
			},
		},
	})

	assert.Equal(t, map[string]interface{}{
		"NormalMap":   []interface{}{"1", "2", "3"},
		"NormalSlice": []interface{}{"1", "2", "3"},
		"Normal":      "b",
		"Map_nest_Map": map[string]interface{}{
			"c": map[string]interface{}{
				"d": "d",
			},
		},
		"Slice_nest_Map": []interface{}{
			[]interface{}{nil, "a", "b", "c"},
		},
		"Slice_nest_Slice": []interface{}{
			[]interface{}{"a", "b", "c"},
		},
		"Map_nest_Slice": map[string]interface{}{
			"c": []interface{}{nil, "a", "b", "c"},
		},
		"MapSlice_nest_Slice": []map[string]interface{}{
			{"a": []interface{}{"1", "2", "3"}},
		},
		"Slice_nest_MapSlice": []interface{}{
			[]map[string]interface{}{
				{"a": "1"},
				{"b": "2"},
				{"c": "3"},
			},
		},
		"DeepMixed": map[string]interface{}{
			"0": []interface{}{
				nil,
				[]interface{}{
					"2",
					[]interface{}{nil, nil, nil, nil, "5"},
				},
				"3",
				nil,
				[]interface{}{
					nil, nil, nil, nil, nil, "6", []interface{}{"7"}, nil, []interface{}{nil, nil, "10"},
				},
			},
			"a": "b",
			"1": []interface{}{
				nil,
				"2",
			},
		},
	}, result)
}
