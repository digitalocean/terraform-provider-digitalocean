package datalist

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func valueMatches(s *schema.Schema, value interface{}, filterValue string) bool {
	switch s.Type {
	case schema.TypeString:
		return strings.EqualFold(filterValue, value.(string))

	case schema.TypeBool:
		if boolValue, err := strconv.ParseBool(filterValue); err == nil {
			return boolValue == value.(bool)
		}

	case schema.TypeInt:
		if intValue, err := strconv.Atoi(filterValue); err == nil {
			return intValue == value.(int)
		}

	case schema.TypeFloat:
		if floatValue, err := strconv.ParseFloat(filterValue, 64); err == nil {
			return floatValue == value.(float64)
		}

	case schema.TypeList:
		listValues := value.([]interface{})
		result := false
		for _, listValue := range listValues {
			result = result || valueMatches(s.Elem.(*schema.Schema), listValue, filterValue)
		}
		return result
	}

	return false
}

func compareValues(s *schema.Schema, value1 interface{}, value2 interface{}) int {
	switch s.Type {
	case schema.TypeString:
		return strings.Compare(value1.(string), value2.(string))

	case schema.TypeBool:
		boolValue1 := value1.(bool)
		boolValue2 := value2.(bool)
		if boolValue1 == boolValue2 {
			return 0
		} else if !boolValue1 {
			return -1
		} else {
			return 1
		}

	case schema.TypeInt:
		intValue1 := value1.(int)
		intValue2 := value2.(int)
		if intValue1 < intValue2 {
			return -1
		} else if intValue1 > intValue2 {
			return 1
		} else {
			return 0
		}

	case schema.TypeFloat:
		floatValue1 := value1.(float64)
		floatValue2 := value2.(float64)
		if floatValue1 < floatValue2 {
			return -1
		} else if floatValue1 > floatValue2 {
			return 1
		} else {
			return 0
		}

	default:
		panic("Illegal state: Unsupported value type for sort")
	}
}
