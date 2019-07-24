package data

import "time"

func parseRow( item interface{}, defaultVal interface{}) interface{}{
	if item == nil{
		return defaultVal
	}
	return item

}

func parseRowBool(row interface{}, defaultValue bool) bool{
	return parseRow(row, defaultValue).(bool)
}
func parseRowInt(row interface{}, defaultValue int64) int64{
	return parseRow(row, defaultValue).(int64)
}
func parseRowFloat(row interface{}, defaultValue float64) float64{
	return parseRow(row, defaultValue).(float64)
}
func parseRowString(row interface{}, defaultValue string) string{
	return parseRow(row, defaultValue).(string)
}
func parseRowBytes(row interface{}, defaultValue []byte) []byte{
	return []byte(parseRowString(row,string(defaultValue)))
}

func parseRowTime(row interface{}, defaultValue time.Time) time.Time{
	return time.Unix(parseRowInt(row,defaultValue.Unix()), 0)
}
