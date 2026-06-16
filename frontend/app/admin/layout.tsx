'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Navbar from '@/components/Navbar';
import Link from 'next/link';

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const [isAdmin, setIsAdmin] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    const userStr = localStorage.getItem('user');
    
    if (!token) {
      router.push('/login');
      return;
    }

    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        if (user.role === 'admin') {
          setIsAdmin(true);
        } else {
          router.push('/');
        }
      } catch (error) {
        router.push('/');
      }
    }
    setLoading(false);
  }, [router]);

  if (loading) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center">
          <div className="text-black">Loading...</div>
        </div>
      </>
    );
  }

  if (!isAdmin) {
    return null;
  }

  return (
    <>
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex gap-8">
          {/* Sidebar */}
          <div className="w-64 bg-white rounded-xl shadow-md p-4 h-fit">
            <h2 className="font-bold text-black mb-4">Admin Menu</h2>
            <nav className="space-y-2">
              <Link href="/admin" className="block px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-100">
                Dashboard
              </Link>
              <Link href="/admin/products" className="block px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-100">
                Products
              </Link>
              <Link href="/admin/categories" className="block px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-100">
                Categories
              </Link>
              <Link href="/admin/orders" className="block px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-100">
                Orders
              </Link>
            </nav>
          </div>

          {/* Main Content */}
          <div className="flex-1">
            {children}
          </div>
        </div>
      </div>
    </>
  );
}
