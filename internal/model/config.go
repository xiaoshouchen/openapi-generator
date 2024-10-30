package model

type Config struct {
	AimType     string  `json:"aim_type"`
	OutPath     string  `json:"out_path"`
	ProjectName string  `json:"package_name"`
	Fetcher     Fetcher `json:"fetcher"`
}

type Fetcher struct {
	From   string `json:"from"`
	ApiFox struct {
		Token     string `json:"token"`
		ProjectId int    `json:"project_id"`
	} `json:"api_fox"`
	File string `json:"file"`
}
