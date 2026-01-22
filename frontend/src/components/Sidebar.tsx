"use client";

import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';

const Sidebar = () => {
  const { user } = useAuth();

  return (
    <aside className="w-64 flex-shrink-0 bg-white dark:bg-gray-800 p-4 border-r border-gray-200 dark:border-gray-700">
      <div className="flex items-center mb-8">
        {/* Placeholder Logo */}
        <div className="p-2 rounded-lg bg-gray-200 dark:bg-gray-700 mr-3">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M16.88 15.54L15.46 14.12L14.05 15.54L12 13.49L9.95 15.54L8.54 14.12L10.59 12L8.54 9.88L9.95 8.46L12 10.51L14.05 8.46L15.46 9.88L13.41 12L15.46 14.12L16.88 15.54M12 2C6.48 2 2 6.48 2 12C2 17.52 6.48 22 12 22C17.52 22 22 17.52 22 12C22 6.48 17.52 2 12 2"/></svg>
        </div>
        <h1 className="text-xl font-bold text-gray-900 dark:text-white">SAGE LMS</h1>
      </div>
      <nav className="flex flex-col space-y-2">
        <h2 className="px-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Menu</h2>
        
        {/* Tautan untuk Guru & Superadmin */}
        {(user?.peran === 'teacher' || user?.peran === 'superadmin') && (
          <Link href="/dashboard/teacher/classes" className="flex items-center px-4 py-2 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700">
            Manajemen Kelas
          </Link>
        )}

        {/* Tautan khusus Superadmin */}
        {user?.peran === 'superadmin' && (
          <>
            <h2 className="px-4 pt-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Admin</h2>
            <Link href="/admin/teachers" className="flex items-center px-4 py-2 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700">
              Manajemen Guru
            </Link>
          </>
        )}
        
        {/* Placeholder untuk link siswa */}
        {user?.peran === 'student' && (
           <p className="px-4 py-2 text-sm text-gray-500">Menu siswa akan ada di sini</p>
        )}

      </nav>
    </aside>
  );
};

export default Sidebar;
