import { Link, NavLink } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

function Navbar() {
  const { isLoggedIn, isAdmin, member, logout } = useAuth();

  const linkClass = ({ isActive }) =>
    'text-sm font-medium transition ' + (isActive ? 'text-emerald-600' : 'text-slate-600 hover:text-slate-900');

  return (
    <header className="border-b border-emerald-100 bg-white">
      <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-3 md:px-6">
        <Link to="/" className="text-lg font-bold text-emerald-600">
          MovieHub
        </Link>
        <nav className="flex items-center gap-4">
          <NavLink to="/" className={linkClass}>หน้าแรก</NavLink>
          <NavLink to="/movies" className={linkClass}>หนังทั้งหมด</NavLink>
          {isLoggedIn && <NavLink to="/me/wishlist" className={linkClass}>อยากดู</NavLink>}
          <NavLink to="/about" className={linkClass}>เกี่ยวกับ</NavLink>
          {isAdmin && <NavLink to="/admin" className={linkClass}>หลังบ้าน</NavLink>}
          {isLoggedIn ? (
            <div className="flex items-center gap-3">
              <span className="text-xs text-slate-500">{member?.displayName}</span>
              <button onClick={logout} className="rounded-lg border border-slate-200 px-3 py-1 text-xs text-slate-600 hover:bg-slate-50">
                ออกจากระบบ
              </button>
            </div>
          ) : (
            <Link to="/login" className="rounded-lg bg-emerald-500 px-3 py-1 text-xs font-medium text-white hover:bg-emerald-600">
              เข้าสู่ระบบ
            </Link>
          )}
        </nav>
      </div>
    </header>
  );
}

export default Navbar;