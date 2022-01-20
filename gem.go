package gem

import (
	"fmt"
	"sort"
)

func Compute(load float64, fuels Fuels, plants []PowerPlant) ([]Item, error) {
	for i := range plants {
		plants[i] = setCost(plants[i], fuels)
	}
	sort.Slice(plants, func(i, j int) bool {
		return plants[i].cost < plants[j].cost && plants[i].Efficiency > plants[j].Efficiency
	})
	var (
		items []Item
		ok    bool
	)
	for i := 100.0; i > 0; i-- {
		items, ok = compute(load, plants, i/100.0)
		if ok {
			break
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("fail to compute power production")
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Power > items[j].Power })
	return items, nil
}

func compute(load float64, plants []PowerPlant, limit float64) ([]Item, bool) {
	var (
		total float64
		items = make([]Item, len(plants))
	)
	for i, p := range plants {
		var (
			pow = p.Max * limit
			sub = total + pow
		)
		items[i].Name = p.Name
		if p.Min > 0 && pow < p.Min || total == load {
			continue
		}
		if sub <= load {
			items[i].Power = Number(pow)
			total = sub
			continue
		}
		if diff := load - total; p.Min > 0 && diff > p.Min {
			items[i].Power = Number(diff)
			total = load
		}
	}
	return items, total == load
}

func setCost(p PowerPlant, fuels Fuels) PowerPlant {
	cost := fuels.Get(p.Type)
	if p.Type == wind {
		p.Efficiency = cost / 100
		p.Max *= p.Efficiency
	} else {
		p.cost = cost
	}
	return p
}
