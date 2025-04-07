package enum

import "testing"

func TestPlatformFromNew(test *testing.T) {
	platform, err := Platform("local").New()
	if err != nil {
		test.Error(err)
	}
	if platform != LOCAL {
		test.Error("platform did not type convert correctly")
		test.Errorf("expected: LOCAL, actual: %s", platform)
	}

	if _, err = Platform("foo").New(); err == nil || err.Error() != "invalid platform enum" {
		test.Error("platform type conversion did not error expectedly")
		test.Errorf("expected: invalid platform enum, actual: %s", err)
	}
}

func TestAuthEngineNew(test *testing.T) {
	authEngine, err := AuthEngine("token").New()
	if err != nil {
		test.Error(err)
	}
	if authEngine != VaultToken {
		test.Error("authengine did not type convert correctly")
		test.Errorf("expected: token, actual: %s", authEngine)
	}

	if _, err = AuthEngine("foo").New(); err == nil || err.Error() != "invalid authengine enum" {
		test.Error("authengine type conversion did not error expectedly")
		test.Errorf("expected: invalid authengine enum, actual: %s", err)
	}
}
