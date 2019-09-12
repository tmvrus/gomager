package text

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalError(t *testing.T) {
	resp := apiResponse{}
	if err := json.Unmarshal(errorData, &resp); err != nil {
		t.Fatalf("unexpected error: %q", err.Error())
	}

	if !resp.HasError {
		t.Fatal("expected HasError be true")
	}

	if len(resp.ErrorMessage) != 1 {
		t.Fatal("unexpected ErrorMessage len")
	}

	if resp.ErrorMessage[0] != "Unable to recognize the file type" {
		t.Fatalf("invalid error message: %q", resp.ErrorMessage[0])
	}
}
func TestUnmarshal(t *testing.T) {
	resp := apiResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("unexpected error: %q", err.Error())
	}

	if len(resp.Result) != 1 {
		t.Fatal("unexpected result len")
	}

	if len(resp.Result[0].Overlay.Lines) != 1 {
		t.Fatal("unexpected lines len")
	}

	if len(resp.Result[0].Overlay.Lines[0].Words) != 2 {
		t.Fatal("unexpected words len")
	}

	if h := resp.Result[0].Overlay.Lines[0].Words[0].Height; h != 13 {
		t.Fatalf("invalid Height value %f", h)
	}

	if w := resp.Result[0].Overlay.Lines[0].Words[1].Width; w != 66.1 {
		t.Fatalf("invalid Width value %f", w)
	}
}

var data = []byte(`{
  "ParsedResults": [
    {
      "TextOverlay": {
        "Lines": [
          {
            "LineText": "4 40:02PM",
            "Words": [
              {
                "WordText": "4",
                "Left": 204,
                "Top": 203,
                "Height": 13,
                "Width": 8
              },
              {
                "WordText": "40:02PM",
                "Left": 219,
                "Top": 202,
                "Height": 15,
                "Width": 66.1
              }
            ],
            "MaxHeight": 15,
            "MinTop": 202
          }
        ],
        "HasOverlay": true,
        "Message": "Total lines: 1"
      },
      "FileParseExitCode": 1,
      "TextOrientation": "0",
      "ParsedText": "4 40:02PM\r\n",
      "ErrorMessage": "",
      "ErrorDetails": ""
    }
  ],
  "OCRExitCode": 1,
  "IsErroredOnProcessing": false,
  "ProcessingTimeInMilliseconds": 0.859,
  "SearchablePDFURL": "Searchable PDF not generated as it was not requested."
}`)

var errorData = []byte(`{"OCRExitCode":99,"IsErroredOnProcessing":true,"ErrorMessage":["Unable to recognize the file type"],"ProcessingTimeInMilliseconds":"23"}`)
