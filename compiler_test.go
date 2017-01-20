package opentis100

import "testing"

func TestCompileMov(t *testing.T) {
	program := `@0
mov 1 acc
	`

	p, err := compile("testProg", program)

	if err != nil {
		t.Errorf("Unexpected compile error: %s", err.Error())
	}

	t.Run("Program has right number of sets", func(t *testing.T) {
		if len(p.Sets) != 1 {
			t.Errorf("Expected %d set(s), got %d", 1, len(p.Sets))
		}

	})

	t.Run("Program set has right number of instructions", func(t *testing.T) {
		set := p.Sets[0]

		if len(set.Instructions) > 1 {
			t.Errorf("Expected %d instruction(s), got %d", 1, len(set.Instructions))
		}
	})
}
