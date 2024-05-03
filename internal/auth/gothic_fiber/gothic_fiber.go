// This package is adapted from the `github.com/markbates/goth/gothic`
// to work with the Fiber web framework.
package gothic_fiber

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/markbates/goth"
)

const SessionName = "_gothic_session"

// Session can/should be set by applications using gothic. The default is a cookie store.
var (
	SessionStore  *session.Store
	ErrSessionNil = errors.New("goth/gothic: no SESSION_SECRET environment variable is set. The default cookie store is not available and any calls will fail. Ignore this warning if you are using a different store")
)

func init() {
	// optional config
	config := session.Config{
		KeyLookup:      fmt.Sprintf("cookie:%s", SessionName),
		Expiration:     10 * time.Minute, // as short as possible for security
		CookieSameSite: "Lax",
		CookieHTTPOnly: true,
		// CookieSecure:   true, // TODO: HTTPS only, dev mode should be false
	}

	SessionStore = session.New(config)
}

// BeginAuthHandler is a convenience handler for starting the authentication process.
// It expects to be able to get the name of the provider from the query parameters
// as either "provider" or ":provider".

// BeginAuthHandler will redirect the user to the appropriate authentication end-point
// for the requested provider.

// See https://github.com/markbates/goth/examples/main.go to see this in action.
func BeginAuthHandler(ctx *fiber.Ctx) error {
	url, err := GetAuthURL(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
func SetState(ctx *fiber.Ctx) string {
	state := ctx.Query("state")
	if len(state) > 0 {
		return state
	}

	// If a state query param is not passed in, generate a random
	// base64-encoded nonce so that the state on the auth URL
	// is unguessable, preventing CSRF attacks, as described in
	//
	// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic("gothic: source of randomness unavailable: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
func GetState(ctx *fiber.Ctx) string {
	return ctx.Query("state")
}

// GetAuthURL starts the authentication process with the requested provided.
// It will return a URL that should be used to send users to.

// It expects to be able to get the name of the provider from the query parameters
// as either "provider" or ":provider".

// I would recommend using the BeginAuthHandler instead of doing all of these steps
// yourself, but that's entirely up to you.
func GetAuthURL(ctx *fiber.Ctx) (string, error) {
	if SessionStore == nil {
		return "", ErrSessionNil
	}

	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}

	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = StoreInSession(providerName, sess.Marshal(), ctx)
	if err != nil {
		return "", err
	}

	return url, err
}

// Options that affect how CompleteUserAuth works.
type CompleteUserAuthOptions struct {
	// True if CompleteUserAuth should automatically end the user's session.
	//
	// Defaults to True.
	ShouldLogout bool
}

// CompleteUserAuth does what it says on the tin. It completes the authentication
// process and fetches all of the basic information about the user from the provider.

// It expects to be able to get the name of the provider from the query parameters
// as either "provider" or ":provider".

// This method automatically ends the session. You can prevent this behavior by
// passing in options. Please note that any options provided in addition to the
// first will be ignored.

// See https://github.com/markbates/goth/examples/main.go to see this in action.
func CompleteUserAuth(ctx *fiber.Ctx, options ...CompleteUserAuthOptions) (goth.User, error) {
	if SessionStore == nil {
		return goth.User{}, ErrSessionNil
	}

	providerName, err := GetProviderName(ctx)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}

	value, err := GetFromSession(providerName, ctx)
	if err != nil {
		return goth.User{}, err
	}

	shouldLogout := true
	if len(options) > 0 && !options[0].ShouldLogout {
		shouldLogout = false
	}

	if shouldLogout {
		defer Logout(ctx)
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	err = validateState(ctx, sess)
	if err != nil {
		return goth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	reqURL, err := url.Parse(ctx.Request().URI().String())
	if err != nil {
		return goth.User{}, err
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, reqURL.Query())
	if err != nil {
		return goth.User{}, err
	}

	err = StoreInSession(providerName, sess.Marshal(), ctx)

	if err != nil {
		return goth.User{}, err
	}

	gu, err := provider.FetchUser(sess)
	return gu, err
}

// validateState ensures that the state token param from the original
// AuthURL matches the one included in the current (callback) request.
func validateState(ctx *fiber.Ctx, sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	originalState := authURL.Query().Get("state")
	if originalState != "" && (originalState != ctx.Query("state")) {
		return errors.New("state token mismatch")
	}
	return nil
}

// Logout invalidates a user session.
func Logout(ctx *fiber.Ctx) error {
	session, err := SessionStore.Get(ctx)
	if err != nil {
		return err
	}

	if err := session.Destroy(); err != nil {
		return err
	}

	return nil
}

// GetProviderName is a function used to get the name of a provider
// for a given request. By default, this provider is fetched from
// the URL query string. If you provide it in a different way,
// assign your own function to this variable that returns the provider
// name for your request.
func GetProviderName(ctx *fiber.Ctx) (string, error) {
	// try to get it from the url param "provider"
	if p := ctx.Query("provider"); p != "" {
		return p, nil
	}

	// try to get it from the url param ":provider"
	if p := ctx.Params("provider"); p != "" {
		return p, nil
	}

	//  try to get it from the Fasthttp context's value of "provider" key
	if p := ctx.Get("provider", ""); p != "" {
		return p, nil
	}

	// As a fallback, loop over the used providers, if we already have a valid session for any provider (ie. user has already begun authentication with a provider), then return that provider name
	providers := goth.GetProviders()
	session, err := SessionStore.Get(ctx)
	if err != nil {
		return "", err
		// or panic?
	}

	for _, provider := range providers {
		p := provider.Name()
		value := session.Get(p)
		if _, ok := value.(string); ok {
			return p, nil
		}
	}

	// if not found then return an empty string with the corresponding error
	return "", errors.New("you must select a provider")
}

// StoreInSession stores a specified key/value pair in the session.
func StoreInSession(key string, value string, ctx *fiber.Ctx) error {
	session, err := SessionStore.Get(ctx)
	if err != nil {
		return err
	}

	if err := updateSessionValue(session, key, value); err != nil {
		return err
	}

	// saved here
	session.Save()
	return nil
}

// GetFromSession retrieves a previously-stored value from the session.
// If no value has previously been stored at the specified key, it will return an error.
func GetFromSession(key string, ctx *fiber.Ctx) (string, error) {
	session, err := SessionStore.Get(ctx)
	if err != nil {
		return "", err
	}

	value, err := getSessionValue(session, key)
	if err != nil {
		return "", errors.New("could not find a matching session for this request")
	}

	return value, nil
}

func getSessionValue(store *session.Session, key string) (string, error) {
	value := store.Get(key)
	if value == nil {
		return "", errors.New("could not find a matching session for this request")
	}

	rdata := strings.NewReader(value.(string))
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return "", err
	}
	s, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

func updateSessionValue(session *session.Session, key, value string) error {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(value)); err != nil {
		return err
	}
	if err := gz.Flush(); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	session.Set(key, b.String())

	return nil
}
