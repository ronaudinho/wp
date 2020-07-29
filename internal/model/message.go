package model

// Message is a model of message
// field CreateAt is not required, but added anyway in case
// query for all message implements filter by datetime
// TODO: query filter (not sure if will implement)
type Message struct {
	// sample use of validate tag
	// notoneof can be used to filter allowed words
	// tag is probably not the best place to do this though
	// would fit better for something that has lower variance
	// the downside is the need to use pointer in this field
	// to ensure that the application does not accept nil value
	// as without using pointer, there is no way to check
	// other than using zero value, which sometimes can be acceptable
	Message *string `json:"message" validate:"required notoneof=fuck,fck"`
	// Message  *string `json:"message" validate:"required"`
	CreateAt int32 `json:"create_at"`
}
