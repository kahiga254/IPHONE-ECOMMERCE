'use client';

import Link from 'next/link';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function Navbar() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const router = useRouter();
  
  useEffect(() => {
    setMounted(true);
    const token = localStorage.getItem('access_token');
    const userStr = localStorage.getItem('user');
    
    setIsLoggedIn(!!token);
    
    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        setIsAdmin(user.role === 'admin');
      } catch (error) {
        setIsAdmin(false);
      }
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    router.push('/login');
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/products?search=${encodeURIComponent(searchQuery.trim())}`);
      setSearchQuery('');
    }
  };

  if (!mounted) {
    return (
      <nav className="bg-white border-b border-gray-200 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="text-2xl font-bold">
              <span className="text-black">Nifnoc</span> <span className="text-blue-600">iMobile</span>
            </div>
          </div>
        </div>
      </nav>
    );
  }

  return (
    <nav className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-wrap items-center justify-between py-2">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2">
            <span className="text-2xl font-bold">
              <span className="text-black">Nifnoc</span>{' '}
              <span className="text-blue-600">iMobile</span>
            </span>
          </Link>

          {/* Search Bar - Desktop */}
          <form onSubmit={handleSearch} className="hidden md:flex flex-1 max-w-md mx-8">
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search products..."
              className="w-full px-4 py-2 border border-gray-300 rounded-l-lg focus:ring-2 focus:ring-black focus:border-transparent text-black outline-none"
            />
            <button
              type="submit"
              className="bg-black text-white px-4 py-2 rounded-r-lg hover:bg-gray-800 transition"
            >
              Search
            </button>
          </form>

          {/* Desktop Menu */}
          <div className="hidden md:flex items-center space-x-6">
            <Link href="/" className="text-gray-800 hover:text-blue-600">Home</Link>
            <Link href="/products" className="text-gray-800 hover:text-blue-600">Products</Link>
            {isLoggedIn && (
              <Link href="/wishlist" className="text-gray-800 hover:text-blue-600">Wishlist</Link>
            )}
            <Link href="/cart" className="text-gray-800 hover:text-blue-600">Cart</Link>
            {isLoggedIn && (
              <Link href="/profile" className="text-gray-800 hover:text-blue-600">Profile</Link>
            )}
            {isAdmin && (
              <Link href="/admin" className="text-gray-800 hover:text-blue-600 font-semibold">Admin</Link>
            )}
            {isLoggedIn ? (
              <button onClick={handleLogout} className="border-2 border-black text-black px-4 py-1 rounded-full font-medium hover:bg-black hover:text-white">
                Logout
              </button>
            ) : (
              <>
                <Link href="/login" className="border-2 border-black text-black px-4 py-1 rounded-full font-medium hover:bg-black hover:text-white">
                  Login
                </Link>
                <Link href="/register" className="bg-black text-white px-4 py-1 rounded-full font-medium hover:bg-gray-800">
                  Register
                </Link>
              </>
            )}
          </div>

          {/* Mobile Menu Button */}
          <button className="md:hidden" onClick={() => setIsMenuOpen(!isMenuOpen)}>
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        </div>

        {/* Mobile Search */}
        <div className="md:hidden pb-3">
          <form onSubmit={handleSearch} className="flex">
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search products..."
              className="flex-1 px-4 py-2 border border-gray-300 rounded-l-lg focus:ring-2 focus:ring-black focus:border-transparent text-black outline-none"
            />
            <button
              type="submit"
              className="bg-black text-white px-4 py-2 rounded-r-lg hover:bg-gray-800 transition"
            >
              Search
            </button>
          </form>
        </div>

        {/* Mobile Menu */}
        {isMenuOpen && (
          <div className="md:hidden py-4 border-t border-gray-200">
            <div className="flex flex-col space-y-3">
              <Link href="/" className="text-gray-800 hover:text-blue-600">Home</Link>
              <Link href="/products" className="text-gray-800 hover:text-blue-600">Products</Link>
              {isLoggedIn && (
                <Link href="/wishlist" className="text-gray-800 hover:text-blue-600">Wishlist</Link>
              )}
              <Link href="/cart" className="text-gray-800 hover:text-blue-600">Cart</Link>
              {isLoggedIn && (
                <Link href="/profile" className="text-gray-800 hover:text-blue-600">Profile</Link>
              )}
              {isAdmin && (
                <Link href="/admin" className="text-gray-800 hover:text-blue-600 font-semibold">Admin</Link>
              )}
              {isLoggedIn ? (
                <button onClick={handleLogout} className="border-2 border-black text-black px-4 py-1 rounded-full w-full text-center">
                  Logout
                </button>
              ) : (
                <div className="flex space-x-3">
                  <Link href="/login" className="border-2 border-black text-black px-4 py-1 rounded-full flex-1 text-center">Login</Link>
                  <Link href="/register" className="bg-black text-white px-4 py-1 rounded-full flex-1 text-center">Register</Link>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </nav>
  );
}
