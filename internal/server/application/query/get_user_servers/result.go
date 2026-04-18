package getuserservers

type Result struct {
	Servers []ResultServer `json:"servers"`
}

type ResultServer struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Members []ResultMember `json:"members"`
}

type ResultMember struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
