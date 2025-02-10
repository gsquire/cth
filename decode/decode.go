package decode

import "encoding/json"

type Scan struct {
	Ip          string          `json:"ip"`
	Port        uint32          `json:"port"`
	Service     string          `json:"service"`
	Timestamp   int64           `json:"timestamp"`
	DataVersion int             `json:"data_version"`
	Data        json.RawMessage `json:"data"`
}

type V1Data struct {
	ResponseBytesUtf8 []byte `json:"response_bytes_utf8"`
}

type V2Data struct {
	ResponseStr string `json:"response_str"`
}

// This type represents a fully decoded message with the response as a string type.
type DecodedMessage struct {
	ParsedScan *Scan
	Response   string
}

// We need to check what data version we get and subsequently parse the proper type.
func DecodeMessage(data []byte) (*DecodedMessage, error) {
	var scan Scan
	err := json.Unmarshal(data, &scan)
	if err != nil {
		return nil, err
	}

	var response string
	if scan.DataVersion == 1 {
		var v1 V1Data
		err = json.Unmarshal(scan.Data, &v1)
		if err != nil {
			return nil, err
		}
		response = string(v1.ResponseBytesUtf8)
	} else {
		var v2 V2Data
		err = json.Unmarshal(scan.Data, &v2)
		if err != nil {
			return nil, err
		}
		response = v2.ResponseStr
	}

	return &DecodedMessage{ParsedScan: &scan, Response: response}, nil
}
