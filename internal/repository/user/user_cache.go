package repository

import "fmt"

const (
	UsernameKeyPrefix = "username" // username:username
	IDKeyPrefix       = "id"       // id:userid
	AdminKey          = "admin"    // admin:userid
	SysAdminKey       = "sysadmin"
)

func GetUserIDKey(id uint) string {
	return fmt.Sprintf("%s:%d", IDKeyPrefix, id)
}

func GetUsernameKey(username string) string {
	return fmt.Sprintf("%s:%s", UsernameKeyPrefix, username)
}

func GetAdminKey(id uint) string {
	return fmt.Sprintf("%s:%d", AdminKey, id)
}

func GetSysAdminKey() string {
	return SysAdminKey
}
