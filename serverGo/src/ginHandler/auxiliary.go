package ginHandler

// ImageBundle The JSON structure to store an encoded base64 image and the corresponding text.
// Content of the text may vary depending on different usages.
type ImageBundle struct {
	EncodedImage string `json:"image"`
	Text string `json:"text"`
}

// ImageBundles A collection of ImageBundle s.
type ImageBundles struct {
	Images []ImageBundle `json:"images"`
}

// JSONShowPictures The JSON structure with offset and the number of pictures requested.
type JSONShowPictures struct {
	Offset int `json:"offset"`
	N int `json:"n"`
}

// JSONLabeledResult The JSON structure containing the name of the picture and the label created by the labeler.
// The label is a string and its meaning can be designed accordingly.
type JSONLabeledResult struct {
	Name string `json:"name"`
	Val string `json:"val"`
}

// JSONLabeledResults A collection of JSONLabeledResult s.
type JSONLabeledResults struct {
	Results []JSONLabeledResult `json:"results"`
}
