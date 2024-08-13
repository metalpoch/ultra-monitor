package entity

type Example struct {
	Message string `bson:"message"`
	Value   uint8  `bson:"value"`
}

type ExampleResponse struct {
	Message string `bson:"message"`
}
