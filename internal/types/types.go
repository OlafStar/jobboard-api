package types

type RegisterCompany struct {
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"password"`
	Country string `json:"country"`
	City    string `json:"city"`
	Address string `json:"address"`
}

type JWTCompany struct {
	Email string `json:"email"`
	Name string `json:"name"`
	PasswordHash string `json:"password"`
	Country string `json:"country"`
	City    string `json:"city"`
	Address string `json:"address"`
}

type LoginCompany struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type Company struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Website string `json:"website"`
	Logo    string `json:"logo"`
}

type ContactDetails struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// type Location struct {
// 	Country string `json:"country"`
// 	City    string `json:"city"`
// 	Address string `json:"address"`
// }