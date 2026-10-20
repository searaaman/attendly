import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import Layout from "./components/Layout";
import RequireAuth from "./components/RequireAuth";
import CreateSession from "./pages/CreateSession";
import Dashboard from "./pages/Dashboard";
import Login from "./pages/Login";
import MemberDetail from "./pages/MemberDetail";
import Members from "./pages/Members";
import Register from "./pages/Register";
import SessionDetail from "./pages/SessionDetail";
import Sessions from "./pages/Sessions";

// On GitHub Pages the app lives under /<repo>/, which Vite exposes as BASE_URL.
const basename = import.meta.env.BASE_URL.replace(/\/$/, "") || undefined;

function App() {
  return (
    <BrowserRouter basename={basename}>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route
          element={
            <RequireAuth>
              <Layout />
            </RequireAuth>
          }
        >
          <Route path="/" element={<Dashboard />} />
          <Route path="/sessions" element={<Sessions />} />
          <Route path="/sessions/new" element={<CreateSession />} />
          <Route path="/sessions/:id" element={<SessionDetail />} />
          <Route path="/members" element={<Members />} />
          <Route path="/members/:id" element={<MemberDetail />} />
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
