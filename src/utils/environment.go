package utils

import "os/user"

func GetUserProfile() (*user.User, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	return usr, nil
}
