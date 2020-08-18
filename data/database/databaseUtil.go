package database

import "time"

func parseProp(item interface{}, defaultVal interface{}) interface{} {
	if item == nil {
		return defaultVal
	}
	return item
}

func parseBoolProp(row interface{}, defaultValue bool) bool {
	return parseProp(row, defaultValue).(bool)
}
func parseIntProp(row interface{}, defaultValue int64) int64 {
	return parseProp(row, defaultValue).(int64)
}
func parseFloatProp(row interface{}, defaultValue float64) float64 {
	return parseProp(row, defaultValue).(float64)
}
func parseStringProp(row interface{}, defaultValue string) string {
	return parseProp(row, defaultValue).(string)
}
func parseByteArrayProp(row interface{}, defaultValue []byte) []byte {
	return []byte(parseStringProp(row, string(defaultValue)))
}
func parseTimeProp(row interface{}, defaultValue time.Time) time.Time {
	return time.Unix(parseIntProp(row, defaultValue.Unix()), 0)
}
