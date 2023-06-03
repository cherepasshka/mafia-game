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
	pdf.Write(16, fmt.Sprintf("Login: %s\nGender: %s\nEmail: %s", user.Login, user.Gender, user.Email))
	return pdf, nil
}
