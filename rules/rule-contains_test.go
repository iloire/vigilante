package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidContainsRule(t *testing.T) {
	assert := assert.New(t)

	rule := Contains{"hello world"}
	body := "hello moon"
	assert.False(rule.Match(body, 200).Success, "Should be unsuccessful")
	assert.Equal(rule.Match(body, 200).Error, "Response doesn't contain hello world", "should fail")
}

func TestValidContainsRule(t *testing.T) {
	assert := assert.New(t)

	rule := Contains{"hello world"}
	body := "this is my own hello world, yeah!"
	assert.True(rule.Match(body, 200).Success, "Should be successful")
	assert.Equal(rule.Match(body, 200).Error, "", "should be ok")
}
