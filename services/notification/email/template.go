package email

import "fmt"

// EmailTemplate defines subject and body for a template
type EmailTemplate interface {
    Subject() string
    Body() string
}

// OnboardingEmail template
type OnboardingEmail struct {
    UserID uint
    Email  string
}

func (e OnboardingEmail) Subject() string {
    return "Welcome to OpenFinstack!"
}

func (e OnboardingEmail) Body() string {
    return fmt.Sprintf(`<h1>Hi!</h1><p>Welcome aboard. Your user ID is %d.</p>`, e.UserID)
}

// Add other types like ResetPasswordEmail, InviteEmail, etc. similarly
