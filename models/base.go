package models

type Gobject struct {
	ID         int `json:"id"`
	Name, Desc string
	Cat_id     int
	Lat        float64
	Lon        float64
}

type LinkTrek struct {
	ID        int `json:"id"`
	Short     string
	Cat_id    int
	Objid     int
	Remote    string
	Published bool
}

func (LinkTrek) TableName() string {
	return "linktrek_link"
}

type Click struct {
	ID                int `json:"id"`
	Created           string
	Link_id           int
	Referrer          string
	Ip_address        string
	User_agent        string
	User_agent_source string
}

func (Click) TableName() string {
	return "linktrek_click"
}

func (Gobject) TableName() string {
	return "objects_gobject"
}
