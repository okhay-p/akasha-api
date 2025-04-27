package req_structs

type GenerateTopicReqBody struct {
	Content      string `json:"content"`
	IsPublic     bool   `json:"is_public"`
	OnlyContent  bool   `json:"only_content"`
	NumOfLessons uint8  `json:"num_of_lessons"`
}
