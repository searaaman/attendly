import { Link } from "react-router-dom";
import { money, timeRange } from "../lib/format";
import { sessionCost, sessionMembers, slotsOf } from "../lib/session";
import type { Session } from "../types";

function SessionList({ sessions }: { sessions: Session[] }) {
  if (sessions.length === 0) {
    return <p className="empty">No sessions yet.</p>;
  }
  return (
    <div className="session-list">
      {sessions.map((session) => {
        const start = new Date(session.start_time);
        const slots = slotsOf(session).length;
        const members = sessionMembers(session).length;
        return (
          <Link key={session.id} to={`/sessions/${session.id}`} className="session-item">
            <div className="date-badge">
              <div className="day">{start.getDate()}</div>
              <div className="mon">
                {start.toLocaleDateString("en-IN", { month: "short", weekday: "short" })}
              </div>
            </div>
            <div>
              <div style={{ fontWeight: 600 }}>{timeRange(session.start_time, session.end_time)}</div>
              <div className="muted small">
                {slots} slot{slots === 1 ? "" : "s"} · {members} member{members === 1 ? "" : "s"} ·{" "}
                {money(session.unit_rate)}/unit/hr
              </div>
            </div>
            <div className="num" style={{ fontWeight: 600 }}>
              {money(sessionCost(session))}
            </div>
          </Link>
        );
      })}
    </div>
  );
}

export default SessionList;
