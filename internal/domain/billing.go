package domain

import "time"

type MemberBill struct {
	MemberID uint    `json:"member_id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
}

// LedgerEntry is one member's running account across every session.
// Balance = TotalPaid - TotalBill, so a negative balance means the member owes money.
type LedgerEntry struct {
	MemberID  uint    `json:"member_id"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	TotalBill float64 `json:"total_bill"`
	TotalPaid float64 `json:"total_paid"`
	Balance   float64 `json:"balance"`
}

// SessionCharge is what a member was billed and paid for one session.
type SessionCharge struct {
	SessionID uint      `json:"session_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	SlotCount int       `json:"slot_count"`
	Bill      float64   `json:"bill"`
	Paid      float64   `json:"paid"`
}

type MemberLedger struct {
	LedgerEntry
	Sessions []SessionCharge `json:"sessions"`
	Payments []Payment       `json:"payments"`
}
