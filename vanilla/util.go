package vanilla

import (
	"fmt"
	"strconv"
)

func ExtractUniqueIds(datas []IIDable, idType string) []int {
	id2bool := make(map[int]bool)

	for _, data := range datas {
		id := data.GetId(idType)
		id2bool[id] = true
	}
	
	ids := make([]int, 0)
	for id := range id2bool {
		ids = append(ids, id)
	}
	
	return ids
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func init() {
}