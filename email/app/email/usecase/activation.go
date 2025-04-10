package usecase

import (
	"path/filepath"
	"socialmedia/email/internal/domain"
)

type activationUsecase struct {
	mailer      domain.Mailer
	templateDir string
}

func NewActivationService(mailer domain.Mailer, templateDir string) ActivationUseCase {
	return &activationUsecase{
		mailer:      mailer,
		templateDir: templateDir,
	}
}
func (uc *activationUsecase) SendActivationEmail(email, code, userName, templateName string) error {
	body, err := RenderTemplate(
		filepath.Join(uc.templateDir, templateName),
		struct {
			Code     string
			UserName string
		}{
			Code:     code,
			UserName: userName,
		},
	)
	if err != nil {
		return err
	}

	return uc.mailer.Send(email, "Account Activation", body)
}

// func (uc *activationUsecase) RenderTemplate(templateFile string, data interface{}) (string, error) {
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
