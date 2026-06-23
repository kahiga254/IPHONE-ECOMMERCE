export const dynamic = 'force-dynamic';
'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Navbar from '@/components/Navbar';
import api from '@/services/api';
import toast from 'react-hot-toast';

export default function WishlistPage() {
  const router = useRouter();
  const [wishlist, setWishlist] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchWishlist();
  }, []);

  const fetchWishlist = async () => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      router.push('/login');
      return;
    }

    try {
      const response = await api.get('/wishlist');
      setWishlist(response.data.data || []);
    } catch (error) {
      console.error('Failed to fetch wishlist:', error);
      toast.error('Failed to load wishlist');
    } finally {
      setLoading(false);
    }
  };

  const removeFromWishlist = async (variantId: string) => {
    try {
      await api.delete(`/wishlist/${variantId}`);
      toast.success('Removed from wishlist');
      fetchWishlist();
    } catch (error) {
      toast.error('Failed to remove');
    }
  };

  if (loading) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center">Loading...</div>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 py-16">
        <h1 className="text-3xl font-bold mb-8">My Wishlist</h1>

        {wishlist.length === 0 ? (
          <div className="text-center py-16">
            <p className="text-gray-500 mb-4">Your wishlist is empty</p>
            <Link href="/products" className="bg-black text-white px-6 py-2 rounded-full hover:bg-gray-800">
              Start Shopping
            </Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {wishlist.map((item: any) => (
              <div key={item.id} className="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg">
                <div className="p-4">
                  <h3 className="font-semibold text-lg mb-1">{item.variant?.sku || 'Product'}</h3>
                  <p className="text-gray-600 text-sm mb-2">
                    Color: {item.variant?.color} | Storage: {item.variant?.storage}
                  </p>
                  <p className="text-blue-600 font-bold text-xl mb-4">
                    KSh {item.variant?.price?.toLocaleString()}
                  </p>
                  <div className="flex gap-3">
                    <button className="flex-1 bg-black text-white py-2 rounded-full text-sm hover:bg-gray-800">
                      Add to Cart
                    </button>
                    <button 
                      onClick={() => removeFromWishlist(item.variant_id)}
                      className="px-4 py-2 border border-red-500 text-red-500 rounded-full text-sm hover:bg-red-500 hover:text-white transition"
                    >
                      Remove
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  );
}
