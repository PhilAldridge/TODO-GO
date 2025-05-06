package router

type V1PutBody struct {
	Label    string
	Deadline string
}

type V1PatchBody struct {
	Id    string
	Field string
	Value string
}

type V1DeleteBody struct {
	Id string
}

type TodoApiHandler struct {
	actor chan<- StoreCommand
}

type TodoApiHandlerV2 struct {
	actor chan<- StoreCommand
	username string
}

type UserPutBody struct {
	Username string
	Password string
}
