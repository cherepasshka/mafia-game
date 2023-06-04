package user

import "time"

type GenderType string

const (
	Undefined GenderType = "undefined"
	Male      GenderType = "male"
	Female    GenderType = "female"
)

type Profile struct {
	Login     string     `json:"login"`
	Email     string     `json:"email"`
	Gender    GenderType `json:"gender"`
	ImageName string     `json:"image"`
}

type User struct {
	Profile
	SessionsCnt      int
	VictoriesCnt     int
	TotalGameTime    time.Duration
	LastSessionEnter time.Time
	ActiveSession    bool
}

func New(login string) User {
	return User{
		Profile: Profile{
			Login:  login,
			Gender: Undefined,
		},
		SessionsCnt:   0,
		VictoriesCnt:  0,
		TotalGameTime: 0,
		ActiveSession: false,
	}
}

func (user User) GetTotalGameTime() time.Duration {
	var delta time.Duration = 0
	if user.ActiveSession {
		delta = time.Now().Sub(user.LastSessionEnter)
	}
	return user.TotalGameTime + delta
}
