import api from "./client";
import type {
  LedgerEntry,
  Payment,
  Member,
  MemberBill,
  MemberLedger,
  Session,
  TimeSlot,
  User,
} from "../types";

// Auth

export async function login(email: string, password: string) {
  const { data } = await api.post<{ token: string; user: User }>("/login", { email, password });
  return data;
}

export async function register(name: string, email: string, password: string) {
  const { data } = await api.post<{ user: User }>("/register", { name, email, password });
  return data;
}

// Sessions

export interface NewSession {
  start_time: string;
  end_time: string;
  unit_rate: number;
  hourly_slots: boolean;
  default_units: number;
}

export async function listSessions() {
  return (await api.get<Session[]>("/sessions")).data;
}

export async function getSession(id: number) {
  return (await api.get<Session>(`/sessions/${id}`)).data;
}

export async function createSession(body: NewSession) {
  return (await api.post<Session>("/sessions", body)).data;
}

export async function deleteSession(id: number) {
  await api.delete(`/sessions/${id}`);
}

export async function getSessionBilling(id: number) {
  return (await api.get<MemberBill[]>(`/sessions/${id}/billing`)).data;
}

// Time slots

export async function createTimeSlot(
  sessionId: number,
  body: { start_time: string; end_time: string; units_booked: number },
) {
  return (await api.post<TimeSlot>(`/sessions/${sessionId}/timeslots`, body)).data;
}

export async function updateTimeSlotUnits(slotId: number, units: number) {
  return (await api.put<TimeSlot>(`/timeslots/${slotId}`, { units_booked: units })).data;
}

export async function deleteTimeSlot(slotId: number) {
  await api.delete(`/timeslots/${slotId}`);
}

export async function assignMember(slotId: number, memberId: number) {
  await api.post(`/timeslots/${slotId}/members`, { member_id: memberId });
}

export async function removeMember(slotId: number, memberId: number) {
  await api.delete(`/timeslots/${slotId}/members/${memberId}`);
}

// Members and ledger

export async function listMembers() {
  return (await api.get<Member[]>("/members")).data;
}

export async function createMember(name: string, phone: string) {
  return (await api.post<Member>("/members", { name, phone })).data;
}

export async function getLedger() {
  return (await api.get<LedgerEntry[]>("/ledger")).data;
}

export async function getMemberLedger(id: number) {
  return (await api.get<MemberLedger>(`/members/${id}/ledger`)).data;
}

// Payments

export async function listPayments(filter: { member_id?: number; session_id?: number } = {}) {
  return (await api.get<Payment[]>("/payments", { params: filter })).data;
}

export async function createPayment(memberId: number, sessionId: number, amount: number) {
  return (
    await api.post<Payment>("/payments", {
      member_id: memberId,
      session_id: sessionId,
      amount,
    })
  ).data;
}

export async function deletePayment(id: number) {
  await api.delete(`/payments/${id}`);
}
