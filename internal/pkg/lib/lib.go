package lib


type Database interface {
	UserCrud
	RefreshTokenDatabase
}