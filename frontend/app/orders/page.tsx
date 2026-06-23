'use client';
export const dynamic = 'force-dynamic';


import Link from 'next/link';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Navbar from '@/components/Navbar';
import api from '@/services/api';
import toast from 'react-hot-toast';

interface Order {
  id: string;
  order_number: string;
  total_amount: number;
  status: string;
  payment_status: string;
  created_at: string;
}

export default function OrdersPage() {
  const router = useRouter();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      router.push('/login');
      return;
    }

    try {
      const response = await api.get('/orders');
      // The response has { success: true, data: { data: [], total: 0, ... } }
      const ordersData = response.data.data;
      
      if (ordersData && Array.isArray(ordersData.data)) {
        setOrders(ordersData.data);
      } else if (Array.isArray(ordersData)) {
        setOrders(ordersData);
      } else {
        console.error('Orders data is not an array:', ordersData);
        setOrders([]);
      }
    } catch (error: any) {
      console.error('Failed to fetch orders:', error);
      toast.error(error.response?.data?.error || 'Failed to load orders');
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'processing':
        return 'bg-blue-100 text-blue-800';
      case 'shipped':
        return 'bg-purple-100 text-purple-800';
      case 'delivered':
        return 'bg-green-100 text-green-800';
      case 'cancelled':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getPaymentStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'paid':
        return 'bg-green-100 text-green-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-KE', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  if (loading) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center text-black">Loading orders...</div>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 py-16">
        <h1 className="text-3xl font-bold text-black mb-8">My Orders</h1>

        {orders.length === 0 ? (
          <div className="text-center py-16">
            <p className="text-black mb-4">You haven't placed any orders yet</p>
            <Link href="/products" className="bg-black text-white px-6 py-2 rounded-full hover:bg-gray-800">
              Start Shopping
            </Link>
          </div>
        ) : (
          <div className="space-y-4">
            {orders.map((order) => (
              <div key={order.id} className="bg-white rounded-xl shadow-md overflow-hidden">
                <div className="p-6">
                  <div className="flex flex-wrap justify-between items-start mb-4">
                    <div>
                      <p className="text-sm text-gray-600">Order #{order.order_number}</p>
                      <p className="text-sm text-gray-600 mt-1">{formatDate(order.created_at)}</p>
                    </div>
                    <div className="text-right">
                      <p className="text-xl font-bold text-black">KSh {order.total_amount.toLocaleString()}</p>
                    </div>
                  </div>

                  <div className="flex flex-wrap gap-3 mb-4">
                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(order.status)}`}>
                      {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                    </span>
                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${getPaymentStatusColor(order.payment_status)}`}>
                      Payment: {order.payment_status.charAt(0).toUpperCase() + order.payment_status.slice(1)}
                    </span>
                  </div>

                  <Link
                    href={`/orders/${order.id}`}
                    className="inline-block text-blue-600 hover:underline text-sm font-medium"
                  >
                    View Details →
                  </Link>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  );
}