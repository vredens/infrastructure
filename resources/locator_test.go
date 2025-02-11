package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	var name string
	var role string
	var err error

	name, role, err = parse("arn://storage/postgres/mydb", "storage", "postgres")
	assert.Equal(t, "mydb", name)
	assert.Equal(t, "", role)
	assert.Nil(t, err)

	name, role, err = parse("arn://storage/postgres/mydb/myrole", "storage", "postgres")
	assert.Equal(t, "mydb", name)
	assert.Equal(t, "myrole", role)
	assert.Nil(t, err)

	name, role, err = parse("arn://storage/postgres/mydb/myrole")
	assert.Equal(t, "", name)
	assert.Equal(t, "", role)
	assert.NotNil(t, err)

	name, role, err = parse("arn://databases", "storage", "postgres")
	assert.Equal(t, "", name)
	assert.Equal(t, "", role)
	assert.NotNil(t, err)
}

func TestParams(t *testing.T) {
	var params = Params{
		"key1": "value1",
		"key2": "value2",
		"key4": 4,
		"key5": 5.0,
		"key6": true,
		"key7": false,
	}

	assert.Equal(t, "value1", params.String("key1"))
	assert.Equal(t, "value2", params.String("key2"))
	assert.Equal(t, "", params.String("key3"))
	assert.Equal(t, 4, params.Int("key4"))
	assert.Equal(t, 0, params.Int("key5"))
	assert.Equal(t, true, params.Bool("key6"))
	assert.Equal(t, false, params.Bool("key7"))
	assert.Equal(t, 5.0, params.Float64("key5"))
}
