import { NavLink } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const links = [
  { to: '/dashboard', label: 'Dashboard', icon: '📊' },
  { to: '/upload', label: 'New Deck', icon: '✨' },
];

export function Sidebar() {
  const { user, logout } = useAuth();

  return (
    <aside className="flex h-full w-64 flex-col border-r border-slate-200 bg-white">
      <div className="border-b border-slate-100 px-6 py-5">
        <NavLink to="/dashboard" className="flex items-center gap-2">
          <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-brand-600 text-lg text-white">
            D
          </span>
          <span className="text-lg font-bold text-slate-900">DeckForge</span>
        </NavLink>
      </div>

      <nav className="flex-1 space-y-1 p-4">
        {links.map((link) => (
          <NavLink
            key={link.to}
            to={link.to}
            className={({ isActive }) =>
              `flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition ${
                isActive
                  ? 'bg-brand-50 text-brand-700'
                  : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
              }`
            }
          >
            <span>{link.icon}</span>
            {link.label}
          </NavLink>
        ))}
      </nav>

      <div className="border-t border-slate-100 p-4">
        <p className="truncate text-sm font-medium text-slate-900">{user?.full_name}</p>
        <p className="truncate text-xs text-slate-500">{user?.email}</p>
        <button
          onClick={logout}
          className="mt-3 w-full rounded-lg px-3 py-2 text-left text-sm text-slate-600 hover:bg-slate-50"
        >
          Sign out
        </button>
      </div>
    </aside>
  );
}
