package rcnb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncode(t *testing.T) {
	d, _ := Encode([]byte("rcnb"))
	assert.Equal(t, d, "ɌcńƁȓČņÞ")
	d, _ = Encode([]byte("The world is based on RC thus RCNB"))
	assert.Equal(t, d, "ȐčnÞȒċƝÞɌČǹƄɌcŅƀȒĉȠƀȒȼȵƄŕƇŅbȒCŅßȒČnƅŕƇņƁȓċȵƀȐĉņþŕƇņÞȒȻŅBɌCňƀȐĉņþƦȻƝƃ")
}
