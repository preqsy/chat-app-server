// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AuthUser struct {
	ID        int32  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type AuthUserCreate struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type AuthUserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUserResponse struct {
	AuthUser *AuthUser `json:"authUser"`
	Token    string    `json:"token"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Mutation struct {
}

type Query struct {
}

type Subscription struct {
}

type Time struct {
	UnixTime  int32  `json:"unixTime"`
	TimeStamp string `json:"timeStamp"`
}

type UserEmail struct {
	Email string `json:"email"`
}
