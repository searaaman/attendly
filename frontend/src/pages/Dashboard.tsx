import { Link, useNavigate } from "react-router-dom";
import { getLedger, listSessions } from "../api/endpoints";
import Balance from "../components/Balance";
import SessionList from "../components/SessionList";
import { ErrorBox, Loading } from "../components/Status";
import { isZero, money } from "../lib/format";
import { useLoad } from "../lib/useLoad";

function Dashboard() {
  const navigate = useNavigate();
  const { data, error } = useLoad(() => Promise.all([getLedger(), listSessions()]));

  if (error) return <ErrorBox message={error} />;
  if (!data) return <Loading />;

  const [ledger, sessions] = data;
  const owing = ledger.filter((entry) => entry.balance < 0 && !isZero(entry.balance));
  const outstanding = owing.reduce((total, entry) => total - entry.balance, 0);
  const collected = ledger.reduce((total, entry) => total + entry.total_paid, 0);

  return (
    <div className="stack">
      <div className="page-head">
        <div>
          <h1>Dashboard</h1>
          <p className="muted">Who owes what, across every session.</p>
        </div>
        <Link to="/sessions/new" className="btn">
          + New session
        </Link>
      </div>

      <div className="stats">
        <div className="stat">
          <div className="stat-label">Outstanding</div>
          <div className={`stat-value ${outstanding > 0 ? "owes" : ""}`}>{money(outstanding)}</div>
        </div>
        <div className="stat">
          <div className="stat-label">Members owing</div>
          <div className="stat-value">{owing.length}</div>
        </div>
        <div className="stat">
          <div className="stat-label">Collected</div>
          <div className="stat-value">{money(collected)}</div>
        </div>
        <div className="stat">
          <div className="stat-label">Sessions</div>
          <div className="stat-value">{sessions.length}</div>
        </div>
      </div>

      <div className="grid-2">
        <section className="card card-flush">
          <div className="card-head">
            <h2>Who owes</h2>
            <Link to="/members" className="small">
              All members →
            </Link>
          </div>
          {owing.length === 0 ? (
            <p className="empty">Everyone is settled up.</p>
          ) : (
            <div className="table-wrap" style={{ marginTop: 12 }}>
              <table>
                <tbody>
                  {owing.slice(0, 8).map((entry) => (
                    <tr
                      key={entry.member_id}
                      className="clickable"
                      onClick={() => navigate(`/members/${entry.member_id}`)}
                    >
                      <td>{entry.name}</td>
                      <td className="num">
                        <Balance amount={entry.balance} />
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>

        <section className="card card-flush">
          <div className="card-head">
            <h2>Recent sessions</h2>
            <Link to="/sessions" className="small">
              All sessions →
            </Link>
          </div>
          <div style={{ marginTop: 12 }}>
            <SessionList sessions={sessions.slice(0, 5)} />
          </div>
        </section>
      </div>
    </div>
  );
}

export default Dashboard;
