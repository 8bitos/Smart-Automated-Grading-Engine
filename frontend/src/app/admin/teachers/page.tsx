'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import AdminRouteGuard from '@/components/AdminRouteGuard';
import { useRouter } from 'next/navigation';

// Define the Teacher type based on expected data
interface Teacher {
  id: number;
  nama_lengkap: string;
  username: string;
  email: string;
}

export default function ManageTeachersPage() {
  const { user, token } = useAuth();
  const router = useRouter();

  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [isFetching, setIsFetching] = useState(true);
  const [fetchError, setFetchError] = useState<string | null>(null);

  const [showAddForm, setShowAddForm] = useState(false);
  const [newTeacherName, setNewTeacherName] = useState('');
  const [newTeacherUsername, setNewTeacherUsername] = useState('');
  const [newTeacherEmail, setNewTeacherEmail] = useState('');
  const [newTeacherPassword, setNewTeacherPassword] = useState('');
  const [addFormError, setAddFormError] = useState<string | null>(null);
  const [addFormLoading, setAddFormLoading] = useState(false);


  // Fetch all teachers
  const fetchTeachers = async () => {
    setIsFetching(true);
    setFetchError(null);
    try {
      const response = await fetch('http://localhost:8080/api/admin/teachers', {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      setTeachers(data || []);
    } catch (error: any) {
      setFetchError(error.message);
    } finally {
      setIsFetching(false);
    }
  };

  // Handle adding a new teacher
  const handleAddTeacher = async (e: React.FormEvent) => {
    e.preventDefault();
    setAddFormLoading(true);
    setAddFormError(null);

    try {
       const response = await fetch('http://localhost:8080/api/admin/teachers', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          nama_lengkap: newTeacherName,
          username: newTeacherUsername,
          email: newTeacherEmail,
          password: newTeacherPassword,
        }),
      });

      const responseData = await response.json();

      if (!response.ok) {
        throw new Error(responseData.message || 'Gagal menambahkan guru.');
      }

      // Reset form and refetch teachers
      setNewTeacherName('');
      setNewTeacherUsername('');
      setNewTeacherEmail('');
      setNewTeacherPassword('');
      setShowAddForm(false);
      fetchTeachers(); // Refetch the list to show the new teacher
    } catch (error: any) {
      setAddFormError(error.message);
    } finally {
      setAddFormLoading(false);
    }
  };

    // Handle deleting a teacher
  const handleDeleteTeacher = async (teacherId: number, teacherName: string) => {
    if (!window.confirm(`Apakah Anda yakin ingin menghapus guru "${teacherName}"?`)) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/api/admin/teachers/${teacherId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Gagal menghapus guru.');
      }
      
      // Filter out the deleted teacher from the state
      setTeachers(teachers.filter(t => t.id !== teacherId));

    } catch (error: any) {
      alert(`Error: ${error.message}`);
    }
  };

  // Fetch teachers on component mount
  useEffect(() => {
    if (token) {
      fetchTeachers();
    }
  }, [token]);

  return (
    <AdminRouteGuard>
      <div className="max-w-4xl mx-auto bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6">
          <h1 className="text-3xl font-bold mb-6 text-center text-gray-800 dark:text-white">Manajemen Guru</h1>

          {/* Tombol Tambah Guru */}
          <div className="mb-6 text-right">
            <button
              onClick={() => setShowAddForm(!showAddForm)}
              className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition-colors"
            >
              {showAddForm ? 'Sembunyikan Form' : 'Tambah Guru Baru'}
            </button>
          </div>

          {/* Form Tambah Guru */}
          {showAddForm && (
            <div className="mb-8 p-6 bg-gray-50 dark:bg-gray-700 rounded-lg shadow-inner">
              <h2 className="text-xl font-semibold mb-4 text-gray-800 dark:text-white">Form Tambah Guru</h2>
              <form onSubmit={handleAddTeacher} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300" htmlFor="newTeacherName">Nama Lengkap</label>
                  <input
                    type="text"
                    id="newTeacherName"
                    value={newTeacherName}
                    onChange={(e) => setNewTeacherName(e.target.value)}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm p-2 bg-white dark:bg-gray-600 text-gray-900 dark:text-white"
                    required
                  />
                </div>
                 <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300" htmlFor="newTeacherUsername">Username</label>
                  <input
                    type="text"
                    id="newTeacherUsername"
                    value={newTeacherUsername}
                    onChange={(e) => setNewTeacherUsername(e.target.value)}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm p-2 bg-white dark:bg-gray-600 text-gray-900 dark:text-white"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300" htmlFor="newTeacherEmail">Email</label>
                  <input
                    type="email"
                    id="newTeacherEmail"
                    value={newTeacherEmail}
                    onChange={(e) => setNewTeacherEmail(e.target.value)}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm p-2 bg-white dark:bg-gray-600 text-gray-900 dark:text-white"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300" htmlFor="newTeacherPassword">Password</label>
                  <input
                    type="password"
                    id="newTeacherPassword"
                    value={newTeacherPassword}
                    onChange={(e) => setNewTeacherPassword(e.target.value)}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm p-2 bg-white dark:bg-gray-600 text-gray-900 dark:text-white"
                    required
                  />
                </div>
                {addFormError && <p className="text-center text-sm text-red-500">{addFormError}</p>}
                <button
                  type="submit"
                  disabled={addFormLoading}
                  className="w-full bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded-lg transition-colors disabled:bg-green-300"
                >
                  {addFormLoading ? 'Menambahkan...' : 'Tambahkan Guru'}
                </button>
              </form>
            </div>
          )}

          {isFetching ? (
            <p className="text-center">Memuat daftar guru...</p>
          ) : fetchError ? (
            <p className="text-center text-red-500">Error: {fetchError}</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full bg-white dark:bg-gray-800 rounded-lg">
                <thead>
                  <tr className="bg-gray-200 dark:bg-gray-700">
                    <th className="py-3 px-4 text-left text-sm font-medium text-gray-600 dark:text-gray-300">Nama Lengkap</th>
                    <th className="py-3 px-4 text-left text-sm font-medium text-gray-600 dark:text-gray-300">Username</th>
                    <th className="py-3 px-4 text-left text-sm font-medium text-gray-600 dark:text-gray-300">Email</th>
                    <th className="py-3 px-4 text-left text-sm font-medium text-gray-600 dark:text-gray-300">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {teachers.length > 0 ? (
                    teachers.map((teacher) => (
                      <tr key={teacher.id} className="border-b dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700">
                        <td className="py-3 px-4">{teacher.nama_lengkap}</td>
                        <td className="py-3 px-4">{teacher.username}</td>
                        <td className="py-3 px-4">{teacher.email}</td>
                        <td className="py-3 px-4">
                          <button
                            onClick={() => handleDeleteTeacher(teacher.id, teacher.nama_lengkap)}
                            className="bg-red-500 hover:bg-red-600 text-white font-bold py-1 px-3 rounded text-xs transition-colors"
                          >
                            Hapus
                          </button>
                        </td>
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <td colSpan={4} className="py-3 px-4 text-center text-gray-500 dark:text-gray-400">Belum ada guru.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </div>
    </AdminRouteGuard>
  );
}