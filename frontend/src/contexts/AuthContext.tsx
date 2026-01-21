"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { jwtDecode } from 'jwt-decode';

// Tipe untuk data pengguna yang diambil dari token
interface User {
  user_id: string;
  peran: 'student' | 'teacher' | 'superadmin';
}

// Tipe untuk nilai yang disediakan oleh AuthContext
interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (token: string) => void;
  logout: () => void;
  isLoading: boolean;
}

// Membuat context dengan nilai default
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Props untuk AuthProvider
interface AuthProviderProps {
  children: ReactNode;
}

// Komponen Provider
export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Cek token saat komponen pertama kali dimuat
    try {
      const storedToken = localStorage.getItem('token');
      if (storedToken) {
        const decodedToken = jwtDecode<User>(storedToken);
        // Cek apakah token masih valid (belum expired)
        // Note: jwt-decode tidak memvalidasi signature, hanya membaca payload.
        // Validasi expiry adalah praktik yang baik.
        // const isExpired = decodedToken.exp * 1000 < Date.now();
        // if (isExpired) {
        //   logout();
        // } else {
          setUser(decodedToken);
          setToken(storedToken);
        // }
      }
    } catch (error) {
      console.error("Gagal memproses token:", error);
      // Jika token tidak valid, bersihkan
      localStorage.removeItem('token');
    } finally {
      setIsLoading(false);
    }
  }, []);

  const login = (newToken: string) => {
    try {
      const decodedToken = jwtDecode<User>(newToken);
      localStorage.setItem('token', newToken);
      setToken(newToken);
      setUser(decodedToken);
    } catch (error) {
      console.error("Gagal menyimpan atau decode token login:", error);
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
    setToken(null);
  };

  const value = { user, token, login, logout, isLoading };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// Custom hook untuk menggunakan AuthContext
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth harus digunakan di dalam AuthProvider');
  }
  return context;
};
