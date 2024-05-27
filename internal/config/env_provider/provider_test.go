package env_provider

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Provider(t *testing.T) {
	provider := Provider("ATK_", map[string]string{
		"ATK_TRUSTED_DOMAINS":          "trusted_domains.$$",
		"ATK_TRUSTED_DOMAINS_$$":       "trusted_domains.$$",
		"ATK_SIMPLE_ARRAY":             "simple_array.$$", // without ATK_SIMPLE_ARRAY_$$ in this map
		"ATK_DB_PORT":                  "db.port",
		"ATK_DB_PREPARE_STMT":          "db.prepare_stmt",
		"ATK_DB_PREPARE__STMT":         "db.prepare_stmt",
		"ATK_LOGIN__TIMEOUT":           "login_timeout",
		"ATK_DEBUG":                    "debug",
		"ATK_ADMIN_USERS_$$_EMAILS":    "admin_users.$$.emails.$$",
		"ATK_ADMIN_USERS_$$_EMAILS_$$": "admin_users.$$.emails.$$",
		"ATK_AUTH_AUTH0_CLIENT_ID":     "auth.auth0.client_id",
	})

	provider.osEnvironFunc = func() []string {
		return []string{
			"ATK_TRUSTED_DOMAINS=1 2 3",
			"ATK_TRUSTED_DOMAINS_0=666",
			"ATK_SIMPLE_ARRAY_0=test0",
			"ATK_DB_PORT=3306",
			"ATK_DB_PREPARE_STMT=true",
			"ATK_DB_PREPARE__STMT=false",
			"ATK_LOGIN__TIMEOUT=10",
			"ATK_DEBUG=true",
			"ATK_ADMIN_USERS_0_EMAILS_1=2",
			"ATK_ADMIN_USERS_0_EMAILS=a b c",
			"ATK_ADMIN_USERS_0_EMAILS_2=2",
			"ATK_AUTH_AUTH0_CLIENT_ID=123456",
		}
	}

	conf, err := provider.Read()
	if err != nil {
		t.Error(err)
	}

	confJson, _ := json.Marshal(conf)

	assert.JSONEq(t, `{
		"admin_users": [
		  {
			"emails": [
			  "a",
			  "b",
			  "2"
			]
		  }
		],
		"auth": {
		  "auth0": {
			"client_id": "123456"
		  }
		},
		"db": {
		  "port": "3306",
		  "prepare_stmt": "false"
		},
		"debug": "true",
		"login_timeout": "10",
		"simple_array": [
		  "test0"
		],
		"trusted_domains": [
		  "666",
		  "2",
		  "3"
		]
	  }`, string(confJson))

	// format print json with tab indent
	b, _ := json.MarshalIndent(conf, "", "  ")
	t.Log(string(b))
}

func Test_ProviderDetailer(t *testing.T) {
	performTest := func(t *testing.T, rules map[string]string, envs []string, expected string) {
		provider := Provider("ATK_", rules)
		provider.osEnvironFunc = func() []string {
			return envs
		}

		conf, err := provider.Read()
		if err != nil {
			t.Error(err)
		}

		confJson, _ := json.Marshal(conf)

		assert.JSONEq(t, expected, string(confJson))
	}

	t.Run("simple array (space split set)", func(t *testing.T) {
		performTest(t, map[string]string{
			"ATK_ARRAY": "array.$$",
		}, []string{
			"ATK_ARRAY=1 2 3",
		}, `{
			"array": [ "1", "2", "3" ]
		}`)
	})

	t.Run("simple array (number set)", func(t *testing.T) {
		performTest(t, map[string]string{
			"ATK_ARRAY": "array.$$",
		}, []string{
			"ATK_ARRAY_0=1",
			"ATK_ARRAY_1=2",
			"ATK_ARRAY_2=3",
		}, `{
			"array": [ "1", "2", "3" ]
		}`)
	})

	t.Run("override simple array", func(t *testing.T) {
		performTest(t, map[string]string{
			"ATK_ARRAY": "array.$$",
		}, []string{
			"ATK_ARRAY_0=1 2 3",
			"ATK_ARRAY=4 5 6",
		}, `{
			"array": [ "4", "5", "6" ]
		}`)
	})

	t.Run("struct array with simple array (number set)", func(t *testing.T) {
		performTest(t, map[string]string{
			"ATK_ARRAY_$$_ITEMS": "array.$$.items.$$",
		}, []string{
			"ATK_ARRAY_0_ITEMS_0=a",
			"ATK_ARRAY_0_ITEMS_1=b",
			"ATK_ARRAY_0_ITEMS_2=c",
			"ATK_ARRAY_1_ITEMS_0=d",
		}, `{
			"array": [
				{
					"items": [ "a", "b", "c" ]
				},
				{
					"items": [ "d" ]
				}
			]
		}`)
	})

	t.Run("struct array with simple array (space split set)", func(t *testing.T) {
		performTest(t, map[string]string{
			"ATK_ARRAY_$$_ITEMS": "array.$$.items.$$",
		}, []string{
			"ATK_ARRAY_0_ITEMS=a b c",
			"ATK_ARRAY_1_ITEMS=d",
		}, `{
			"array": [
				{
					"items": [
						"a",
						"b",
						"c"
					]
				},
				{
					"items": [
						"d"
					]
				}
			]
		}`)
	})

	t.Run("struct array (number set)", func(t *testing.T) {
		performTest(t, map[string]string{
			"ATK_ARRAY_$$_ITEMS": "array.$$.items",
		}, []string{
			"ATK_ARRAY_0_ITEMS=a",
			"ATK_ARRAY_1_ITEMS=b",
			"ATK_ARRAY_2_ITEMS=c",
			"ATK_ARRAY_1_ITEMS=d",
		}, `{
			"array": [
				{
					"items": "a"
				},
				{
					"items": "d"
				},
				{
					"items": "c"
				}
			]
		}`)
	})

	t.Run("override struct array", func(t *testing.T) {
		performTest(t, map[string]string{
			"ATK_ARRAY_$$_ITEMS": "array.$$.items",
		}, []string{
			"ATK_ARRAY_0_ITEMS=a",
			"ATK_ARRAY_1_ITEMS=d",
			"ATK_ARRAY_0_ITEMS=a b c",
		}, `{
			"array": [
				{
					"items": "a b c"
				},
				{
					"items": "d"
				}
			]
		}`)
	})
}

