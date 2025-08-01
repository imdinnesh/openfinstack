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

type KYCStatusEmail struct {
    UserID uint
    Email  string
    Status string // e.g. "Approved", "Rejected", "Pending"
    Reason string // Optional: reason for rejection or additional info
}

func (e KYCStatusEmail) Subject() string {
    return fmt.Sprintf("Your KYC Status: %s", e.Status)
}

func (e KYCStatusEmail) Body() string {
    body := fmt.Sprintf(`
        <div style="font-family: Arial, sans-serif; line-height: 1.6;">
            <h2>KYC Status Update 📄</h2>
            <p>Hi there,</p>
            <p>Your KYC process has been <strong>%s</strong>.</p>
            <p><strong>User ID:</strong> %d</p>`, e.Status, e.UserID)

    if e.Status == "Rejected" && e.Reason != "" {
        body += fmt.Sprintf(`
            <p><strong>Reason:</strong> %s</p>
            <p>Please review the above and resubmit your KYC information if needed.</p>`, e.Reason)
    }

    body += `
            <p>If you have any questions, feel free to reach out to our support team.</p>
            <p>Thank you,<br/>The OpenFinstack Team</p>
        </div>
    `

    return body
}
