package auth

import (
	"testing"
	"time"
)

func newTestManager() TokenManager {
	return NewJWTManager("0123456789abcdef0123456789abcdef", time.Minute, time.Hour)
}

func TestAccessAndRefreshTokensAreNotInterchangeable(t *testing.T) {
	m := newTestManager()

	access, err := m.GenerateAccess("user-1")
	if err != nil {
		t.Fatal(err)
	}
	refresh, err := m.GenerateRefresh("user-1")
	if err != nil {
		t.Fatal(err)
	}

	if c, err := m.ParseAccess(access.Token); err != nil || c.UserID != "user-1" || c.TokenID != access.ID {
		t.Fatalf("ParseAccess(access) = %+v, %v", c, err)
	}
	if c, err := m.ParseRefresh(refresh.Token); err != nil || c.TokenID != refresh.ID {
		t.Fatalf("ParseRefresh(refresh) = %+v, %v", c, err)
	}
	if _, err := m.ParseAccess(refresh.Token); err == nil {
		t.Fatal("refresh token accepted as access token")
	}
	if _, err := m.ParseRefresh(access.Token); err == nil {
		t.Fatal("access token accepted as refresh token")
	}
}

func TestParseRejectsExpiredAndForeignTokens(t *testing.T) {
	expired := NewJWTManager("0123456789abcdef0123456789abcdef", -time.Minute, -time.Minute)
	tok, _ := expired.GenerateAccess("user-1")
	if _, err := newTestManager().ParseAccess(tok.Token); err == nil {
		t.Fatal("expired token accepted")
	}

	other := NewJWTManager("ffffffffffffffffffffffffffffffff", time.Minute, time.Hour)
	tok, _ = other.GenerateAccess("user-1")
	if _, err := newTestManager().ParseAccess(tok.Token); err == nil {
		t.Fatal("token signed with another secret accepted")
	}
}
