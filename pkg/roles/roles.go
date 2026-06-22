package roles

type Role string

const (
	RoleTeacher Role = "teacher"
	RoleStudent Role = "student"
	RoleAdmin   Role = "admin"
	RoleUnknown Role = ""
)

func (r Role) IsTeacher() bool { return r == RoleTeacher }
func (r Role) IsStudent() bool { return r == RoleStudent }
func (r Role) IsAdmin() bool   { return r == RoleAdmin }
func (r Role) IsUnknown() bool { return r == RoleUnknown }

func IsRole(role string) bool {
	r := Role(role)
	if r.IsTeacher() ||
		r.IsStudent() ||
		r.IsAdmin() ||
		r.IsUnknown() {
		return true
	}
	return false
}

func GetRole(role string) Role {
	r := Role(role)
	if r.IsAdmin() ||
		r.IsStudent() ||
		r.IsTeacher() {
		return r
	}
	return RoleUnknown
}

func (r Role) String() string {
	return string(r)
}

func (r *Role) PtrString() *string {
	if r == nil {
		return nil
	}
	str := string(*r)
	return &str
}
