package client

type infoApp struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Info struct {
	App infoApp `json:"app"`
}
