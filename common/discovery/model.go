package discovery

import "encoding/json"

type EndpointInfo struct {
	IP       string         `json:"ip"`
	Port     string         `json:"port"`
	MetaData map[string]any `json:"meta"`
}

func (e *EndpointInfo) Marshal() string {
	marshalData, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(marshalData)
}
func UnMarshal(data []byte) (*EndpointInfo, error) {
	endpointStruck := &EndpointInfo{}
	err := json.Unmarshal(data, endpointStruck)
	if err != nil {
		return nil, err
	}
	return endpointStruck, nil
}
