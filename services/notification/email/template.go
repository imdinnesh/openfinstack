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
    return fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; line-height: 1.6;">
            <h2>Welcome to OpenFinstack 🎉</h2>
            <p>Hi there,</p>
            <p>We're thrilled to have you on board! Your journey into the future of finance starts now.</p>
            <p><strong>Your User ID:</strong> %d</p>
            <p>If you have any questions or need support, feel free to reach out to us anytime.</p>
            <p>Happy building,<br/>The OpenFinstack Team</p>
        </div>
    `, e.UserID)
}

// Add other types like ResetPasswordEmail, InviteEmail, etc. similarly
