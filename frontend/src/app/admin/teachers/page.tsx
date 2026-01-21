"use client";

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import AdminRouteGuard from '@/components/AdminRouteGuard';

interface Teacher {
  id: string;
  nama_lengkap: string;
  email: string;
}

export default function AdminTeachersPage() {
  const { token, user, isLoading } = useAuth();
  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [fetchError, setFetchError] = useState('');
  const [isFetching, setIsFetching] = useState(true);

  // State for Add Teacher Form
  const [showAddForm, setShowAddForm] = useState(false);
  const [newTeacherName, setNewTeacherName] = useState('');
  const [newTeacherEmail, setNewTeacherEmail] = useState('');
  const [newTeacherPassword, setNewTeacherPassword] = useState('');
  const [addFormLoading, setAddFormLoading] = useState(false);
  const [addFormError, setAddFormError] = useState('');

  // Fungsi untuk mengambil daftar guru
  const fetchTeachers = async () => {
    if (!token) {
      setFetchError('Tidak ada token otentikasi.');
      setIsFetching(false);
      return;
    }

    setIsFetching(true);
    setFetchError('');
    try {
      const res = await fetch('http://localhost:8080/api/admin/teachers', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.message || 'Gagal mengambil daftar guru.');
      }
      setTeachers(data);
    } catch (err: any) {
      setFetchError(err.message);
    } finally {
      setIsFetching(false);
    }
  };

  // Fungsi untuk menambah guru baru
  const handleAddTeacher = async (e: React.FormEvent) => {
    e.preventDefault();
    setAddFormError('');
    setAddFormLoading(true);

    if (!token) {
      setAddFormError('Tidak ada token otentikasi untuk menambah guru.');
      setAddFormLoading(false);
      return;
    }

    try {
      const res = await fetch('http://localhost:8080/api/admin/teachers', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          nama_lengkap: newTeacherName,
          email: newTeacherEmail,
          password: newTeacherPassword,
        }),
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.message || 'Gagal menambah guru.');
      }

      alert('Guru berhasil ditambahkan!');
      setNewTeacherName('');
      setNewTeacherEmail('');
      setNewTeacherPassword('');
      setShowAddForm(false);
      fetchTeachers(); // Refresh daftar guru
    } catch (err: any) {
      setAddFormError(err.message);
    } finally {
      setAddFormLoading(false);
    }
  };

  // Fungsi untuk menghapus guru
  const handleDeleteTeacher = async (teacherId: string, teacherName: string) => {
    if (!token) {
      alert('Tidak ada token otentikasi untuk menghapus guru.');
      return;
    }

    if (!confirm(`Apakah Anda yakin ingin menghapus guru ${teacherName}?`)) {
      return;
    }

    try {
      const res = await fetch(`http://localhost:8080/api/admin/teachers/${teacherId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.message || 'Gagal menghapus guru.');
      }

      alert('Guru berhasil dihapus!');
      fetchTeachers(); // Refresh daftar guru
    } catch (err: any) {
      alert(`Error menghapus guru: ${err.message}`);
    }
  };


  useEffect(() => {
    if (!isLoading && user && user.peran === 'superadmin') {
      fetchTeachers();
    }
  }, [isLoading, user, token]); 

  // Tampilan halaman admin
  return (
    <AdminRouteGuard>
      <main className="min-h-screen p-8 bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-white">
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
                    <th className="py-3 px-4 text-left text-sm font-medium text-gray-600 dark:text-gray-300">Email</th>
                    <th className="py-3 px-4 text-left text-sm font-medium text-gray-600 dark:text-gray-300">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {teachers.length > 0 ? (
                    teachers.map((teacher) => (
                      <tr key={teacher.id} className="border-b dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700">
                        <td className="py-3 px-4">{teacher.nama_lengkap}</td>
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
                      <td colSpan={3} className="py-3 px-4 text-center text-gray-500 dark:text-gray-400">Belum ada guru.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </main>
    </AdminRouteGuard>
  );
}
