package model

// Message
type Message struct {
	// sample use of validate tag
	// notoneof can be used to filter allowed words
	// tag is probably not the best place to do this though
	// would fit better for something that has lower variance
	Message *string `json:"message" validate:"required notoneof=fuck,fck"`
	// Message  *string `json:"message" validate:"required"`
	CreateAt int32 `json:"create_at"`
}
