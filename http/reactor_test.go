package http

func init() {
	Accept("/index", User{}, "Index")
	Run()
}

type User struct {
	Reactor
}

func (u *User) Index() {

}
