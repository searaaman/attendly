import type { Member, Session, TimeSlot } from "../types";

export function slotsOf(session: Session): TimeSlot[] {
  return session.time_slots ?? [];
}

export function membersOf(slot: TimeSlot): Member[] {
  return slot.members ?? [];
}

// Everyone who attended at least one slot, sorted by name.
export function sessionMembers(session: Session): Member[] {
  const byId = new Map<number, Member>();
  for (const slot of slotsOf(session)) {
    for (const member of membersOf(slot)) byId.set(member.id, member);
  }
  return [...byId.values()].sort((a, b) => a.name.localeCompare(b.name));
}

export function sessionCost(session: Session): number {
  return slotsOf(session).reduce(
    (total, slot) => total + slot.units_booked * session.unit_rate,
    0,
  );
}

// Cost of slots that someone actually attended; empty slots aren't billed.
export function billedCost(session: Session): number {
  return slotsOf(session)
    .filter((slot) => membersOf(slot).length > 0)
    .reduce((total, slot) => total + slot.units_booked * session.unit_rate, 0);
}

export function perMemberCost(session: Session, slot: TimeSlot): number | null {
  const count = membersOf(slot).length;
  if (count === 0) return null;
  return (slot.units_booked * session.unit_rate) / count;
}
