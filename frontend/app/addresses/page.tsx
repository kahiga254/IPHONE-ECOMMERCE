export const dynamic = 'force-dynamic';
'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Navbar from '@/components/Navbar';
import api from '@/services/api';
import toast from 'react-hot-toast';

interface Address {
  id: string;
  full_name: string;
  phone: string;
  street: string;
  county: string;
  town: string;
  is_default: boolean;
}

export default function AddressesPage() {
  const router = useRouter();
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    full_name: '',
    phone: '',
    street: '',
    county: '',
    town: '',
    is_default: false,
  });

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      router.push('/login');
      return;
    }
    fetchAddresses();
  }, []);

  const fetchAddresses = async () => {
    try {
      const response = await api.get('/addresses');
      setAddresses(response.data.data || []);
    } catch (error) {
      console.error('Failed to fetch addresses:', error);
      toast.error('Failed to load addresses');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      if (editingId) {
        await api.put(`/addresses/${editingId}`, formData);
        toast.success('Address updated successfully');
      } else {
        await api.post('/addresses', formData);
        toast.success('Address added successfully');
      }
      
      setShowForm(false);
      setEditingId(null);
      setFormData({
        full_name: '',
        phone: '',
        street: '',
        county: '',
        town: '',
        is_default: false,
      });
      fetchAddresses();
    } catch (error: any) {
      console.error('Failed to save address:', error);
      toast.error(error.response?.data?.error || 'Failed to save address');
    }
  };

  const handleEdit = (address: Address) => {
    setEditingId(address.id);
    setFormData({
      full_name: address.full_name,
      phone: address.phone,
      street: address.street,
      county: address.county,
      town: address.town,
      is_default: address.is_default,
    });
    setShowForm(true);
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this address?')) return;
    
    try {
      await api.delete(`/addresses/${id}`);
      toast.success('Address deleted successfully');
      fetchAddresses();
    } catch (error: any) {
      console.error('Failed to delete address:', error);
      toast.error(error.response?.data?.error || 'Failed to delete address');
    }
  };

  const handleSetDefault = async (id: string) => {
    try {
      await api.patch(`/addresses/${id}/default`);
      toast.success('Default address updated');
      fetchAddresses();
    } catch (error: any) {
      console.error('Failed to set default:', error);
      toast.error(error.response?.data?.error || 'Failed to set default address');
    }
  };

  if (loading) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center text-black">Loading addresses...</div>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 py-16">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-black">My Addresses</h1>
          {!showForm && (
            <button
              onClick={() => setShowForm(true)}
              className="bg-black text-white px-4 py-2 rounded-full hover:bg-gray-800 transition"
            >
              + Add New Address
            </button>
          )}
        </div>

        {/* Address Form */}
        {showForm && (
          <div className="bg-white rounded-xl shadow-md p-6 mb-8">
            <h2 className="text-xl font-bold text-black mb-4">
              {editingId ? 'Edit Address' : 'Add New Address'}
            </h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Full Name *
                  </label>
                  <input
                    type="text"
                    required
                    value={formData.full_name}
                    onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-black focus:border-transparent text-black"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Phone Number *
                  </label>
                  <input
                    type="tel"
                    required
                    value={formData.phone}
                    onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-black focus:border-transparent text-black"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Street Address *
                </label>
                <input
                  type="text"
                  required
                  value={formData.street}
                  onChange={(e) => setFormData({ ...formData, street: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-black focus:border-transparent text-black"
                />
              </div>

              <div className="grid md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    County *
                  </label>
                  <input
                    type="text"
                    required
                    value={formData.county}
                    onChange={(e) => setFormData({ ...formData, county: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-black focus:border-transparent text-black"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Town/City *
                  </label>
                  <input
                    type="text"
                    required
                    value={formData.town}
                    onChange={(e) => setFormData({ ...formData, town: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-black focus:border-transparent text-black"
                  />
                </div>
              </div>

              <div className="flex items-center gap-3">
                <input
                  type="checkbox"
                  id="is_default"
                  checked={formData.is_default}
                  onChange={(e) => setFormData({ ...formData, is_default: e.target.checked })}
                  className="w-4 h-4"
                />
                <label htmlFor="is_default" className="text-sm text-gray-700">
                  Set as default address
                </label>
              </div>

              <div className="flex gap-3 pt-2">
                <button
                  type="submit"
                  className="bg-black text-white px-6 py-2 rounded-full hover:bg-gray-800 transition"
                >
                  {editingId ? 'Update Address' : 'Save Address'}
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setShowForm(false);
                    setEditingId(null);
                    setFormData({
                      full_name: '',
                      phone: '',
                      street: '',
                      county: '',
                      town: '',
                      is_default: false,
                    });
                  }}
                  className="border border-gray-300 text-gray-700 px-6 py-2 rounded-full hover:bg-gray-100 transition"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}

        {/* Addresses List */}
        {addresses.length === 0 && !showForm ? (
          <div className="text-center py-16 bg-white rounded-xl shadow-md">
            <p className="text-black mb-4">You haven't saved any addresses yet</p>
            <button
              onClick={() => setShowForm(true)}
              className="bg-black text-white px-6 py-2 rounded-full hover:bg-gray-800"
            >
              Add Your First Address
            </button>
          </div>
        ) : (
          <div className="grid md:grid-cols-2 gap-6">
            {addresses.map((address) => (
              <div key={address.id} className="bg-white rounded-xl shadow-md p-6 relative">
                {address.is_default && (
                  <span className="absolute top-4 right-4 bg-green-100 text-green-800 text-xs px-2 py-1 rounded">
                    Default
                  </span>
                )}
                <h3 className="font-semibold text-lg text-black mb-2">{address.full_name}</h3>
                <p className="text-gray-600 text-sm mb-1">{address.street}</p>
                <p className="text-gray-600 text-sm mb-1">{address.town}, {address.county}</p>
                <p className="text-gray-600 text-sm mb-4">Phone: {address.phone}</p>
                
                <div className="flex gap-3">
                  <button
                    onClick={() => handleEdit(address)}
                    className="text-blue-600 hover:underline text-sm"
                  >
                    Edit
                  </button>
                  {!address.is_default && (
                    <>
                      <button
                        onClick={() => handleSetDefault(address.id)}
                        className="text-green-600 hover:underline text-sm"
                      >
                        Set Default
                      </button>
                      <button
                        onClick={() => handleDelete(address.id)}
                        className="text-red-600 hover:underline text-sm"
                      >
                        Delete
                      </button>
                    </>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  );
}
