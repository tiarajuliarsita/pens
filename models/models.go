package models

type Pen struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func (m Pen) Validate() error {
	cErr := []string{}
	if m.Name == "" {
		cErr = append(cErr, "name cannot nnull")
	}
	if m.Price == 0.0 {
		cErr = append(cErr, "PRICE cannot null")
	}
	if len(cErr) > 0 {
		return NewBadRequest(cErr)
	}
	return nil
}

func (m Pen) Exist() error {
	if m.ID == 0 {
		return NewNotFound("PEN NOT FOUND")
	}
	return nil
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
