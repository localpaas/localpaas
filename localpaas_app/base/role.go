package base

var (
	//nolint
	mapRoleValue = map[UserRole]int{
		UserRoleAdmin:  100,
		UserRoleMember: 10,
	}
)

func RoleCmp(r1, r2 UserRole) int {
	return mapRoleValue[r1] - mapRoleValue[r2]
}
