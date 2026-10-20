import { useState, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { errorMessage } from "../api/client";
import { createMember, getLedger } from "../api/endpoints";
import Balance from "../components/Balance";
import { ErrorBox, Loading } from "../components/Status";
import { money } from "../lib/format";
import { useLoad } from "../lib/useLoad";

function Members() {
  const navigate = useNavigate();
  const { data: ledger, error: loadError, reload } = useLoad(getLedger);
  const [search, setSearch] = useState("");
  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  const handleAdd = async (e: FormEvent) => {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await createMember(name.trim(), phone.trim());
      setName("");
      setPhone("");
      reload();
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  };

  const rows = (ledger ?? [])
    .filter((entry) => entry.name.toLowerCase().includes(search.toLowerCase()))
    .sort((a, b) => a.name.localeCompare(b.name));

  return (
    <div className="stack">
      <div className="page-head">
        <div>
          <h1>Members</h1>
          <p className="muted">Balances include every session attended and every payment made.</p>
        </div>
      </div>

      <form className="card" onSubmit={handleAdd}>
        <div className="card-head">
          <h2>Add member</h2>
        </div>
        <ErrorBox message={error} />
        <div className="row" style={{ alignItems: "flex-end" }}>
          <label style={{ flex: 2, minWidth: 160 }}>
            Name
            <input required value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label style={{ flex: 1, minWidth: 140 }}>
            Phone (optional)
            <input type="tel" value={phone} onChange={(e) => setPhone(e.target.value)} />
          </label>
          <button className="btn" disabled={busy || !name.trim()}>
            Add
          </button>
        </div>
      </form>

      <section className="card card-flush">
        <div className="card-head">
          <h2>All members</h2>
          <input
            type="search"
            placeholder="Search"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            style={{ maxWidth: 220 }}
          />
        </div>
        <ErrorBox message={loadError} />
        {!ledger && !loadError && <Loading />}
        {ledger && rows.length === 0 && <p className="empty">No members found.</p>}
        {rows.length > 0 && (
          <div className="table-wrap" style={{ marginTop: 12 }}>
            <table>
              <thead>
                <tr>
                  <th>Name</th>
                  <th className="hide-sm">Phone</th>
                  <th className="num hide-sm">Billed</th>
                  <th className="num hide-sm">Paid</th>
                  <th className="num">Balance</th>
                </tr>
              </thead>
              <tbody>
                {rows.map((entry) => (
                  <tr
                    key={entry.member_id}
                    className="clickable"
                    onClick={() => navigate(`/members/${entry.member_id}`)}
                  >
                    <td style={{ fontWeight: 500 }}>{entry.name}</td>
                    <td className="muted hide-sm">{entry.phone || "—"}</td>
                    <td className="num hide-sm">{money(entry.total_bill)}</td>
                    <td className="num hide-sm">{money(entry.total_paid)}</td>
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
    </div>
  );
}

export default Members;
