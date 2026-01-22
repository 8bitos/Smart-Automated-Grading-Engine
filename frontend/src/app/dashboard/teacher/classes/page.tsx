"use client";

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import TeacherRouteGuard from '@/components/TeacherRouteGuard';

interface Class {
  id: string;
  nama_kelas: string;
  deskripsi: string;
  created_at: string;
}

export default function TeacherClassesPage() {
  const { token } = useAuth();
  const [classes, setClasses] = useState<Class[]>([]);
  const [fetchError, setFetchError] = useState('');
  const [isFetching, setIsFetching] = useState(true);

  // State for Add Class Form
  const [showAddForm, setShowAddForm] = useState(false);
  const [newClassName, setNewClassName] = useState('');
  const [newClassDesc, setNewClassDesc] = useState('');
  const [addFormLoading, setAddFormLoading] = useState(false);
  const [addFormError, setAddFormError] = useState('');

  const fetchClasses = async () => {
    if (!token) return;
    setIsFetching(true);
    try {
      const res = await fetch('http://localhost:8080/api/classes', {
        headers: { 'Authorization': `Bearer ${token}` },
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.message || 'Gagal mengambil data kelas.');
      setClasses(data || []);
    } catch (err: any) {
      setFetchError(err.message);
    } finally {
      setIsFetching(false);
    }
  };

  const handleAddClass = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;
    setAddFormLoading(true);
    setAddFormError('');
    try {
      const res = await fetch('http://localhost:8080/api/classes', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          nama_kelas: newClassName,
          deskripsi: newClassDesc,
        }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.message || 'Gagal membuat kelas.');
      
      alert('Kelas berhasil dibuat!');
      setShowAddForm(false);
      setNewClassName('');
      setNewClassDesc('');
      fetchClasses(); // Refresh
    } catch (err: any) {
      setAddFormError(err.message);
    } finally {
      setAddFormLoading(false);
    }
  };

  useEffect(() => {
    fetchClasses();
  }, [token]);

  return (
    <TeacherRouteGuard>
        <div className="max-w-4xl mx-auto">
          <h1 className="text-3xl font-bold mb-6 text-gray-800 dark:text-white">Manajemen Kelas Anda</h1>

          <div className="mb-6 text-right">
            <button
              onClick={() => setShowAddForm(!showAddForm)}
              className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg"
            >
              {showAddForm ? 'Batal' : '+ Buat Kelas Baru'}
            </button>
          </div>

          {showAddForm && (
            <div className="mb-8 p-6 bg-white dark:bg-gray-800 rounded-lg shadow-md">
              <h2 className="text-xl font-semibold mb-4">Form Kelas Baru</h2>
              <form onSubmit={handleAddClass} className="space-y-4">
                <div>
                  <label htmlFor="newClassName" className="block text-sm font-medium">Nama Kelas</label>
                  <input
                    id="newClassName"
                    type="text"
                    value={newClassName}
                    onChange={(e) => setNewClassName(e.target.value)}
                    className="mt-1 block w-full rounded-md p-2 bg-gray-100 dark:bg-gray-700"
                    required
                  />
                </div>
                <div>
                  <label htmlFor="newClassDesc" className="block text-sm font-medium">Deskripsi</label>
                  <textarea
                    id="newClassDesc"
                    value={newClassDesc}
                    onChange={(e) => setNewClassDesc(e.target.value)}
                    className="mt-1 block w-full rounded-md p-2 bg-gray-100 dark:bg-gray-700"
                    rows={3}
                  />
                </div>
                {addFormError && <p className="text-red-500">{addFormError}</p>}
                <button type="submit" disabled={addFormLoading} className="w-full bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded-lg disabled:bg-green-400">
                  {addFormLoading ? 'Menyimpan...' : 'Simpan Kelas'}
                </button>
              </form>
            </div>
          )}

          <h2 className="text-2xl font-semibold mb-4">Daftar Kelas</h2>
          {isFetching ? <p>Loading...</p> : fetchError ? <p className="text-red-500">{fetchError}</p> : (
            <div className="space-y-4">
              {classes.length > 0 ? classes.map(c => (
                <div key={c.id} className="p-4 bg-white dark:bg-gray-800 border dark:border-gray-700 rounded-lg shadow-sm hover:shadow-md transition-shadow">
                  <h3 className="font-bold text-lg">{c.nama_kelas}</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">{c.deskripsi}</p>
                </div>
              )) : <p>Anda belum memiliki kelas.</p>}
            </div>
          )}
        </div>
    </TeacherRouteGuard>
  );
}
