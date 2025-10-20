package base

type UserRole string

const (
	UserRoleOwner  = UserRole("owner")
	UserRoleAdmin  = UserRole("admin")
	UserRoleMember = UserRole("member")
)

var (
	AllUserRoles = []UserRole{UserRoleOwner, UserRoleAdmin, UserRoleMember}
)

type UserStatus string

const (
	UserStatusInvited  = UserStatus("invited")
	UserStatusActive   = UserStatus("active")
	UserStatusDisabled = UserStatus("disabled")
)

var (
	AllUserStatuses = []UserStatus{UserStatusInvited, UserStatusActive, UserStatusDisabled}
)

type UserSecurityOption string

const (
	UserSecurityEnforceSSO   = UserSecurityOption("enforce-sso")
	UserSecurityPassword2FA  = UserSecurityOption("password-2fa")
	UserSecurityPasswordOnly = UserSecurityOption("password-only")
)

var (
	AllUserSecurityOptions = []UserSecurityOption{UserSecurityEnforceSSO, UserSecurityPassword2FA,
		UserSecurityPasswordOnly}
)
