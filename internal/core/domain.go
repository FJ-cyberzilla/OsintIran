// internal/core/domain.go
package core

type User struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// internal/services/user_service.go
type UserService struct {
    repo    repositories.UserRepository
    cache   cache.Cache
    eventBus event.Bus
}

func (s *UserService) CreateUser(ctx context.Context, user *core.User) error {
    // Business logic with events
    s.eventBus.Publish(UserCreatedEvent{User: user})
}
