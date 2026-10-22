package services

import (
	"testing"
	"time"

	"github.com/searaaman/attendly/internal/domain"
)

func TestBuildHourlySlots(t *testing.T) {
	start := time.Date(2026, 9, 1, 18, 0, 0, 0, time.UTC)

	slots := BuildHourlySlots(start, start.Add(4*time.Hour), 2)

	if len(slots) != 4 {
		t.Fatalf("expected 4 slots, got %d", len(slots))
	}
	for i, slot := range slots {
		wantStart := start.Add(time.Duration(i) * time.Hour)
		if !slot.StartTime.Equal(wantStart) || !slot.EndTime.Equal(wantStart.Add(time.Hour)) {
			t.Errorf("slot %d: got %v-%v", i, slot.StartTime, slot.EndTime)
		}
		if slot.UnitsBooked != 2 {
			t.Errorf("slot %d: expected 2 units, got %d", i, slot.UnitsBooked)
		}
	}
}

func TestBuildHourlySlotsPartialLastHour(t *testing.T) {
	start := time.Date(2026, 9, 1, 18, 0, 0, 0, time.UTC)
	end := start.Add(90 * time.Minute)

	slots := BuildHourlySlots(start, end, 1)

	if len(slots) != 2 {
		t.Fatalf("expected 2 slots, got %d", len(slots))
	}
	if !slots[1].EndTime.Equal(end) {
		t.Errorf("expected last slot to end at %v, got %v", end, slots[1].EndTime)
	}
}

func TestNewSessionRejectsEndBeforeStart(t *testing.T) {
	start := time.Date(2026, 9, 1, 18, 0, 0, 0, time.UTC)

	_, err := NewSession(domain.CreateSessionRequest{StartTime: start, EndTime: start})

	if err != ErrInvalidTimeRange {
		t.Errorf("expected ErrInvalidTimeRange, got %v", err)
	}
}

func TestNewSessionWithoutHourlySlots(t *testing.T) {
	start := time.Date(2026, 9, 1, 18, 0, 0, 0, time.UTC)

	session, err := NewSession(domain.CreateSessionRequest{
		StartTime: start, EndTime: start.Add(2 * time.Hour), UnitRate: 300,
	})

	if err != nil {
		t.Fatal(err)
	}
	if len(session.TimeSlots) != 0 {
		t.Errorf("expected no slots, got %d", len(session.TimeSlots))
	}
}

func TestValidateTimeSlot(t *testing.T) {
	start := time.Date(2026, 9, 1, 18, 0, 0, 0, time.UTC)
	session := domain.Session{StartTime: start, EndTime: start.Add(4 * time.Hour)}

	cases := []struct {
		name       string
		start, end time.Time
		want       error
	}{
		{"inside", start.Add(time.Hour), start.Add(2 * time.Hour), nil},
		{"whole session", start, start.Add(4 * time.Hour), nil},
		{"end before start", start.Add(2 * time.Hour), start.Add(time.Hour), ErrInvalidTimeRange},
		{"starts early", start.Add(-time.Hour), start.Add(time.Hour), ErrSlotOutsideSession},
		{"ends late", start.Add(3 * time.Hour), start.Add(5 * time.Hour), ErrSlotOutsideSession},
	}
	for _, tc := range cases {
		if got := ValidateTimeSlot(session, tc.start, tc.end); got != tc.want {
			t.Errorf("%s: expected %v, got %v", tc.name, tc.want, got)
		}
	}
}
