package api

type StatRequest struct {
	Name  string `json:"name"`
	Check bool   `json:"check"`
}

type StatResponse struct {
	Name  string `json:"name"`
	Check bool   `json:"check"`
}

func APIStat(req StatRequest) (StatResponse, error) {
	return StatResponse{
		Name:  req.Name,
		Check: req.Check,
	}, nil
}
