package mails

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mailData"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Specified header should be filtered out.
*/
func TestRepliesFilter_IsAutomaticResponse_AutoRespondHeaderSet(t *testing.T) {
	// Filter only looking at headers, so the rest can be ignored:
	testMail := mailData.Mail{Headers: []string{"x-autorespond"}}

	testee := RepliesFilter{}
	assert.True(t, testee.IsAutomaticResponse(testMail), "Mail should be filtered")
}

/*
	Specified header should be filtered out.
*/
func TestRepliesFilter_IsAutomaticResponse_PrecendenceSetToAutoReply(t *testing.T) {
	// Filter only looking at headers, so the rest can be ignored:
	testMail := mailData.Mail{Headers: []string{"x-precedence:auto_reply"}}

	testee := RepliesFilter{}
	assert.True(t, testee.IsAutomaticResponse(testMail), "Mail should be filtered")
}

/*
	Specified header should be filtered out.
*/
func TestRepliesFilter_IsAutomaticResponse_PrecendenceSetToBulk(t *testing.T) {
	// Filter only looking at headers, so the rest can be ignored:
	testMail := mailData.Mail{Headers: []string{"precedence:bulk"}}

	testee := RepliesFilter{}
	assert.True(t, testee.IsAutomaticResponse(testMail), "Mail should be filtered")
}

/*
	Specified header should be filtered out.
*/
func TestRepliesFilter_IsAutomaticResponse_PrecendenceSetToJunk(t *testing.T) {
	// Filter only looking at headers, so the rest can be ignored:
	testMail := mailData.Mail{Headers: []string{"precedence:junk"}}

	testee := RepliesFilter{}
	assert.True(t, testee.IsAutomaticResponse(testMail), "Mail should be filtered")
}

/*
	Specified header should be filtered out.
*/
func TestRepliesFilter_IsAutomaticResponse_AutoSubmittedSetToAutoReplied(t *testing.T) {
	// Filter only looking at headers, so the rest can be ignored:
	testMail := mailData.Mail{Headers: []string{"auto-submitted:auto-replied"}}

	testee := RepliesFilter{}
	assert.True(t, testee.IsAutomaticResponse(testMail), "Mail should be filtered")
}

/*
	Valid mail should not be filtered.
*/
func TestRepliesFilter_IsAutomaticResponse_ReturnsFalse(t *testing.T) {
	// Filter only looking at headers, so the rest can be ignored:
	testMail := mailData.Mail{Headers: []string{}}

	testee := RepliesFilter{}
	assert.False(t, testee.IsAutomaticResponse(testMail), "Mail should not be filtered")
}
