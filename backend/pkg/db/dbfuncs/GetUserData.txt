package dbfuncs

type User struct {
	id        string
	Email     string
	FirstName string
	LastName  string
	userName  string
}

func GetUserData(id string) (User, error) {
	user := User{}
	err := database.QueryRow("SELECT Id, Password, Nickname,  Profile,  FROM Users WHERE Id=?", id).Scan(&user.id)
	if err != nil {
		return User{}, err
	}

}
