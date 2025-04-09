package utils

import (
	"sort"

	"github.com/metalpoch/olt-blueprint/common/model"
)

type TrafficStates []*model.TrafficState

// Len devuelve la longitud del slice.
func (ts TrafficStates) Len() int {
	return len(ts)
}

// Less devuelve true si el elemento en el índice i tiene un valor Out mayor que el elemento en el índice j.
func (ts TrafficStates) Less(i, j int) bool {
	return ts[i].Out > ts[j].Out
}

// Swap intercambia los elementos en los índices i y j.
func (ts TrafficStates) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

// SortTrafficStatesByOut ordena un slice de punteros a model.TrafficState de forma descendente
// según el valor del campo Out y devuelve los primeros n elementos.
func SortTrafficStatesByOut(trafficStates []*model.TrafficState, n int) []*model.TrafficState {
	sortedTrafficStates := make(TrafficStates, len(trafficStates))
	copy(sortedTrafficStates, trafficStates)
	sort.Sort(sortedTrafficStates)

	if n > 0 && n < len(sortedTrafficStates) {
		return sortedTrafficStates[:n]
	} else if n >= len(sortedTrafficStates) {
		return sortedTrafficStates
	} else {
		return []*model.TrafficState{}
	}
}
