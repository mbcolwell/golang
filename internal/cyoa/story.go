package cyoa

import "encoding/json"

func ParseJson(filename string) (story, error) {
	storyData, err := jsonFiles.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var s story
	err = json.Unmarshal(storyData, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
