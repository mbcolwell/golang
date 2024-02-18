package cyoa

type story map[string]arc

type arc struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []option `json:"options"`
}

type option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
