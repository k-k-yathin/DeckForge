import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';
import type { User } from '../types';
import * as authApi from '../api/auth';
import { getErrorMessage } from '../api/client';

interface AuthContextValue {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, fullName: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

/**
 * AuthProvider stores JWT + user in React state and localStorage.
 * Wrap the app so any component can call useAuth().
 */
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Restore session on page refresh
  useEffect(() => {
    const savedToken = localStorage.getItem('deckforge_token');
    const savedUser = localStorage.getItem('deckforge_user');
    if (savedToken && savedUser) {
      setToken(savedToken);
      try {
        setUser(JSON.parse(savedUser));
      } catch {
        localStorage.removeItem('deckforge_user');
      }
    }
    setIsLoading(false);
  }, []);

  const persist = useCallback((authUser: User, authToken: string) => {
    setUser(authUser);
    setToken(authToken);
    localStorage.setItem('deckforge_token', authToken);
    localStorage.setItem('deckforge_user', JSON.stringify(authUser));
  }, []);

  const login = useCallback(
    async (email: string, password: string) => {
      const data = await authApi.login(email, password);
      persist(data.user, data.token);
    },
    [persist]
  );

  const register = useCallback(
    async (email: string, password: string, fullName: string) => {
      const data = await authApi.register(email, password, fullName);
      persist(data.user, data.token);
    },
    [persist]
  );

  const logout = useCallback(() => {
    setUser(null);
    setToken(null);
    localStorage.removeItem('deckforge_token');
    localStorage.removeItem('deckforge_user');
  }, []);

  const value = useMemo(
    () => ({ user, token, isLoading, login, register, logout }),
    [user, token, isLoading, login, register, logout]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}

// Re-export for forms that need error messages
export { getErrorMessage };
