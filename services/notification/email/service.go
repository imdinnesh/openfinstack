package email

import "log"

type Service struct {
    sender Sender
}

func NewService(sender Sender) *Service {
    return &Service{sender: sender}
}

func (s *Service) Send(to string, template EmailTemplate) error {
    log.Printf("[EmailService] Sending '%s' to %s", template.Subject(), to)
    return s.sender.Send(to, template.Subject(), template.Body())
}
