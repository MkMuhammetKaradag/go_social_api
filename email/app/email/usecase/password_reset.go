package usecase

import (
	"path/filepath"
	"socialmedia/email/internal/domain"
)

type passwordResetUsecase struct {
	mailer      domain.Mailer
	templateDir string
}

func NewPasswordResetService(mailer domain.Mailer, templateDir string) PasswordResetUseCase {
	return &passwordResetUsecase{
		mailer:      mailer,
		templateDir: templateDir,
	}
}

func (uc *passwordResetUsecase) SendPasswordResetEmail(email, resetLink, userName, templateName string) error {
	body, err := RenderTemplate(
		filepath.Join(uc.templateDir, templateName),
		struct {
			ResetLink string
			UserName  string
		}{
			ResetLink: resetLink,
			UserName:  userName,
		},
	)
	if err != nil {
		return err
	}

	return uc.mailer.Send(email, "Password Reset", body)
}

// func (uc *passwordResetUsecase) renderTemplate(templateFile string, data interface{}) (string, error) {
// 	tmpl, err := template.ParseFiles(templateFile)
// 	if err != nil {
// 		return "", err
// 	}

// 	var buf bytes.Buffer
// 	if err := tmpl.Execute(&buf, data); err != nil {
// 		return "", err
// 	}

// 	return buf.String(), nil
// }
