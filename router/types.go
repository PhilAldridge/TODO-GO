package router

import "github.com/PhilAldridge/TODO-GO/store"

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

type TodoApiHandlerV2 struct {
	actor    *store.StoreActor
	username string
}

type UserPutBody struct {
	Username string
	Password string
}
