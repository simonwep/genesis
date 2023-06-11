package main

import "golang.org/x/exp/slices"

func ValidateUserName(name string) bool {
	users := Env().AppAllowedUsers
	return len(users) == 0 || slices.Contains(users, name)
}
