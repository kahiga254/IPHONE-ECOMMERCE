'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Navbar from '@/components/Navbar';
import api from '@/services/api';
import toast from 'react-hot-toast';

interface CartItem {
  id: string;
  name: string;
  price: number;
  quantity: number;
  image: string;
  variant: string;
}

interface Address {
  id: string;
  full_name: string;
  phone: string;
  street: string;
  county: string;
  town: string;
  is_default: boolean;
}

export default function CheckoutPage() {
  const router = useRouter();
  const [cartItems, setCartItems] = useState<CartItem[]>([]);
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [selectedAddress, setSelectedAddress] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [processing, setProcessing] = useState(false);
  const [paymentMethod, setPaymentMethod] = useState('mpesa');
  const [mpesaPhone, setMpesaPhone] = useState('');

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      router.push('/login');
      return;
    }
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const savedCart = localStorage.getItem('cart');
      if (savedCart) {
        const items = JSON.parse(savedCart);
        setCartItems(items);
      }

      const response = await api.get('/addresses');
      const addressData = response.data.data || [];
      setAddresses(addressData);
      
      const defaultAddr = addressData.find((addr: Address) => addr.is_default);
      if (defaultAddr) {
        setSelectedAddress(defaultAddr.id);
        setMpesaPhone(defaultAddr.phone);
      } else if (addressData.length > 0) {
        setSelectedAddress(addressData[0].id);
        setMpesaPhone(addressData[0].phone);
      }
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  };

  const subtotal = cartItems.reduce((sum, item) => sum + item.price * item.quantity, 0);
  const shipping = subtotal > 0 ? 200 : 0;
  const total = subtotal + shipping;

  const handlePlaceOrder = async () => {
    if (cartItems.length === 0) {
      toast.error('Your cart is empty');
      return;
    }

    if (!selectedAddress) {
      toast.error('Please select a shipping address');
      return;
    }

    if (paymentMethod === 'mpesa' && !mpesaPhone) {
      toast.error('Please enter M-Pesa phone number');
      return;
    }

    setProcessing(true);

    try {
      const token = localStorage.getItem('access_token');
      
      // Create order first
      const orderData = {
        items: cartItems.map(item => ({
          variant_id: item.id,
          quantity: item.quantity,
          price: item.price
        })),
        address_id: selectedAddress,
        subtotal: subtotal,
        shipping_fee: shipping,
        total: total,
        payment_method: paymentMethod
      };

      const orderResponse = await api.post('/orders', orderData);
      const order = orderResponse.data.data;
      
      if (paymentMethod === 'mpesa') {
        // Initiate M-Pesa payment
        const paymentResponse = await api.post('/payments/mpesa/stkpush', {
          order_id: order.id,
          phone_number: mpesaPhone,
          amount: total
        });
        
        if (paymentResponse.data.success) {
          toast.success('Payment initiated! Check your phone for M-Pesa prompt');
          // Clear cart
          localStorage.removeItem('cart');
          // Redirect to order confirmation
          router.push(`/orders/${order.id}?payment=pending`);
        } else {
          toast.error('Payment initiation failed');
        }
      } else {
        // Cash on delivery
        toast.success('Order placed successfully! You will pay on delivery');
        localStorage.removeItem('cart');
        router.push(`/orders/${order.id}`);
      }
    } catch (error: any) {
      console.error('Order failed:', error);
      toast.error(error.response?.data?.error || 'Failed to place order');
    } finally {
      setProcessing(false);
    }
  };

  if (loading) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center text-black">Loading checkout...</div>
      </>
    );
  }

  if (cartItems.length === 0) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center">
          <p className="text-black mb-4">Your cart is empty</p>
          <Link href="/products" className="bg-black text-white px-6 py-2 rounded-full hover:bg-gray-800">
            Continue Shopping
          </Link>
        </div>
      </>
    );
  }

  const selectedAddressObj = addresses.find(a => a.id === selectedAddress);

  return (
    <>
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 py-16">
        <h1 className="text-3xl font-bold text-black mb-8">Checkout</h1>

        <div className="grid lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-6">
            {/* Shipping Address */}
            <div className="bg-white rounded-xl shadow-md p-6">
              <h2 className="text-xl font-bold text-black mb-4">Shipping Address</h2>
              
              {addresses.length === 0 ? (
                <div className="text-center py-4">
                  <p className="text-gray-600 mb-3">No addresses saved</p>
                  <Link href="/addresses" className="text-blue-600 hover:underline">
                    Add New Address →
                  </Link>
                </div>
              ) : (
                <div className="space-y-3">
                  {addresses.map((addr) => (
                    <label key={addr.id} className="flex items-start gap-3 p-3 border border-gray-200 rounded-lg cursor-pointer hover:bg-gray-50">
                      <input
                        type="radio"
                        name="address"
                        value={addr.id}
                        checked={selectedAddress === addr.id}
                        onChange={(e) => {
                          setSelectedAddress(e.target.value);
                          setMpesaPhone(addr.phone);
                        }}
                        className="mt-1 cursor-pointer"
                      />
                      <div>
                        <p className="font-medium text-black">{addr.full_name}</p>
                        <p className="text-sm text-gray-600">{addr.street}</p>
                        <p className="text-sm text-gray-600">{addr.town}, {addr.county}</p>
                        <p className="text-sm text-gray-600">Phone: {addr.phone}</p>
                      </div>
                    </label>
                  ))}
                  <Link href="/addresses" className="text-blue-600 hover:underline text-sm inline-block mt-2">
                    + Add New Address
                  </Link>
                </div>
              )}
            </div>

            {/* Payment Method */}
            <div className="bg-white rounded-xl shadow-md p-6">
              <h2 className="text-xl font-bold text-black mb-4">Payment Method</h2>
              <div className="space-y-3">
                <label className="flex items-center gap-3 p-3 border border-gray-200 rounded-lg cursor-pointer hover:bg-gray-50">
                  <input
                    type="radio"
                    name="payment"
                    value="mpesa"
                    checked={paymentMethod === 'mpesa'}
                    onChange={(e) => setPaymentMethod(e.target.value)}
                    className="cursor-pointer"
                  />
                  <div>
                    <p className="font-medium text-black">M-Pesa</p>
                    <p className="text-sm text-gray-600">Pay via M-Pesa STK Push</p>
                  </div>
                </label>
                
                {paymentMethod === 'mpesa' && (
                  <div className="ml-8 mt-3">
                    <label className="block text-sm font-medium text-black mb-1">
                      M-Pesa Phone Number
                    </label>
                    <input
                      type="tel"
                      value={mpesaPhone}
                      onChange={(e) => setMpesaPhone(e.target.value)}
                      placeholder="0712345678"
                      className="w-full max-w-sm px-4 py-2 border border-gray-300 rounded-lg text-black"
                    />
                    <p className="text-xs text-gray-500 mt-1">You will receive a prompt on this phone</p>
                  </div>
                )}
                
                <label className="flex items-center gap-3 p-3 border border-gray-200 rounded-lg cursor-pointer hover:bg-gray-50">
                  <input
                    type="radio"
                    name="payment"
                    value="cash"
                    checked={paymentMethod === 'cash'}
                    onChange={(e) => setPaymentMethod(e.target.value)}
                    className="cursor-pointer"
                  />
                  <div>
                    <p className="font-medium text-black">Cash on Delivery</p>
                    <p className="text-sm text-gray-600">Pay when you receive the order</p>
                  </div>
                </label>
              </div>
            </div>

            {/* Order Items Summary */}
            <div className="bg-white rounded-xl shadow-md p-6">
              <h2 className="text-xl font-bold text-black mb-4">Order Items</h2>
              <div className="space-y-3">
                {cartItems.map((item) => (
                  <div key={item.id} className="flex justify-between items-center py-2 border-b border-gray-100">
                    <div>
                      <p className="font-medium text-black">{item.name}</p>
                      <p className="text-sm text-gray-600">Qty: {item.quantity} | {item.variant}</p>
                    </div>
                    <p className="font-semibold text-black">KSh {(item.price * item.quantity).toLocaleString()}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Order Summary */}
          <div className="bg-white rounded-xl shadow-md p-6 h-fit">
            <h2 className="text-xl font-bold text-black mb-4">Order Summary</h2>
            <div className="space-y-2 border-b border-gray-200 pb-4">
              <div className="flex justify-between text-black">
                <span>Subtotal</span>
                <span>KSh {subtotal.toLocaleString()}</span>
              </div>
              <div className="flex justify-between text-black">
                <span>Shipping</span>
                <span>KSh {shipping.toLocaleString()}</span>
              </div>
            </div>
            <div className="flex justify-between font-bold text-xl mt-4 text-black">
              <span>Total</span>
              <span className="text-blue-600">KSh {total.toLocaleString()}</span>
            </div>

            {selectedAddressObj && (
              <div className="mt-6 pt-4 border-t border-gray-200">
                <h3 className="font-semibold text-black mb-2">Deliver to:</h3>
                <p className="text-gray-600 text-sm">{selectedAddressObj.full_name}</p>
                <p className="text-gray-600 text-sm">{selectedAddressObj.street}</p>
                <p className="text-gray-600 text-sm">{selectedAddressObj.town}, {selectedAddressObj.county}</p>
              </div>
            )}

            <button
              onClick={handlePlaceOrder}
              disabled={processing || addresses.length === 0 || cartItems.length === 0}
              className={`w-full py-3 rounded-full mt-6 font-medium transition-all ${
                processing 
                  ? 'bg-gray-400 cursor-not-allowed' 
                  : 'bg-black text-white hover:bg-gray-800 cursor-pointer'
              }`}
            >
              {processing ? 'Processing...' : `Place Order • KSh ${total.toLocaleString()}`}
            </button>
          </div>
        </div>
      </div>
    </>
  );
}
