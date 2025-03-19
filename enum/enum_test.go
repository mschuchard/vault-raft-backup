package enum

import "testing"

func TestPlatformFromString(test *testing.T) {
	platform, err := Platform("").From("local")
	if err != nil {
		test.Error(err)
	}
	if platform != LOCAL {
		test.Error("platform did not type convert correctly")
		test.Errorf("expected: LOCAL, actual: %s", platform)
	}

	if _, err = Platform("").From("foo"); err == nil || err.Error() != "invalid platform enum" {
		test.Error("platform type conversion did not error expectedly")
		test.Errorf("expected: invalid platform enum, actual: %s", err)
	}
}

func TestAuthEngineFromString(test *testing.T) {
	authEngine, err := AuthEngine("").From("token")
	if err != nil {
		test.Error(err)
	}
	if authEngine != VaultToken {
		test.Error("authengine did not type convert correctly")
		test.Errorf("expected: token, actual: %s", authEngine)
	}

	if _, err = AuthEngine("").From("foo"); err == nil || err.Error() != "invalid authengine enum" {
		test.Error("authengine type conversion did not error expectedly")
		test.Errorf("expected: invalid authengine enum, actual: %s", err)
	}
}
