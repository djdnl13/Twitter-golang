package model

type Tweet struct {
	Id       string `json:"id"`
        Text     string `json:"text"`
        LikesCount    string `json:"likesCount"`
        AccountId  string `json:"accountId"`
	ServedBy string `json:"servedBy"`
}

func (a *Tweet) ToString() string {
	return a.Id + " " + a.Text
}
