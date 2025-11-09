package lib

import "encoding/json"

func UnmarshalInto[T any](unmarshalled any, params *T) (*T, error) {
	marshalled, err := json.Marshal(unmarshalled)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(marshalled, &params)
	return params, nil
}
