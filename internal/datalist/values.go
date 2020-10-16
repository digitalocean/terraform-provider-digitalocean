package datalist

import (
	"math"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func floatApproxEquals(a, b float64) bool {
	return math.Abs(a-b) < 0.000001
}

func valueMatches(s *schema.Schema, value interface{}, filterValue interface{}, matchBy string) bool {
	switch s.Type {
	case schema.TypeString:
		switch matchBy {
		case "exact":
			return strings.EqualFold(filterValue.(string), value.(string))
		case "substring":
			return strings.Contains(value.(string), filterValue.(string))
		case "re":
			return filterValue.(*regexp.Regexp).MatchString(value.(string))
		}

	case schema.TypeBool:
		return filterValue.(bool) == value.(bool)

	case schema.TypeInt:
		return filterValue.(int) == value.(int)

	case schema.TypeFloat:
		return floatApproxEquals(filterValue.(float64), value.(float64))

	case schema.TypeList:
		listValues := value.([]interface{})
		result := false
		for _, listValue := range listValues {
			valueDoesMatch := valueMatches(s.Elem.(*schema.Schema), listValue, filterValue, matchBy)
			result = result || valueDoesMatch
		}
		return result

	case schema.TypeSet:
		setValue := value.(*schema.Set)
		listValues := setValue.List()
		result := false
		for _, listValue := range listValues {
			valueDoesMatch := valueMatches(s.Elem.(*schema.Schema), listValue, filterValue, matchBy)
			result = result || valueDoesMatch
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
		if floatApproxEquals(floatValue1, floatValue2) {
			return 0
		} else if floatValue1 < floatValue2 {
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
