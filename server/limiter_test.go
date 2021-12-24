package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimiter_AllowRequestFrom(t *testing.T) {
	aUser := "test-user-1"
	anotherUser := "test-user-2"

	limiterWithZeroLimit := NewTokenBucketLimiter(0, 10*time.Millisecond)

	limiterWithQuota := NewTokenBucketLimiter(10, 10*time.Millisecond)

	limiterWithUsedQuota := NewTokenBucketLimiter(1, 10*time.Second)
	limiterWithUsedQuota.AllowRequestFrom(aUser)

	limiterWithUsedQuotaButTinyWindow := NewTokenBucketLimiter(1, 1*time.Nanosecond)
	limiterWithUsedQuotaButTinyWindow.AllowRequestFrom(aUser)

	tests := []struct {
		name    string
		limiter Limiter
		userid  string
		wantRes bool
	}{
		{
			name:    "zero limit",
			limiter: limiterWithZeroLimit,
			userid:  aUser,
			wantRes: false,
		},
		{
			name:    "with quota",
			limiter: limiterWithQuota,
			userid:  aUser,
			wantRes: true,
		},
		{
			name:    "used quota",
			limiter: limiterWithUsedQuota,
			userid:  aUser,
			wantRes: false,
		},

		{
			name:    "used quota but another user",
			limiter: limiterWithUsedQuota,
			userid:  anotherUser,
			wantRes: true,
		},
		{
			name:    "used quota but tiny window",
			limiter: limiterWithUsedQuotaButTinyWindow,
			userid:  aUser,
			wantRes: true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			res := test.limiter.AllowRequestFrom(test.userid)
			assert.Equal(t, test.wantRes, res)
		})
	}
}
