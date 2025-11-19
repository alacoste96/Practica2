package main

import (
	"testing"
)

func addMechsWithSkills(t *testing.T, g *Garage, nMech, nElec, nBody int) {
	t.Helper()

	nextID := len(g.mechanics) + 1

	for i := 0; i < nMech; i++ {
		m := Mechanic{
			id:    MechanicID(nextID),
			skill: MECHSKILL,
		}
		nextID++
		g.mechanics = append(g.mechanics, &m)
	}
	for i := 0; i < nElec; i++ {
		m := Mechanic{
			id:    MechanicID(nextID),
			skill: ELECTRICSKILL,
		}
		nextID++
		g.mechanics = append(g.mechanics, &m)
	}
	for i := 0; i < nBody; i++ {
		m := Mechanic{
			id:    MechanicID(nextID),
			skill: BODYSKILL,
		}
		nextID++
		g.mechanics = append(g.mechanics, &m)
	}
}

func Test3MechsVS1others(t *testing.T) {
	g := newGarage()

	blocks := 2
	nMech := 3 * blocks
	nElec := 1 * blocks
	nBody := 1 * blocks

	addMechsWithSkills(t, g, nMech, nElec, nBody)

	gotMech := g.countSkill(MECHSKILL)
	gotElec := g.countSkill(ELECTRICSKILL)
	gotBody := g.countSkill(BODYSKILL)
	gotTotal := len(g.mechanics)

	// Comprobación básica de conteos
	if gotMech != nMech || gotElec != nElec || gotBody != nBody {
		t.Fatalf("Conteos incorrectos: mech=%d (want %d), elec=%d (want %d), body=%d (want %d)",
			gotMech, nMech, gotElec, nElec, gotBody, nBody)
	}

	if gotTotal != nMech+nElec+nBody {
		t.Fatalf("Total de mecánicos = %d; esperado %d", gotTotal, nMech+nElec+nBody)
	}

	// Comparativas de ratios:
	// 3 mecánicos de mecánica por cada 1 de eléctrica
	if gotMech != 3*gotElec {
		t.Errorf("Se esperaba relación 3:1 entre mecánica y eléctrica; got mech=%d, elec=%d", gotMech, gotElec)
	}

	// 3 mecánicos de mecánica por cada 1 de carrocería
	if gotMech != 3*gotBody {
		t.Errorf("Se esperaba relación 3:1 entre mecánica y carrocería; got mech=%d, body=%d", gotMech, gotBody)
	}

	// misma cantidad de eléctricos que de carroceros
	if gotElec != gotBody {
		t.Errorf("Se esperaba mismo nº de eléctricos y carrocería; got elec=%d, body=%d", gotElec, gotBody)
	}
}

func Test1MechVS3others(t *testing.T) {
	g := newGarage()

	blocks := 2
	nMech := 1 * blocks
	nElec := 3 * blocks
	nBody := 3 * blocks

	addMechsWithSkills(t, g, nMech, nElec, nBody)

	gotMech := g.countSkill(MECHSKILL)
	gotElec := g.countSkill(ELECTRICSKILL)
	gotBody := g.countSkill(BODYSKILL)
	gotTotal := len(g.mechanics)

	if gotMech != nMech || gotElec != nElec || gotBody != nBody {
		t.Fatalf("Conteos incorrectos: mech=%d (want %d), elec=%d (want %d), body=%d (want %d)",
			gotMech, nMech, gotElec, nElec, gotBody, nBody)
	}

	if gotTotal != nMech+nElec+nBody {
		t.Fatalf("Total de mecánicos = %d; esperado %d", gotTotal, nMech+nElec+nBody)
	}

	// 1 mecánico de mecánica por cada 3 de eléctrica
	if gotElec != 3*gotMech {
		t.Errorf("Se esperaba relación 1:3 entre mecánica y eléctrica; got mech=%d, elec=%d", gotMech, gotElec)
	}

	// 1 mecánico de mecánica por cada 3 de carrocería
	if gotBody != 3*gotMech {
		t.Errorf("Se esperaba relación 1:3 entre mecánica y carrocería; got mech=%d, body=%d", gotMech, gotBody)
	}

	// misma cantidad de eléctricos que de carroceros
	if gotElec != gotBody {
		t.Errorf("Se esperaba mismo nº de eléctricos y carrocería; got elec=%d, body=%d", gotElec, gotBody)
	}
}
