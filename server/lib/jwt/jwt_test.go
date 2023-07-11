package jwt

import "testing"

func TestCreateToken(t *testing.T) {
	token, err := CreateToken()
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal("token is empty")
	}
	t.Log("token: ", token)
}
