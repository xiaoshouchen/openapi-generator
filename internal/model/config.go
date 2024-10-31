package model

type Config struct {
	AimType     string  `json:"aim_type"`
	OutPath     string  `json:"out_path"`
	ProjectName string  `json:"package_name"`
	Fetcher     Fetcher `json:"fetcher"`
}

type Fetcher struct {
	From   string `json:"from"`
	Apifox struct {
		Token     string `json:"token"`
		ProjectId int64  `json:"project_id"`
	} `json:"apifox"`
	File string `json:"file"`
}
