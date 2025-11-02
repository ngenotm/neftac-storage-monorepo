'use client';
import { useState, useEffect } from 'react';
import axios from 'axios';

export default function Dashboard() {
  const [buckets, setBuckets] = useState<string[]>([]);
  const [selected, setSelected] = useState('');
  const [files, setFiles] = useState<string[]>([]);
  const token = typeof window !== 'undefined' ? localStorage.getItem('jwt') || '' : '';

  useEffect(() => {
    if (!token) return;
    axios.get('http://localhost:8080/v1/buckets', {
      headers: { Authorization: `Bearer ${token}` }
    }).then(r => setBuckets(r.data));
  }, [token]);

  const listFiles = (bucket: string) => {
    setSelected(bucket);
    axios.get(`http://localhost:8080/v1/buckets/${bucket}/objects`, {
      headers: { Authorization: `Bearer ${token}` }
    }).then(r => setFiles(r.data));
  };

  const upload = async (e: any) => {
    const file = e.target.files[0];
    await axios.put(`http://localhost:8080/v1/buckets/${selected}/objects/${file.name}`, file, {
      headers: { Authorization: `Bearer ${token}` }
    });
    listFiles(selected);
  };

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <h1 className="text-4xl font-bold text-gray-800 mb-8">Neftac Object Storage</h1>
      <div className="flex gap-6">
        <div className="w-64 bg-white p-4 rounded-lg shadow">
          <h2 className="font-semibold text-lg mb-3">Buckets</h2>
          {buckets.map(b => (
            <div key={b} onClick={() => listFiles(b)} className="p-3 hover:bg-blue-50 cursor-pointer rounded mb-1">
              {b}
            </div>
          ))}
        </div>
        <div className="flex-1 bg-white p-6 rounded-lg shadow">
          <h2 className="text-2xl font-semibold mb-4">{selected || 'Select a bucket'}</h2>
          {selected && (
            <>
              <input type="file" onChange={upload} className="mb-6 p-2 border rounded" />
              <ul className="space-y-2">
                {files.map(f => <li key={f} className="text-blue-600 hover:underline">{f}</li>)}
              </ul>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
