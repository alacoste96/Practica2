package main

import (
	"fmt"
	"testing"
)

func genGarage(t *testing.T, n int, kind IssueType) *Garage {
	t.Helper()

	g := newGarage()

	for i := 0; i < n; i++ {
		vid := VehicleID(
			fmt.Sprintf("TEST-%d", i+1),
		)

		incID := IncidenceID(i + 1)

		v := Vehicle{
			id:     vid,
			issues: []IncidenceID{incID},
		}
		inc := Incidence{
			id:        incID,
			vehicleID: vid,
			kind:      kind,
		}

		g.vehicles = append(g.vehicles, &v)
		g.issues = append(g.issues, &inc)
	}

	return g
}

func TestDupCarsWithSameIssue(t *testing.T) {
	const baseN = 3

	g1 := genGarage(t, baseN, MECHTYPE)
	got1 := g1.CountVehiclesWithSingleIssue(MECHTYPE)

	if got1 != baseN {
		t.Fatalf("para %d coches esperábamos contador=%d; got=%d", baseN, baseN, got1)
	}

	g2 := genGarage(t, 2*baseN, MECHTYPE)
	got2 := g2.CountVehiclesWithSingleIssue(MECHTYPE)

	if got2 != 2*baseN {
		t.Fatalf("para %d coches esperábamos contador=%d; got=%d", 2*baseN, 2*baseN, got2)
	}

	if got2 != 2*got1 {
		t.Errorf("se esperaba que al duplicar coches (%d -> %d), el contador también se duplicase (%d -> %d), pero obtuvimos got1=%d, got2=%d",
			baseN, 2*baseN, baseN, 2*baseN, got1, got2)
	}

}
