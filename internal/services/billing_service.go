package services

import (
	"sort"

	"github.com/searaaman/attendly/internal/domain"
)

// CalculateSessionBills splits each slot's unit cost evenly between the members
// in that slot and totals the result per member. Empty slots are not billed.
func CalculateSessionBills(session domain.Session) []domain.MemberBill {
	memberBills := make(map[uint]*domain.MemberBill)
	for _, slot := range session.TimeSlots {
		if len(slot.Members) == 0 {
			continue
		}
		slotCost := float64(slot.UnitsBooked) * session.UnitRate
		costPerMember := slotCost / float64(len(slot.Members))
		for _, member := range slot.Members {
			bill, exists := memberBills[member.ID]
			if exists {
				bill.Amount += costPerMember
			} else {
				memberBills[member.ID] = &domain.MemberBill{
					MemberID: member.ID,
					Name:     member.Name,
					Amount:   costPerMember,
				}
			}
		}
	}

	bills := make([]domain.MemberBill, 0, len(memberBills))
	for _, bill := range memberBills {
		bills = append(bills, *bill)
	}
	// Map iteration order is random; sort so API responses are stable.
	sort.Slice(bills, func(i, j int) bool {
		if bills[i].Name != bills[j].Name {
			return bills[i].Name < bills[j].Name
		}
		return bills[i].MemberID < bills[j].MemberID
	})
	return bills
}
