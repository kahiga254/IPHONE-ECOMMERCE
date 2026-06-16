'use client';

import { useEffect, useState } from 'react';
import api from '@/services/api';
import toast from 'react-hot-toast';

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
    { label: 'Total Products', value: stats.totalProducts, color: 'bg-blue-500' },
    { label: 'Total Orders', value: stats.totalOrders, color: 'bg-green-500' },
    { label: 'Total Users', value: stats.totalUsers, color: 'bg-purple-500' },
    { label: 'Revenue', value: `KSh ${stats.totalRevenue.toLocaleString()}`, color: 'bg-yellow-500' },
  ];

  if (loading) {
    return <div className="text-center py-16 text-black">Loading dashboard...</div>;
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-black mb-8">Dashboard</h1>

      {/* Stats Cards */}
      <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {statCards.map((card, idx) => (
          <div key={idx} className="bg-white rounded-xl shadow-md p-6">
            <p className="text-gray-600 text-sm mb-2">{card.label}</p>
            <p className="text-2xl font-bold text-black">{card.value}</p>
          </div>
        ))}
      </div>

      {/* Recent Orders */}
      <div className="bg-white rounded-xl shadow-md p-6">
        <h2 className="text-xl font-bold text-black mb-4">Recent Orders</h2>
        {recentOrders.length === 0 ? (
          <p className="text-gray-500 text-center py-8">No orders yet</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="border-b border-gray-200">
                <tr className="text-left">
                  <th className="pb-3 text-gray-600 font-medium">Order #</th>
                  <th className="pb-3 text-gray-600 font-medium">Date</th>
                  <th className="pb-3 text-gray-600 font-medium">Total</th>
                  <th className="pb-3 text-gray-600 font-medium">Status</th>
                </tr>
              </thead>
              <tbody>
                {recentOrders.map((order: any) => (
                  <tr key={order.id} className="border-b border-gray-100">
                    <td className="py-3 text-black">{order.order_number}</td>
                    <td className="py-3 text-gray-600">{new Date(order.created_at).toLocaleDateString()}</td>
                    <td className="py-3 text-black">KSh {order.total_amount?.toLocaleString()}</td>
                    <td className="py-3">
                      <span className="px-2 py-1 bg-yellow-100 text-yellow-800 rounded text-xs">
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
