package roles

type Role string

const (
	RoleTeacher Role = "teacher"
	RoleStudent Role = "student"
	RoleAdmin   Role = "admin"
	RoleUser    Role = "user"
)

func (r Role) IsTeacher() bool { return r == RoleTeacher }
func (r Role) IsStudent() bool { return r == RoleStudent }
func (r Role) IsAdmin() bool   { return r == RoleAdmin }
func (r Role) IsUser() bool    { return r == RoleUser }

func (r Role) IsRole() bool {
	if r.IsTeacher() ||
		r.IsStudent() ||
		r.IsAdmin() ||
		r.IsUser() {
		return true
	}
	return false
}
