package abc

import (
	"dartech-course"
	"github.com/magiconair/properties/assert"
	"testing"
)


var expected = dartech_course.song{"Paradise a", "Paradise group a", nil, nil}

func TestCreateSong_isSame(t *testing.T) {
	n := dartech_course.CreateSong("Paradise a", "Paradise group a", nil, nil)
	assert.Equal(t, n, expected)
}
