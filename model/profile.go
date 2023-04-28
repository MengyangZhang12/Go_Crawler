package model

import "encoding/json"

type Profile struct {
	Name      string
	UserUrl   string
	UserId    string
	Age       string
	Gender    string
	Marriage  string
	Location  string
	Height    string
	Education string
	Income    string
	ImageUrl  string
	ImageList []string
	Introduce string
}

func FromJsonObj(o interface{}) (Profile, error) {
	var profile Profile
	s, err := json.Marshal(o)
	if err != nil {
		return profile, err
	}
	err = json.Unmarshal(s, &profile)
	return profile, err
}
