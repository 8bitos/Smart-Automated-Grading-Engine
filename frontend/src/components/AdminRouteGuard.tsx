"use client";

import { ReactNode, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

interface AdminRouteGuardProps {
  children: ReactNode;
}

const AdminRouteGuard = ({ children }: AdminRouteGuardProps) => {
  const { user, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    // Jika selesai loading dan ternyata user tidak ada atau bukan superadmin, redirect
    if (!isLoading && user?.peran !== 'superadmin') {
      alert('Akses ditolak. Anda bukan Superadmin.');
      router.replace('/'); // atau ke halaman login: router.replace('/login');
    }
  }, [user, isLoading, router]);

  // Selama loading, tampilkan spinner atau halaman loading
  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p>Loading...</p>
      </div>
    );
  }

  // Jika user adalah superadmin, tampilkan konten halaman
  if (user?.peran === 'superadmin') {
    return <>{children}</>;
  }

  // Fallback, seharusnya tidak akan pernah tercapai karena useEffect akan redirect
  return null;
};

export default AdminRouteGuard;
