'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Navbar from '@/components/Navbar';
import api from '@/services/api';
import toast from 'react-hot-toast';

interface OrderItem {
  id: string;
  product_name: string;
  quantity: number;
  price: number;
}

interface Order {
  id: string;
  order_number: string;
  total_amount: number;
  status: string;
  payment_status: string;
  shipping_address: string;
  created_at: string;
  items: OrderItem[];
}

export default function OrderDetailPage() {
  const { id } = useParams();
  const router = useRouter();
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchOrder();
  }, [id]);

  const fetchOrder = async () => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      router.push('/login');
      return;
    }

    try {
      const response = await api.get(`/orders/${id}`);
      setOrder(response.data.data);
    } catch (error: any) {
      console.error('Failed to fetch order:', error);
      toast.error(error.response?.data?.error || 'Failed to load order');
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

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-KE', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center text-black">Loading order details...</div>
      </>
    );
  }

  if (!order) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center text-black">Order not found</div>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 py-16">
        <Link href="/orders" className="text-blue-600 hover:underline mb-6 inline-block">
          ← Back to Orders
        </Link>

        <h1 className="text-3xl font-bold text-black mb-2">Order #{order.order_number}</h1>
        <p className="text-gray-600 mb-6">Placed on {formatDate(order.created_at)}</p>

        <div className="grid lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-6">
            {/* Order Items */}
            <div className="bg-white rounded-xl shadow-md p-6">
              <h2 className="text-xl font-bold text-black mb-4">Order Items</h2>
              <div className="space-y-3">
                {order.items?.map((item) => (
                  <div key={item.id} className="flex justify-between items-center py-3 border-b border-gray-100 last:border-0">
                    <div>
                      <p className="font-medium text-black">{item.product_name}</p>
                      <p className="text-sm text-gray-600">Quantity: {item.quantity}</p>
                    </div>
                    <p className="font-semibold text-black">KSh {item.price.toLocaleString()}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* Order Status */}
            <div className="bg-white rounded-xl shadow-md p-6">
              <h2 className="text-xl font-bold text-black mb-4">Order Status</h2>
              <div className="flex items-center gap-3 mb-4">
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(order.status)}`}>
                  {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                </span>
                <span className="text-gray-600">•</span>
                <span className="text-gray-600">Payment: {order.payment_status}</span>
              </div>
            </div>
          </div>

          {/* Order Summary */}
          <div className="bg-white rounded-xl shadow-md p-6 h-fit">
            <h2 className="text-xl font-bold text-black mb-4">Order Summary</h2>
            <div className="space-y-2 border-b border-gray-200 pb-4">
              <div className="flex justify-between text-black">
                <span>Subtotal</span>
                <span className="font-semibold">KSh {order.total_amount.toLocaleString()}</span>
              </div>
              <div className="flex justify-between text-black">
                <span>Shipping</span>
                <span className="font-semibold">KSh 200</span>
              </div>
            </div>
            <div className="flex justify-between font-bold text-xl mt-4 text-black">
              <span>Total</span>
              <span className="text-blue-600">KSh {order.total_amount.toLocaleString()}</span>
            </div>

            {order.shipping_address && (
              <div className="mt-6 pt-4 border-t border-gray-200">
                <h3 className="font-semibold text-black mb-2">Shipping Address</h3>
                <p className="text-gray-600 text-sm">{order.shipping_address}</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
