package domain

type Book struct {
	ID              string
	Title           string
	Author          string
	PublicationYear int32
	Genre           string
}

type BookFilter struct {
	Author          *string
	PublicationYear *int32
	Genre           *string
}
