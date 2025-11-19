package main

import (
	"testing"
)

func (g *Garage) gen3Mechs(t *testing.T, start int) {
	t.Helper()

	var nskill SkillType
	for i := start; i < (start + 3); i++ {
		mid := MechanicID(i + 1)
		switch i % 3 {
		case 0:
			nskill = MECHSKILL
		case 1:
			nskill = ELECTRICSKILL
		case 2:
			nskill = BODYSKILL
		}
		mech := Mechanic{
			id:    mid,
			skill: nskill,
		}
		g.mechanics = append(g.mechanics, &mech)
		g.genSlots()
	}
}

func (g *Garage) sameSkillQuantity(nskills int) bool {
	return g.countSkill(MECHSKILL) == nskills &&
		g.countSkill(BODYSKILL) == nskills &&
		g.countSkill(ELECTRICSKILL) == nskills
}

func TestDupMechs(t *testing.T) {
	const mt1 = 3
	const mt2 = 6
	const sk1 = 1
	const sk2 = 2
	const s1 = 6
	const s2 = 12
	var got int

	g := newGarage()
	g.gen3Mechs(t, 0)

	got = len(g.mechanics)
	if got != mt1 {
		t.Fatalf("Para %d mecánicos esperábamos contador=%d; got=%d", mt1, mt1, got)
	}

	if !g.sameSkillQuantity(sk1) {
		if g.countSkill(MECHSKILL) != sk1 {
			t.Fatalf("Mecánica fuera de rango")
		}
		if g.countSkill(BODYSKILL) != sk1 {
			t.Fatalf("Carrocería fuera de rango")
		}
		if g.countSkill(ELECTRICSKILL) != sk1 {
			t.Fatalf("Eléctrica fuera de rango")
		}
	}
	got = g.availableSlots()
	if got != s1 {
		t.Fatalf("Para %d mecánicos esperábamos %d slots: contador=%d; got=%d", mt1, s1, s1, got)
	}

	g.gen3Mechs(t, 3)

	got = len(g.mechanics)
	if got != mt2 {
		t.Fatalf("Para %d mecánicos esperábamos contador=%d; got=%d", mt2, mt2, got)
	}
	if !g.sameSkillQuantity(sk2) {
		if g.countSkill(MECHSKILL) != sk2 {
			t.Fatalf("Mecánica fuera de rango")
		}
		if g.countSkill(BODYSKILL) != sk2 {
			t.Fatalf("Carrocería fuera de rango")
		}
		if g.countSkill(ELECTRICSKILL) != sk2 {
			t.Fatalf("Eléctrica fuera de rango")
		}
	}
	got = g.availableSlots()
	if got != s2 {
		t.Fatalf("Para %d mecánicos esperábamos %d slots: contador=%d; got=%d", mt2, s2, s2, got)
	}
}
