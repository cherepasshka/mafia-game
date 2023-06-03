package pdf

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"

	"soa.mafia-game/game-server/domain/models/user"
)

func WriteUser(pdf *gofpdf.Fpdf, user user.User) (*gofpdf.Fpdf, error) {
	if pdf == nil {
		pdf = gofpdf.New("P", "mm", "A4", "")
	}
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Write(16, "Profile\n")
	pdf.SetFont("Arial", "", 16)
	pdf.Write(16, fmt.Sprintf("Login: %s\nGender: %s\nEmail: %s\n", user.Login, user.Gender, user.Email))

	pdf.SetFont("Arial", "B", 16)
	pdf.Write(16, "Statistics\n")
	pdf.SetFont("Arial", "", 16)
	pdf.Write(16, fmt.Sprintf("Total amount of sessions: %v\nAmount of victories: %v\n", user.SessionsCnt, user.VictoriesCnt))
	pdf.Write(16, fmt.Sprintf("Amount of defeats: %v\nTime spent in game: %v", user.SessionsCnt-user.VictoriesCnt, user.GetTotalGameTime()))
	return pdf, nil
}
