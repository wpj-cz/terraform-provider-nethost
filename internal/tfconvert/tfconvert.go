package tfconvert

import (
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Int64FromAny(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case json.Number:
		i, err := v.Int64()
		if err == nil {
			return i
		}
		f, err := v.Float64()
		if err == nil {
			return int64(f)
		}
		return 0
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return i
		}
		return 0
	default:
		return 0
	}
}

func StringsToTypes(values []string) []types.String {
	result := make([]types.String, len(values))
	for i, value := range values {
		result[i] = types.StringValue(value)
	}
	return result
}

func TypesToStrings(values []types.String) []string {
	if values == nil {
		return nil
	}

	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.ValueString())
	}
	return result
}
