package user_service

type GetByIdRequest struct {
	Id int `json:"id"`
}

type GetByIdResponse struct {
	Data string `json:"data"`
}
