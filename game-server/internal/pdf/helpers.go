package pdf

import (
	"fmt"
	"os"

	"github.com/jung-kurt/gofpdf"

	"soa.mafia-game/game-server/domain/models/user"
)

func WriteUser(pdf *gofpdf.Fpdf, user user.User) (*gofpdf.Fpdf, error) {
	if pdf == nil {
		pdf = gofpdf.New("P", "mm", "A4", "")
	}
	pdf.AddPage()

	pdf = writeCommonInfo(pdf, user)
	pdf = writeStatistics(pdf, user)
	if len(user.ImageName) == 0 {
		return pdf, nil
	}
	var err error
	pdf, err = addPicture(pdf, user)
	return pdf, err
}

func writeCommonInfo(pdf *gofpdf.Fpdf, user user.User) *gofpdf.Fpdf {
	pdf.SetFont("Arial", "B", 16)
	pdf.Write(16, "Profile\n")
	pdf.SetFont("Arial", "", 16)
	pdf.Write(16, fmt.Sprintf("Login: %s\nGender: %s\nEmail: %s\n", user.Login, user.Gender, user.Email))
	return pdf
}

func writeStatistics(pdf *gofpdf.Fpdf, user user.User) *gofpdf.Fpdf {
	pdf.SetFont("Arial", "B", 16)
	pdf.Write(16, "Statistics\n")
	pdf.SetFont("Arial", "", 16)
	pdf.Write(16, fmt.Sprintf("Total amount of sessions: %v\nAmount of victories: %v\n", user.SessionsCnt, user.VictoriesCnt))
	pdf.Write(16, fmt.Sprintf("Amount of defeats: %v\nTime spent in game: %v", user.SessionsCnt-user.VictoriesCnt, user.GetTotalGameTime()))
	return pdf
}

func addPicture(pdf *gofpdf.Fpdf, user user.User) (*gofpdf.Fpdf, error) {
	imgName := fmt.Sprintf("images/%s", user.ImageName)
	file, err := os.Open(imgName)
	if err != nil {
		return pdf, err
	}
	defer file.Close()
	imageOptions := gofpdf.ImageOptions{ImageType: "jpg", ReadDpi: true}
	pdf.RegisterImageOptionsReader(imgName, imageOptions, file)
	pdf.ImageOptions(imgName, 100, 10, 50, 50, false, imageOptions, 0, "")
	return pdf, nil
}
