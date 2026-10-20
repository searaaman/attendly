// Shapes returned by the Go API. Keep in sync with internal/domain.

export interface Member {
  id: number;
  name: string;
  phone: string;
}

export interface TimeSlot {
  id: number;
  session_id: number;
  start_time: string;
  end_time: string;
  units_booked: number;
  members: Member[] | null;
}

export interface Session {
  id: number;
  start_time: string;
  end_time: string;
  unit_rate: number;
  time_slots: TimeSlot[] | null;
}

export interface MemberBill {
  member_id: number;
  name: string;
  amount: number;
}

export interface LedgerEntry {
  member_id: number;
  name: string;
  phone: string;
  total_bill: number;
  total_paid: number;
  // total_paid - total_bill: negative means the member owes money.
  balance: number;
}

export interface SessionCharge {
  session_id: number;
  start_time: string;
  end_time: string;
  slot_count: number;
  bill: number;
  paid: number;
}

export interface MemberLedger extends LedgerEntry {
  sessions: SessionCharge[];
  payments: Payment[];
}

export interface Payment {
  id: number;
  member_id: number;
  session_id: number;
  amount: number;
  created_at: string;
  member?: Member;
}

export interface User {
  id: number;
  name: string;
  email: string;
}
