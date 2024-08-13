package api

type EchoRequest struct {
	Name  string `json:"name"`
	Check bool   `json:"check"`
}

type EchoResponse struct {
	Name  string `json:"name"`
	Check bool   `json:"check"`
}

func APIEcho(req EchoRequest) (EchoResponse, error) {
	return EchoResponse{
		Name:  req.Name,
		Check: req.Check,
	}, nil
}
