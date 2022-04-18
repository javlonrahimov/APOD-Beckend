package mock

import "javlonrahimov/apod/internal/data"

var permissions = make(map[int64][]string)

type permissionModelMock struct{}

func NewPermissionsMock() data.PermissonService {
	return &permissionModelMock{}
}

func (m permissionModelMock) GetAllForUser(userId int64) (data.Permissions, error) {

	userPermissions, ok := permissions[userId]
	if ok {
		return userPermissions, nil
	}

	return nil, nil
}

func (m permissionModelMock) AddForUser(userId int64, codes ...string) error {
	permissions[userId] = codes
	return nil
}
