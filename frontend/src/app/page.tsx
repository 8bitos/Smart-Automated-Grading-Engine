"use client";

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';

export default function Home() {
  const [message, setMessage] = useState('');
  const [loadingBackend, setLoadingBackend] = useState(true);
  const [errorBackend, setErrorBackend] = useState(null);
  const { user, isLoading, logout } = useAuth(); // Ambil user, isLoading, logout dari AuthContext

  useEffect(() => {
    fetch('http://localhost:8080/api/hello')
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(data => {
        setMessage(data.message);
        setLoadingBackend(false);
      })
      .catch(error => {
        setErrorBackend(error.message);
        setLoadingBackend(false);
      });
  }, []);

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24 bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-white">
      <div className="z-10 w-full max-w-5xl items-center justify-center font-mono text-sm lg:flex">
        <div className="text-center">
          <h1 className="text-4xl font-bold mb-4">Automatic Essay Scoring System</h1>
          <p className="mt-4 text-lg">
            {loadingBackend && 'Loading message from backend...'}
            {errorBackend && `Error from backend: ${errorBackend}`}
            {message && `Message from backend: "${message}"`}
          </p>

          <div className="mt-8 flex justify-center gap-4">
            {!isLoading && user ? ( // Jika user sudah login
              <>
                <p className="text-lg">Halo, {user.peran}!</p>
                {user.peran === 'superadmin' && ( // Tampilkan link Admin jika superadmin
                  <Link href="/admin/teachers" className="rounded-md bg-purple-600 px-4 py-2 text-white hover:bg-purple-700">
                    Admin Dashboard
                  </Link>
                )}
                <button
                  onClick={logout}
                  className="rounded-md bg-red-600 px-4 py-2 text-white hover:bg-red-700"
                >
                  Logout
                </button>
              </>
            ) : ( // Jika belum login
              <>
                <Link href="/login" className="rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700">
                  Login
                </Link>
                <Link href="/register" className="rounded-md bg-gray-600 px-4 py-2 text-white hover:bg-gray-700">
                  Register
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
