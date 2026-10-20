import { isZero, money } from "../lib/format";

// Balance is total_paid - total_bill, so negative means the member owes.
function Balance({ amount }: { amount: number }) {
  if (isZero(amount)) return <span className="pill pill-neutral">Settled</span>;
  if (amount < 0) return <span className="pill pill-owes">Owes {money(-amount)}</span>;
  return <span className="pill pill-credit">Credit {money(amount)}</span>;
}

export default Balance;
