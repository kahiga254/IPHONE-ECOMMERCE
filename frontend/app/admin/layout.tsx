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
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

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
      
      {/* Mobile Sidebar Toggle */}
      <div className="md:hidden bg-gray-100 px-4 py-2 border-b">
        <button
          onClick={() => setIsSidebarOpen(!isSidebarOpen)}
          className="flex items-center gap-2 text-black"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          </svg>
          <span>Menu</span>
        </button>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-4 md:py-8">
        <div className="flex flex-col md:flex-row gap-4 md:gap-8">
          {/* Sidebar - Mobile overlay */}
          <div className={`
            ${isSidebarOpen ? 'block' : 'hidden'} 
            md:block
            fixed md:relative
            inset-0 md:inset-auto
            z-50 md:z-auto
            bg-black/50 md:bg-transparent
            md:w-64
          `}>
            <div className="
              bg-white rounded-xl shadow-md p-4 w-64 h-full md:h-auto
              absolute left-0 top-0 md:relative
              overflow-y-auto
            ">
              <div className="flex justify-between items-center mb-4 md:hidden">
                <h2 className="font-bold text-black">Admin Menu</h2>
                <button onClick={() => setIsSidebarOpen(false)} className="text-black">✕</button>
              </div>
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
                <Link href="/admin/reviews" className="block px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-100">
                  Reviews
                </Link>
              </nav>
            </div>
          </div>

          {/* Main Content */}
          <div className="flex-1 min-w-0 overflow-x-auto">
            {children}
          </div>
        </div>
      </div>
    </>
  );
}
