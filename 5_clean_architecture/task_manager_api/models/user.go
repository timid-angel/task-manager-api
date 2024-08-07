package models

/*
This is the definition of the user struct used for the authentication
and authorization aspects of the project. The email and user name will
be unique through entries. Along with the field names, the json labels
are provided to facilitate the binding process between the model itself
and the JSON format.
*/
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
