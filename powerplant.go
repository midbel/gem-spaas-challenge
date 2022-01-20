package gem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	gas  = "gasfired"
	jet  = "turbojet"
	wind = "windturbine"
	co   = "co2"
)

type PowerPlant struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Min        float64 `json:"pmin"`
	Max        float64 `json:"pmax"`
	Efficiency float64 `json:"efficiency"`

	cost float64
}

type Fuels struct {
	prices map[string]float64
}

func (f *Fuels) Get(str string) float64 {
	return f.prices[str]
}

func (f *Fuels) MarshalJSON() ([]byte, error) {
	c := struct {
		Gas      float64 `json:"gas(euro/MWh)"`
		Kerosine float64 `json:"kerosine(euro/MWh)"`
		Co       float64 `json:"co2(euro/ton)"`
		Wind     float64 `json:"wind(%)"`
	}{
		Gas:      f.prices[gas],
		Kerosine: f.prices[jet],
		Co:       f.prices[co],
		Wind:     f.prices[wind],
	}
	return json.Marshal(c)
}

func (f *Fuels) UnmarshalJSON(b []byte) error {
	prices := make(map[string]float64)
	if err := json.Unmarshal(b, &prices); err != nil {
		return err
	}
	f.prices = make(map[string]float64)
	for k, v := range prices {
		var n string
		switch {
		case strings.HasPrefix(k, "gas"):
			n = gas
		case strings.HasPrefix(k, "kerosine"):
			n = jet
		case strings.HasPrefix(k, "co2"):
			n = co
		case strings.HasPrefix(k, "wind"):
			n = wind
		default:
			return fmt.Errorf("unsupported/unknown type %s", k)
		}
		f.prices[n] = v
	}
	return nil
}

type Number float64

func (n Number) MarshalJSON() ([]byte, error) {
	b :=  strconv.AppendFloat(nil, float64(n), 'f', 1, 64)
	if bytes.HasSuffix(b, []byte(".0")) {
		b = bytes.TrimSuffix(b, []byte(".0"))
	}
	return b, nil
}

type Item struct {
	Name  string `json:"name"`
	Power Number `json:"p"`
}
