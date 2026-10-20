import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { errorMessage } from "../api/client";
import { createSession } from "../api/endpoints";
import { ErrorBox } from "../components/Status";
import { money, time, todayInputValue } from "../lib/format";

// Works out the session's start and end. An end time at or before the start
// time is taken to be after midnight, e.g. 21:00 to 01:00.
function sessionRange(day: string, start: string, end: string) {
  const startDate = new Date(`${day}T${start}`);
  const endDate = new Date(`${day}T${end}`);
  if (endDate <= startDate) endDate.setDate(endDate.getDate() + 1);
  return { startDate, endDate };
}

function hourlyPreview(startDate: Date, endDate: Date): string[] {
  const labels: string[] = [];
  for (let s = new Date(startDate); s < endDate; s = new Date(s.getTime() + 3_600_000)) {
    const e = new Date(Math.min(s.getTime() + 3_600_000, endDate.getTime()));
    labels.push(`${time(s.toISOString())}–${time(e.toISOString())}`);
    if (labels.length > 24) break;
  }
  return labels;
}

function CreateSession() {
  const navigate = useNavigate();
  const [day, setDay] = useState(todayInputValue());
  const [start, setStart] = useState("18:00");
  const [end, setEnd] = useState("22:00");
  const [unitRate, setUnitRate] = useState("300");
  const [hourly, setHourly] = useState(true);
  const [units, setUnits] = useState("2");
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  const valid = day && start && end;
  const { startDate, endDate } = valid
    ? sessionRange(day, start, end)
    : { startDate: null, endDate: null };
  const preview = startDate && endDate && hourly ? hourlyPreview(startDate, endDate) : [];
  const estimated = preview.length * Number(units || 0) * Number(unitRate || 0);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!startDate || !endDate) return;
    setBusy(true);
    setError(null);
    try {
      const session = await createSession({
        start_time: startDate.toISOString(),
        end_time: endDate.toISOString(),
        unit_rate: Number(unitRate),
        hourly_slots: hourly,
        default_units: Number(units),
      });
      navigate(`/sessions/${session.id}`);
    } catch (err) {
      setError(errorMessage(err));
      setBusy(false);
    }
  };

  return (
    <div>
      <div className="page-head">
        <div>
          <p className="small">
            <Link to="/sessions">← Sessions</Link>
          </p>
          <h1>New session</h1>
        </div>
      </div>

      <form className="card" onSubmit={handleSubmit} style={{ maxWidth: 640 }}>
        <ErrorBox message={error} />
        <div className="form-grid">
          <label>
            Date
            <input type="date" required value={day} onChange={(e) => setDay(e.target.value)} />
          </label>
          <label>
            Starts
            <input type="time" required value={start} onChange={(e) => setStart(e.target.value)} />
          </label>
          <label>
            Ends
            <input type="time" required value={end} onChange={(e) => setEnd(e.target.value)} />
          </label>
          <label>
            Unit price per hour (₹)
            <input
              type="number"
              min="0"
              step="any"
              required
              value={unitRate}
              onChange={(e) => setUnitRate(e.target.value)}
            />
          </label>
        </div>

        <div className="stack" style={{ marginTop: 20, gap: 12 }}>
          <label className="inline">
            <input type="checkbox" checked={hourly} onChange={(e) => setHourly(e.target.checked)} />
            Split into one-hour slots
          </label>
          {hourly && (
            <>
              <label style={{ maxWidth: 200 }}>
                Units booked per slot (courts, rooms, desks…)
                <input
                  type="number"
                  min="0"
                  required
                  value={units}
                  onChange={(e) => setUnits(e.target.value)}
                />
              </label>
              <div>
                <div className="slot-preview">
                  {preview.map((label) => (
                    <span key={label} className="slot-chip">
                      {label}
                    </span>
                  ))}
                </div>
                <p className="muted small" style={{ marginBottom: 0 }}>
                  {preview.length} slot{preview.length === 1 ? "" : "s"} · booking cost {money(estimated)}.
                  You can change units per slot afterwards.
                </p>
              </div>
            </>
          )}
          {!hourly && (
            <p className="hint" style={{ margin: 0 }}>
              You'll add time slots yourself on the next screen.
            </p>
          )}
        </div>

        <div className="form-actions">
          <button className="btn" disabled={busy || !valid}>
            {busy ? "Creating…" : "Create session"}
          </button>
          <Link to="/sessions" className="btn btn-secondary">
            Cancel
          </Link>
        </div>
      </form>
    </div>
  );
}

export default CreateSession;
