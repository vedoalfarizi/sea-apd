package user_role

type UserRole int

const (
	CUSTOMER = iota
	MERCHANT
	ADMIN
	OTHER
)

var UserRoleList = []string{
	"customer",
	"merchant",
	"admin",
	"other",
}

func ToString(ts UserRole) string {
	if ts < CUSTOMER || ts > OTHER {
		return ""
	}
	return UserRoleList[ts]
}

func ParseToEnum(src string) UserRole {
	userRoleMap := map[string]UserRole{
		"customer": CUSTOMER,
		"merchant": MERCHANT,
		"admin":    ADMIN,
		"other":    OTHER,
	}
	if val, exist := userRoleMap[src]; exist {
		return val
	}
	return userRoleMap["other"]
}

