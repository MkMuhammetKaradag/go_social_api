package usecase

type ActivationUseCase interface {
	SendActivationEmail(email, resetLink, userName, templateName string) error
}
type PasswordResetUseCase interface {
	SendPasswordResetEmail(email, resetLink, userName, templateName string) error
}
