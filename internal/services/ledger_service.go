package services

import (
	"sort"

	"github.com/searaaman/attendly/internal/domain"
)

// BuildLedger totals every member's bills across all sessions and their payments.
// Because it works over all history, any unpaid amount naturally carries forward.
func BuildLedger(members []domain.Member, sessions []domain.Session, payments []domain.Payment) []domain.LedgerEntry {
	entries := make(map[uint]*domain.LedgerEntry, len(members))
	for _, member := range members {
		entries[member.ID] = &domain.LedgerEntry{
			MemberID: member.ID,
			Name:     member.Name,
			Phone:    member.Phone,
		}
	}

	for _, session := range sessions {
		for _, bill := range CalculateSessionBills(session) {
			if entry, ok := entries[bill.MemberID]; ok {
				entry.TotalBill += bill.Amount
			}
		}
	}
	for _, payment := range payments {
		if entry, ok := entries[payment.MemberID]; ok {
			entry.TotalPaid += payment.Amount
		}
	}

	ledger := make([]domain.LedgerEntry, 0, len(entries))
	for _, entry := range entries {
		entry.Balance = entry.TotalPaid - entry.TotalBill
		ledger = append(ledger, *entry)
	}
	// Biggest debts first, then alphabetical.
	sort.Slice(ledger, func(i, j int) bool {
		if ledger[i].Balance != ledger[j].Balance {
			return ledger[i].Balance < ledger[j].Balance
		}
		return ledger[i].Name < ledger[j].Name
	})
	return ledger
}

// BuildMemberLedger is the ledger for a single member, broken down by session.
// Sessions appear if the member attended them or paid towards them.
func BuildMemberLedger(member domain.Member, sessions []domain.Session, payments []domain.Payment) domain.MemberLedger {
	paidBySession := make(map[uint]float64)
	var totalPaid float64
	for _, payment := range payments {
		if payment.MemberID != member.ID {
			continue
		}
		paidBySession[payment.SessionID] += payment.Amount
		totalPaid += payment.Amount
	}

	charges := []domain.SessionCharge{}
	var totalBill float64
	for _, session := range sessions {
		charge := domain.SessionCharge{
			SessionID: session.ID,
			StartTime: session.StartTime,
			EndTime:   session.EndTime,
			Paid:      paidBySession[session.ID],
		}
		for _, bill := range CalculateSessionBills(session) {
			if bill.MemberID == member.ID {
				charge.Bill = bill.Amount
			}
		}
		for _, slot := range session.TimeSlots {
			for _, p := range slot.Members {
				if p.ID == member.ID {
					charge.SlotCount++
				}
			}
		}
		if charge.SlotCount == 0 && charge.Paid == 0 {
			continue
		}
		totalBill += charge.Bill
		charges = append(charges, charge)
	}
	// Most recent session first.
	sort.Slice(charges, func(i, j int) bool {
		return charges[i].StartTime.After(charges[j].StartTime)
	})

	memberPayments := []domain.Payment{}
	for _, payment := range payments {
		if payment.MemberID == member.ID {
			memberPayments = append(memberPayments, payment)
		}
	}

	return domain.MemberLedger{
		LedgerEntry: domain.LedgerEntry{
			MemberID:  member.ID,
			Name:      member.Name,
			Phone:     member.Phone,
			TotalBill: totalBill,
			TotalPaid: totalPaid,
			Balance:   totalPaid - totalBill,
		},
		Sessions: charges,
		Payments: memberPayments,
	}
}
