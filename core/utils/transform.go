package utils

func DeleteDuplicate(ids []uint) []uint {
	uniqueIDs := make(map[uint]bool)
	uniqueIDList := []uint{}

	// Obtener IDs únicos
	for _, id := range ids {
		if _, ok := uniqueIDs[id]; !ok {
			uniqueIDs[id] = true
			uniqueIDList = append(uniqueIDList, id)
		}
	}
	return uniqueIDList
}
