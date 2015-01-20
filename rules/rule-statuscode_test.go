package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidStatusCodeRule(t *testing.T) {
	assert := assert.New(t)

	rule := StatusCode{200}
	assert.False(rule.Match("dummy content", 302).Success, "should be unsuccessful")
	assert.Equal(rule.Match("dummy content", 302).Error, "Invalid status code", "should be unsuccessful")
}

func TestValidStatusCodeRule(t *testing.T) {
	assert := assert.New(t)

	rule := StatusCode{200}
	assert.True(rule.Match("dummy content", 200).Success, "should be successful")
	assert.Equal(rule.Match("dummy content", 200).Error, "", "should be succeessful")
}
