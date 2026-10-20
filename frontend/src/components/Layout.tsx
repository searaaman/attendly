import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { clearSession, getUserName } from "../api/client";

function Layout() {
  const navigate = useNavigate();

  const logout = () => {
    clearSession();
    navigate("/login");
  };

  return (
    <div className="app">
      <header className="topbar">
        <div className="topbar-inner">
          <NavLink to="/" className="brand">
            <span className="brand-mark" aria-hidden="true" />
            Attendly
          </NavLink>
          <nav className="nav">
            <NavLink to="/" end>
              Dashboard
            </NavLink>
            <NavLink to="/sessions">Sessions</NavLink>
            <NavLink to="/members">Members</NavLink>
          </nav>
          <div className="topbar-user">
            <span className="muted hide-sm">{getUserName()}</span>
            <button className="btn btn-ghost btn-sm" onClick={logout}>
              Log out
            </button>
          </div>
        </div>
      </header>
      <main className="container">
        <Outlet />
      </main>
    </div>
  );
}

export default Layout;
