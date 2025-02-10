package decode

import "testing"

func TestDecoding(t *testing.T) {
	inputs := []string{`{"ip":"1.1.1.1","port":7777,"service":"DNS","timestamp":123456789,"data_version":1,"data":{"response_bytes_utf8":[116,101,115,116]}}`,
		`{"ip":"1.1.1.1","port":7777,"service":"DNS","timestamp":123456789,"data_version":2,"data":{"response_str":"test"}}`}

	for _, tt := range inputs {
		decoded, err := DecodeMessage([]byte(tt))
		if err != nil {
			t.Fatalf("error parsing message: %s", err)
		}

		scan := decoded.ParsedScan

		if scan.Ip != "1.1.1.1" {
			t.Fatalf("IP address is wrong; got %s, want 1.1.1.1", scan.Ip)
		}

		if scan.Port != 7777 {
			t.Fatalf("port is wrong; got %d, want 7777", scan.Port)
		}

		if decoded.Response != "test" {
			t.Fatalf("the response string is wrong; got %s, want 'test'", decoded.Response)
		}
	}
}
