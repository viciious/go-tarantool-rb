package tnt

import (
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	assert := assert.New(t)

	tarantoolConfig := `
    s = box.schema.space.create('tester', {id = 42})
    s:create_index('primary', {
        type = 'hash',
        parts = {1, 'NUM'}
    })
    t = s:insert({1})
    t = s:insert({2, 'Music'})
    t = s:insert({3, 'Length', 93})
    `

	box, err := NewBox(tarantoolConfig, nil)
	if !assert.NoError(err) {
		return
	}
	defer box.Close()

	// unkwnown user
	conn, err := box.Connect(nil)
	assert.NoError(err)
	assert.NotNil(conn)

	data, err := conn.Execute(&Select{
		Space: 42,
		Key:   3,
	})

	if assert.NoError(err) {
		assert.Equal([]interface{}{
			[]interface{}{
				uint32(0x3),
				"Length",
				uint32(0x5d),
			},
		}, data)
	}

	// select with string space
	data, err = conn.Execute(&Select{
		Space: "tester",
		Key:   3,
	})

	if assert.NoError(err) {
		assert.Equal([]interface{}{
			[]interface{}{
				uint32(0x3),
				"Length",
				uint32(0x5d),
			},
		}, data)
	}
}

func TestSelectSpaceSchema(t *testing.T) {
	assert := assert.New(t)

	tarantoolConfig := `
    s = box.schema.space.create('tester', {id = 42})
    s:create_index('primary', {
        type = 'hash',
        parts = {1, 'NUM'}
    })
    t = s:insert({1})
    t = s:insert({2, 'Music'})
    t = s:insert({3, 'Length', 93})

    box.schema.user.create("tester", {password = "12345678"})
    box.schema.user.grant('tester', 'read', 'space', 'tester')
    `
	box, err := NewBox(tarantoolConfig, nil)
	if !assert.NoError(err) {
		return
	}
	defer box.Close()

	conn, err := box.Connect(&Options{
		User:     "tester",
		Password: "12345678",
	})
	assert.NoError(err)
	assert.NotNil(conn)

	data, err := conn.Execute(&Select{
		Space:    ViewSpace,
		Key:      0,
		Iterator: IterGt,
		Limit:    1000000,
	})
	pp.Println(data, err)
}

func BenchmarkSelectPack(b *testing.B) {
	d, _ := newPackData(42)

	for i := 0; i < b.N; i += 1 {
		(&Select{Key: 3}).Pack(0, d)
	}
}
