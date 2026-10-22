import { useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { errorMessage } from "../api/client";
import { deletePayment, getMemberLedger } from "../api/endpoints";
import Balance from "../components/Balance";
import { ErrorBox, Loading } from "../components/Status";
import { date, dateTime, isZero, money, shortDate, timeRange } from "../lib/format";
import { useLoad } from "../lib/useLoad";

// Builds a wa.me link, assuming Indian numbers when no country code is given.
function whatsappLink(phone: string, message: string): string | null {
  let digits = phone.replace(/\D/g, "");
  if (digits.length === 10) digits = `91${digits}`;
  if (digits.length < 11) return null;
  return `https://wa.me/${digits}?text=${encodeURIComponent(message)}`;
}

function MemberDetail() {
  const { id } = useParams();
  const memberId = Number(id);
  const navigate = useNavigate();
  const { data: ledger, error: loadError, reload } = useLoad(() => getMemberLedger(memberId), memberId);
  const [error, setError] = useState<string | null>(null);

  if (loadError) return <ErrorBox message={loadError} />;
  if (!ledger) return <Loading />;

  const owes = ledger.balance < 0 && !isZero(ledger.balance);
  const reminder = whatsappLink(
    ledger.phone,
    `Hi ${ledger.name}, your outstanding balance is ${money(-ledger.balance)}. Thanks!`,
  );

  const handleDelete = async (paymentId: number, amount: number) => {
    if (!confirm(`Delete ${money(amount)} payment?`)) return;
    try {
      await deletePayment(paymentId);
      reload();
    } catch (err) {
      setError(errorMessage(err));
    }
  };

  return (
    <div className="stack">
      <div className="page-head">
        <div>
          <p className="small">
            <Link to="/members">← Members</Link>
          </p>
          <h1>{ledger.name}</h1>
          {ledger.phone && <p className="muted">{ledger.phone}</p>}
        </div>
        {owes && reminder && (
          <a className="btn btn-secondary" href={reminder} target="_blank" rel="noreferrer">
            Send WhatsApp reminder
          </a>
        )}
      </div>

      <ErrorBox message={error} />

      <div className="stats">
        <div className="stat">
          <div className="stat-label">Balance</div>
          <div className="stat-value" style={{ fontSize: "1.1rem", paddingTop: 6 }}>
            <Balance amount={ledger.balance} />
          </div>
        </div>
        <div className="stat">
          <div className="stat-label">Total billed</div>
          <div className="stat-value">{money(ledger.total_bill)}</div>
        </div>
        <div className="stat">
          <div className="stat-label">Total paid</div>
          <div className="stat-value">{money(ledger.total_paid)}</div>
        </div>
        <div className="stat">
          <div className="stat-label">Sessions</div>
          <div className="stat-value">{ledger.sessions.length}</div>
        </div>
      </div>

      <div className="grid-2">
        <section className="card card-flush">
          <div className="card-head">
            <h2>Sessions</h2>
          </div>
          {ledger.sessions.length === 0 ? (
            <p className="empty">No attendance yet.</p>
          ) : (
            <div className="table-wrap" style={{ marginTop: 12 }}>
              <table>
                <thead>
                  <tr>
                    <th>Session</th>
                    <th className="num">Bill</th>
                    <th className="num">Paid</th>
                  </tr>
                </thead>
                <tbody>
                  {ledger.sessions.map((charge) => (
                    <tr
                      key={charge.session_id}
                      className="clickable"
                      onClick={() => navigate(`/sessions/${charge.session_id}`)}
                    >
                      <td>
                        <div style={{ fontWeight: 500 }}>{date(charge.start_time)}</div>
                        <div className="muted small">
                          {timeRange(charge.start_time, charge.end_time)} · {charge.slot_count} slot
                          {charge.slot_count === 1 ? "" : "s"}
                        </div>
                      </td>
                      <td className="num">{money(charge.bill)}</td>
                      <td className="num">{money(charge.paid)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>

        <section className="card card-flush">
          <div className="card-head">
            <h2>Payments</h2>
          </div>
          {ledger.payments.length === 0 ? (
            <p className="empty">No payments yet. Record them from a session page.</p>
          ) : (
            <div className="table-wrap" style={{ marginTop: 12 }}>
              <table>
                <tbody>
                  {ledger.payments.map((payment) => {
                    const session = ledger.sessions.find((s) => s.session_id === payment.session_id);
                    return (
                      <tr key={payment.id}>
                        <td>
                          <div>{dateTime(payment.created_at)}</div>
                          {session && (
                            <div className="muted small">
                              for <Link to={`/sessions/${payment.session_id}`}>{shortDate(session.start_time)}</Link>{" "}
                              session
                            </div>
                          )}
                        </td>
                        <td className="num">{money(payment.amount)}</td>
                        <td className="num">
                          <button
                            className="btn btn-ghost btn-sm"
                            onClick={() => handleDelete(payment.id, payment.amount)}
                          >
                            Delete
                          </button>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>
    </div>
  );
}

export default MemberDetail;
