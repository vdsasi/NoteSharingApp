import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { useLocation } from 'react-router-dom';

import Login from './pages/Login';
import Register from './pages/Register';
import Notes from './pages/Notes';
import UserMenu from './components/UserMenu';
import Profile from './pages/Profile';
import ChangePassword from './pages/ChangePassword';
import Dashboard from './pages/Dashboard';
import CreateNote from "./pages/CreateNote";
import EditNote from "./pages/EditNote";
import ViewNote from "./pages/ViewNote";
import TrashPage from "./pages/TrashPage";
import SidebarLayout from './components/SidebarLayout';
import './app.css';

// Wrapper to conditionally show UserMenu
const AppContent = () => {
  const location = useLocation();
  const hideUserMenu = location.pathname === '/login' || location.pathname === '/register';

  return (
    <div>
      {/* {!hideUserMenu && <UserMenu />} */}
      <Routes>
        <Route path="/login" element={<Login/>}></Route>
        <Route path="/register" element={<Register/>}></Route>
        <Route path="/profile" element={<SidebarLayout><Profile /></SidebarLayout>} />
        <Route path="/change-password" element={<SidebarLayout><ChangePassword /></SidebarLayout>} />
        <Route path="/dashboard" element={<SidebarLayout><Dashboard /></SidebarLayout>} />
        <Route path="/create" element={<SidebarLayout><CreateNote /></SidebarLayout>} />
        <Route path="/edit/:id" element={<SidebarLayout><EditNote /></SidebarLayout>} />
        <Route path="/view/:id" element={<SidebarLayout><ViewNote /></SidebarLayout>} />
        <Route path="/trash" element={<SidebarLayout><TrashPage /></SidebarLayout>} />
      </Routes>
    </div>
  );
};

function App() {
  return (
    <Router>
      <AppContent />
    </Router>
  );
}

export default App;
