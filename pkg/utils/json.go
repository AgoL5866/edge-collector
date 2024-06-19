package utils

import jsoniter "github.com/json-iterator/go"

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	// Marshal is exported by gin/json package.
	JSONMarshal = json.Marshal
	// Unmarshal is exported by gin/json package.
	JSONUnmarshal = json.Unmarshal
	// MarshalIndent is exported by gin/json package.
	JSONMarshalIndent = json.MarshalIndent
	// NewDecoder is exported by gin/json package.
	JSONNewDecoder = json.NewDecoder
	// NewEncoder is exported by gin/json package.
	JSONNewEncoder = json.NewEncoder
)
