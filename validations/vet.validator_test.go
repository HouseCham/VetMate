package validations

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var hasSpecialTestNames = []string{
	"John123",
	"John!",
	"John@",
	"John#",
	"John$",
	"John%",
	"John^",
	"John&",
	"John*",
	"John(",
	"John)",
	"John-",
	"John_",
	"John+",
	"John=",
	"John{",
	"John[",
	"John}",
	"John]",
	"John|",
	"John:",
	"John;",
	"John'",
	"John<",
	"John,",
	"John>",
	"John.",
	"John?",
	"John/",
	"John ",
}
var hasNoSpecialTestNames = []string{
	"John",
	"JohnDoe",
	"JohnDoe",
	"Ramses",
	"Julio",
	"JulioCesar",
	"Romina",
	"RominaGonzalez",
	"Martion",
}

func TestHasSpecialCharacters(t *testing.T) {
	for _, name := range hasSpecialTestNames {
		require.True(t, HasSpecialCharacters(name))
	}

	for _, name := range hasNoSpecialTestNames {
		require.False(t, HasSpecialCharacters(name))
	}
}