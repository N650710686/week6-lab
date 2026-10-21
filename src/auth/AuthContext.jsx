import { createContext, useContext, useState } from 'react';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem('moviehub.token'));
  const [member, setMember] = useState(() => {
    const raw = localStorage.getItem('moviehub.member');
    try {
      return raw ? JSON.parse(raw) : null;
    } catch {
      return null;
    }
  });

  const login = (data) => {
    setToken(data.token);
    setMember(data.member);
    localStorage.setItem('moviehub.token', data.token);
    localStorage.setItem('moviehub.member', JSON.stringify(data.member));
  };

  const register = login;

  const logout = () => {
    setToken(null);
    setMember(null);
    localStorage.removeItem('moviehub.token');
    localStorage.removeItem('moviehub.member');
  };

  const isAdmin = member?.role === 'admin';

  const value = { token, member, isLoggedIn: Boolean(token), isAdmin, login, register, logout };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return useContext(AuthContext);
}