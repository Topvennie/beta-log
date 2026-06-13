// Package api contains all admin api routes
package api

import (
	"errors"
	"fmt"

	"github.com/Topvennie/beta-log/internal/server/dto"
	"github.com/Topvennie/beta-log/internal/server/service"
	"github.com/Topvennie/beta-log/pkg/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/shareed2k/goth_fiber/v2"
	"go.uber.org/zap"
)

type auth struct {
	router fiber.Router
	user   service.User

	redirectURL string
}

func newAuth(router fiber.Router, service service.Service) (*auth, error) {
	openidConnect, err := openidConnect.New(
		config.GetString("auth.oidc.client_id"),
		config.GetString("auth.oidc.client_secret"),
		config.GetString("auth.oidc.callback_url"),
		config.GetString("auth.oidc.discovery_url"),
		openidConnect.NameClaim,
		openidConnect.ProfileClaim,
	)
	if err != nil {
		return nil, fmt.Errorf("initialize oidc %w", err)
	}

	goth.UseProviders(openidConnect)

	api := &auth{
		router:      router.Group("/auth"),
		user:        *service.NewUser(),
		redirectURL: config.GetDefaultString("auth.redirect_url", "/"),
	}

	api.routes()

	return api, nil
}

func (r *auth) routes() {
	r.router.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	r.router.Get("/callback/:provider", r.loginCallbackHandler)
	r.router.Post("/logout", r.logoutHandler)
}

func (r *auth) loginCallbackHandler(c fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c, goth_fiber.CompleteUserAuthOptions{ShouldLogout: false})
	if err != nil {
		return fmt.Errorf("complete user auth %w", err)
	}

	dtoUser, err := r.user.GetByUID(c.RequestCtx(), user.UserID)
	if err != nil && !errors.Is(err, fiber.ErrNotFound) {
		return fmt.Errorf("get user by uid %w", err)
	}

	if errors.Is(err, fiber.ErrNotFound) {
		dtoUser = dto.User{
			UID:  user.UserID,
			Name: user.Name,
		}
		dtoUser, err = r.user.Save(c.RequestCtx(), dtoUser)
		if err != nil {
			return fmt.Errorf("save new user %w", err)
		}
	}

	sess := session.FromContext(c)
	if sess == nil {
		return errors.New("tried to authenticate user without active session")
	}

	sess.Set("id", dtoUser.ID)
	sess.Set("uid", dtoUser.UID)

	return c.Redirect().To(r.redirectURL)
}

func (r *auth) logoutHandler(c fiber.Ctx) error {
	sess := session.FromContext(c)
	if sess == nil {
		return errors.New("no session found")
	}

	if err := sess.Destroy(); err != nil {
		return fmt.Errorf("destroy session %w", err)
	}

	if err := goth_fiber.Logout(c); err != nil {
		zap.S().Warnf("logout failed: %v", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
