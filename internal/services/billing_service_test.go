package services

import (
	"github.com/searaaman/attendly/internal/domain"
	"testing"
)

func TestCalculateSessionBillsTwoMembers(t *testing.T) {
	rahul := domain.Member{
		ID:   1,
		Name: "Rahul",
	}

	anand := domain.Member{
		ID:   2,
		Name: "Anand",
	}

	slot := domain.TimeSlot{
		UnitsBooked: 1,
		Members: []domain.Member{
			rahul,
			anand,
		},
	}

	session := domain.Session{
		UnitRate: 300,
		TimeSlots: []domain.TimeSlot{
			slot,
		},
	}

	bills := CalculateSessionBills(session)

	if len(bills) != 2 {
		t.Fatalf("expected 2 bills, got %d", len(bills))
	}

	billMap := make(map[uint]float64)

	for _, bill := range bills {
		billMap[bill.MemberID] = bill.Amount
	}

	if billMap[rahul.ID] != 150 {
		t.Errorf("expected Rahul to owe 150, got %.2f", billMap[rahul.ID])
	}

	if billMap[anand.ID] != 150 {
		t.Errorf("expected Anand to owe 150, got %.2f", billMap[anand.ID])
	}
}

func TestCalculateSessionBillsMultipleTimeSlots(t *testing.T) {
	rahul := domain.Member{
		ID:   1,
		Name: "Rahul",
	}

	anand := domain.Member{
		ID:   2,
		Name: "Anand",
	}

	slot1 := domain.TimeSlot{
		UnitsBooked: 1,
		Members: []domain.Member{
			rahul,
			anand,
		},
	}

	slot2 := domain.TimeSlot{
		UnitsBooked: 1,
		Members: []domain.Member{
			rahul,
		},
	}

	session := domain.Session{
		UnitRate: 300,
		TimeSlots: []domain.TimeSlot{
			slot1,
			slot2,
		},
	}

	bills := CalculateSessionBills(session)

	if len(bills) != 2 {
		t.Fatalf("expected 2 bills, got %d", len(bills))
	}

	billMap := make(map[uint]float64)

	for _, bill := range bills {
		billMap[bill.MemberID] = bill.Amount
	}

	if billMap[rahul.ID] != 450 {
		t.Errorf("expected Rahul to owe 450, got %.2f", billMap[rahul.ID])
	}

	if billMap[anand.ID] != 150 {
		t.Errorf("expected Anand to owe 150, got %.2f", billMap[anand.ID])
	}
}

func TestCalculateSessionBillsMultipleUnits(t *testing.T) {

	rahul := domain.Member{
		ID:   1,
		Name: "Rahul",
	}

	anand := domain.Member{
		ID:   2,
		Name: "Anand",
	}

	slot := domain.TimeSlot{
		UnitsBooked: 2,
		Members: []domain.Member{
			rahul,
			anand,
		},
	}

	session := domain.Session{
		UnitRate: 300,
		TimeSlots: []domain.TimeSlot{
			slot,
		},
	}

	bills := CalculateSessionBills(session)

	if len(bills) != 2 {
		t.Fatalf("expected 2 bills, got %d", len(bills))
	}

	billMap := make(map[uint]float64)

	for _, bill := range bills {
		billMap[bill.MemberID] = bill.Amount
	}

	if billMap[rahul.ID] != 300 {
		t.Errorf("expected Rahul to owe 300, got %.2f", billMap[rahul.ID])
	}

	if billMap[anand.ID] != 300 {
		t.Errorf("expected Anand to owe 300, got %.2f", billMap[anand.ID])
	}
}

func TestCalculateSessionBillsEmptyTimeSlot(t *testing.T) {

	slot := domain.TimeSlot{
		UnitsBooked: 1,
		Members:     []domain.Member{},
	}

	session := domain.Session{
		UnitRate: 300,
		TimeSlots: []domain.TimeSlot{
			slot,
		},
	}

	bills := CalculateSessionBills(session)

	if len(bills) != 0 {
		t.Fatalf("expected 0 bills, got %d", len(bills))
	}
}

func TestCalculateSessionBillsEmptySession(t *testing.T) {

	session := domain.Session{
		UnitRate:  300,
		TimeSlots: []domain.TimeSlot{},
	}

	bills := CalculateSessionBills(session)

	if len(bills) != 0 {
		t.Fatalf("expected 0 bills, got %d", len(bills))
	}
}

// The session we verified by hand through the API.
func TestCalculateSessionBillsRealSession(t *testing.T) {
	rahul := domain.Member{ID: 1, Name: "Rahul"}
	akhil := domain.Member{ID: 2, Name: "Akhil"}
	anand := domain.Member{ID: 3, Name: "Anand"}

	session := domain.Session{
		UnitRate: 300,
		TimeSlots: []domain.TimeSlot{
			{UnitsBooked: 2, Members: []domain.Member{rahul, akhil, anand}},
			{UnitsBooked: 1, Members: []domain.Member{rahul, anand}},
		},
	}

	bills := CalculateSessionBills(session)

	expected := []domain.MemberBill{
		{MemberID: 2, Name: "Akhil", Amount: 200},
		{MemberID: 3, Name: "Anand", Amount: 350},
		{MemberID: 1, Name: "Rahul", Amount: 350},
	}
	if len(bills) != len(expected) {
		t.Fatalf("expected %d bills, got %d", len(expected), len(bills))
	}
	// Bills come back sorted by name.
	for i := range expected {
		if bills[i] != expected[i] {
			t.Errorf("bill %d: expected %+v, got %+v", i, expected[i], bills[i])
		}
	}
}
