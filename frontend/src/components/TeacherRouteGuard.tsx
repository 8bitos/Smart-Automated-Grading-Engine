"use client";

import { ReactNode, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

interface TeacherRouteGuardProps {
  children: ReactNode;
}

const TeacherRouteGuard = ({ children }: TeacherRouteGuardProps) => {
  const { user, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading) {
      const isAuthorized = user?.peran === 'teacher' || user?.peran === 'superadmin';
      if (!isAuthorized) {
        alert('Akses ditolak. Anda harus login sebagai Guru atau Superadmin.');
        router.replace('/');
      }
    }
  }, [user, isLoading, router]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p>Loading user data...</p>
      </div>
    );
  }

  const isAuthorized = user?.peran === 'teacher' || user?.peran === 'superadmin';
  if (isAuthorized) {
    return <>{children}</>;
  }

  return null;
};

export default TeacherRouteGuard;