func Test_handleEnvPathMap(t *testing.T) {
	assert.Equal(t, map[string]string{
		"ATK_ARRAY":    "array.$$",
		"ATK_ARRAY_$$": "array.$$",
	}, handleEnvPathMap(map[string]string{
		"ATK_ARRAY": "array.$$",
	}), "it should append _$$ if path suffix is .$$ which represents simple el array")

	assert.Equal(t, map[string]string{
		"ATK_ARRAY_$$_A": "array.$$.a",
	}, handleEnvPathMap(map[string]string{
		"ATK_ARRAY_$$_A": "array.$$.a",
	}), "it should not append _$$ if path without .$$ suffix (even if has .$$. in middle) which represents nested struct el array")

	assert.Equal(t, map[string]string{
		"ATK_ARRAY_$$_A":    "array.$$.a.$$",
		"ATK_ARRAY_$$_A_$$": "array.$$.a.$$",
	}, handleEnvPathMap(map[string]string{
		"ATK_ARRAY_$$_A": "array.$$.a.$$",
	}), "it should append _$$ if path suffix is .$$ and has .$$. in middle which represents also simple el array")

	assert.Equal(t, map[string]string{
		"ATK_ARRAY_$$_A":    "array.$$.a.$$",
		"ATK_ARRAY_$$_A_$$": "array.$$.a.$$",
	}, handleEnvPathMap(map[string]string{
		"ATK_ARRAY_$$_A_$$": "array.$$.a.$$",
	}), "it should append non-suffix _$$ if suffix is .$$ which represents also simple el array")
}

func Test_getKeyInEnvPathMap(t *testing.T) {
	key, numbers := getKeyInEnvPathMap("ATK_ARRAY_$$_A")
	assert.Equal(t, "ATK_ARRAY_$$_A", key)
	assert.Equal(t, []string{}, numbers)

	key, numbers = getKeyInEnvPathMap("ATK_ARRAY_1_A")
	assert.Equal(t, "ATK_ARRAY_$$_A", key)
	assert.Equal(t, []string{"1"}, numbers)

	key, numbers = getKeyInEnvPathMap("ATK_ARRAY_1_A_2_B")
	assert.Equal(t, "ATK_ARRAY_$$_A_$$_B", key)
	assert.Equal(t, []string{"1", "2"}, numbers)
}

func Test_recoverNumbersInPath(t *testing.T) {
	path := recoverNumbersInPath("array.$$.a.$$", []string{"2", "1"})
	assert.Equal(t, "array.2.a.1", path)

	path = recoverNumbersInPath("array.$$.a.$$", []string{"3"})
	assert.Equal(t, "array.3.a.$$", path)

	path = recoverNumbersInPath("array.$$.a.$$", []string{})
	assert.Equal(t, "array.$$.a.$$", path)
}

func Test_getSimpleElemArrayPaths(t *testing.T) {
	envPathMap := map[string]string{
		"ATK_ARRAY":    "array.$$",
		"ATK_ARRAY_$$": "array.$$",
		"ATK_ARRAY_A":  "array.a",
		"ATK_ARRAY_B":  "array.b",
		"ATK_ARRAY_C":  "array.c",
	}

	paths := getSimpleElemArrayPaths(envPathMap)
	assert.Equal(t, []string{"array.$$"}, paths)
}
