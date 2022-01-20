package gem

import (
	"encoding/json"
	"strings"
	"testing"
)

const sample = `
{
  "load": 480,
  "fuels":
  {
    "gas(euro/MWh)": 13.4,
    "kerosine(euro/MWh)": 50.8,
    "co2(euro/ton)": 20,
    "wind(%)": 0
  },
  "powerplants": [
    {
      "name": "gasfiredbig1",
      "type": "gasfired",
      "efficiency": 0.53,
      "pmin": 100,
      "pmax": 460
    },
    {
      "name": "gasfiredbig2",
      "type": "gasfired",
      "efficiency": 0.53,
      "pmin": 100,
      "pmax": 460
    },
    {
      "name": "gasfiredsomewhatsmaller",
      "type": "gasfired",
      "efficiency": 0.37,
      "pmin": 40,
      "pmax": 210
    },
    {
      "name": "tj1",
      "type": "turbojet",
      "efficiency": 0.3,
      "pmin": 0,
      "pmax": 16
    },
    {
      "name": "windpark1",
      "type": "windturbine",
      "efficiency": 1,
      "pmin": 0,
      "pmax": 150
    },
    {
      "name": "windpark2",
      "type": "windturbine",
      "efficiency": 1,
      "pmin": 0,
      "pmax": 36
    }
  ]
}
`

func TestCompute(t *testing.T) {
	in := struct {
		Load   float64      `json:"load"`
		Fuels  Fuels        `json:"fuels"`
		Plants []PowerPlant `json:"powerplants"`
	}{}
	if err := json.NewDecoder(strings.NewReader(sample)).Decode(&in); err != nil {
		t.Fatalf("fail to decode input data")
	}
	items, _ := Compute(in.Load, in.Fuels, in.Plants)
	if len(items) != len(in.Plants) {
		t.Fatalf("not enough items in return array! want %d, got %d", len(in.Plants), len(items))
	}
	var total float64
	for i := range items {
		total += float64(items[i].Power)
	}
	if total != in.Load {
		t.Fatalf("results mismatched! expected load %.0f, got %.0f", in.Load, total)
	}
}
