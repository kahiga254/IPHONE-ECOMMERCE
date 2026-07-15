'use client';

import { useEffect, useState } from 'react';
import api from '@/services/api';
import toast from 'react-hot-toast';
import Link from 'next/link';

interface Stats {
  totalProducts: number;
  totalOrders: number;
  totalUsers: number;
  totalRevenue: number;
}

export default function AdminDashboard() {
  const [stats, setStats] = useState<Stats>({
    totalProducts: 0,
    totalOrders: 0,
    totalUsers: 0,
    totalRevenue: 0,
  });
  const [loading, setLoading] = useState(true);
  const [recentOrders, setRecentOrders] = useState([]);

  useEffect(() => {
    fetchStats();
    fetchRecentOrders();
  }, []);

  const fetchStats = async () => {
    try {
      const products = await api.get('/products');
      const orders = await api.get('/admin/orders');
      
      setStats({
        totalProducts: products.data.data?.data?.length || 0,
        totalOrders: orders.data.data?.data?.length || 0,
        totalUsers: 0,
        totalRevenue: 0,
      });
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchRecentOrders = async () => {
    try {
      const response = await api.get('/admin/orders?limit=5');
      setRecentOrders(response.data.data?.data || []);
    } catch (error) {
      console.error('Failed to fetch recent orders:', error);
    }
  };

  const statCards = [
    { label: 'Total Products', value: stats.totalProducts, color: 'bg-blue-50 border-blue-200', textColor: 'text-blue-700' },
    { label: 'Total Orders', value: stats.totalOrders, color: 'bg-green-50 border-green-200', textColor: 'text-green-700' },
    { label: 'Total Users', value: stats.totalUsers, color: 'bg-purple-50 border-purple-200', textColor: 'text-purple-700' },
    { label: 'Revenue', value: `KSh ${stats.totalRevenue.toLocaleString()}`, color: 'bg-yellow-50 border-yellow-200', textColor: 'text-yellow-700' },
  ];

  if (loading) {
    return <div className="text-center py-8 md:py-16 text-black">Loading dashboard...</div>;
  }

  return (
    <div>
      <h1 className="text-xl md:text-2xl lg:text-3xl font-bold text-black mb-4 md:mb-8">Dashboard</h1>

      {/* Stats Cards - Mobile responsive grid */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-6 mb-4 md:mb-8">
        {statCards.map((card, idx) => (
          <div key={idx} className={`${card.color} rounded-xl shadow-sm p-3 md:p-6 border`}>
            <p className="text-xs md:text-sm text-gray-600 mb-1">{card.label}</p>
            <p className={`text-base md:text-2xl font-bold ${card.textColor}`}>{card.value}</p>
          </div>
        ))}
      </div>

      {/* Recent Orders - Responsive */}
      <div className="bg-white rounded-xl shadow-md p-4 md:p-6 overflow-x-auto">
        <h2 className="text-base md:text-xl font-bold text-black mb-3 md:mb-4">Recent Orders</h2>
        {recentOrders.length === 0 ? (
          <p className="text-gray-500 text-center py-4 md:py-8">No orders yet</p>
        ) : (
          <div className="overflow-x-auto -mx-4 md:mx-0">
            <table className="w-full text-sm md:text-base">
              <thead className="border-b border-gray-200">
                <tr className="text-left">
                  <th className="pb-2 md:pb-3 px-2 md:px-0 text-gray-600 font-medium">Order #</th>
                  <th className="pb-2 md:pb-3 px-2 md:px-0 text-gray-600 font-medium hidden sm:table-cell">Date</th>
                  <th className="pb-2 md:pb-3 px-2 md:px-0 text-gray-600 font-medium text-right">Total</th>
                  <th className="pb-2 md:pb-3 px-2 md:px-0 text-gray-600 font-medium text-right">Status</th>
                </tr>
              </thead>
              <tbody>
                {recentOrders.map((order: any) => (
                  <tr key={order.id} className="border-b border-gray-100">
                    <td className="py-2 md:py-3 px-2 md:px-0 text-black font-mono text-xs md:text-sm">{order.order_number}</td>
                    <td className="py-2 md:py-3 px-2 md:px-0 text-gray-600 text-xs md:text-sm hidden sm:table-cell">
                      {new Date(order.created_at).toLocaleDateString()}
                    </td>
                    <td className="py-2 md:py-3 px-2 md:px-0 text-black text-right text-sm md:text-base">
                      KSh {order.total_amount?.toLocaleString()}
                    </td>
                    <td className="py-2 md:py-3 px-2 md:px-0 text-right">
                      <span className={`px-2 py-0.5 md:px-3 md:py-1 rounded text-xs font-medium ${
                        order.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                        order.status === 'processing' ? 'bg-blue-100 text-blue-800' :
                        order.status === 'shipped' ? 'bg-purple-100 text-purple-800' :
                        order.status === 'delivered' ? 'bg-green-100 text-green-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {order.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
