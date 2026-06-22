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

export default function CheckoutPage() {
  const router = useRouter();
  const [cartItems, setCartItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [processing, setProcessing] = useState(false);
  const [paymentMethod, setPaymentMethod] = useState('mpesa');
  const [mpesaPhone, setMpesaPhone] = useState('');
  const [isGuest, setIsGuest] = useState(true);
  const [guestInfo, setGuestInfo] = useState({
    full_name: '',
    email: '',
    phone: '',
    address: '',
    city: '',
    county: ''
  });

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const savedCart = localStorage.getItem('cart');
      if (savedCart) {
        const items = JSON.parse(savedCart);
        setCartItems(items);
      }
      
      // Check if user is logged in
      const token = localStorage.getItem('access_token');
      if (token) {
        setIsGuest(false);
        try {
          const response = await api.get('/auth/me');
          const user = response.data.data;
          setGuestInfo({
            ...guestInfo,
            full_name: user.name || '',
            email: user.email || '',
            phone: user.phone || '',
          });
        } catch (error) {
          console.log('User not logged in');
          setIsGuest(true);
        }
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

    if (paymentMethod === 'mpesa' && !mpesaPhone) {
      toast.error('Please enter M-Pesa phone number');
      return;
    }

    if (isGuest) {
      if (!guestInfo.full_name || !guestInfo.email || !guestInfo.phone || !guestInfo.address) {
        toast.error('Please fill in all guest details');
        return;
      }
    }

    setProcessing(true);

    try {
      const orderData = {
        items: cartItems.map(item => ({
          variant_id: item.id,
          quantity: item.quantity,
          price: item.price
        })),
        subtotal: subtotal,
        shipping_fee: shipping,
        total: total,
        payment_method: paymentMethod,
        is_guest: isGuest,
        guest_info: isGuest ? guestInfo : null,
        phone: mpesaPhone
      };

      const response = await api.post('/orders/guest', orderData);
      const order = response.data.data;
      
      if (paymentMethod === 'mpesa') {
        // Initiate M-Pesa payment
        const paymentResponse = await api.post(
          isGuest ? '/payments/mpesa/guest/stkpush' : '/payments/mpesa/stkpush',
          {
            order_id: order.id,
            phone: mpesaPhone,
            amount: total,
            is_guest: isGuest
        }
      );
        
        if (paymentResponse.data.success) {
          toast.success('Payment initiated! Check your phone for M-Pesa prompt');
          localStorage.removeItem('cart');
          router.push(`/orders/${order.id}?payment=processing`);
        } else {
          toast.error(paymentResponse.data.error || 'Payment initiation failed');
        }
      } else {
        toast.success('Order placed successfully!');
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

  return (
    <>
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 py-16">
        <h1 className="text-3xl font-bold text-black mb-8">Checkout</h1>

        <div className="grid lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-6">
            {/* Guest Checkout Form */}
            {isGuest && (
              <div className="bg-white rounded-xl shadow-md p-6">
                <h2 className="text-xl font-bold text-black mb-4">Guest Details</h2>
                <div className="space-y-4">
                  <div className="grid md:grid-cols-2 gap-4">
                    <input
                      type="text"
                      placeholder="Full Name *"
                      value={guestInfo.full_name}
                      onChange={(e) => setGuestInfo({ ...guestInfo, full_name: e.target.value })}
                      className="w-full px-4 py-2 border rounded-lg text-black"
                    />
                    <input
                      type="email"
                      placeholder="Email *"
                      value={guestInfo.email}
                      onChange={(e) => setGuestInfo({ ...guestInfo, email: e.target.value })}
                      className="w-full px-4 py-2 border rounded-lg text-black"
                    />
                  </div>
                  <input
                    type="text"
                    placeholder="Phone Number *"
                    value={guestInfo.phone}
                    onChange={(e) => setGuestInfo({ ...guestInfo, phone: e.target.value })}
                    className="w-full px-4 py-2 border rounded-lg text-black"
                  />
                  <input
                    type="text"
                    placeholder="Delivery Address *"
                    value={guestInfo.address}
                    onChange={(e) => setGuestInfo({ ...guestInfo, address: e.target.value })}
                    className="w-full px-4 py-2 border rounded-lg text-black"
                  />
                  <div className="grid md:grid-cols-2 gap-4">
                    <input
                      type="text"
                      placeholder="City"
                      value={guestInfo.city}
                      onChange={(e) => setGuestInfo({ ...guestInfo, city: e.target.value })}
                      className="w-full px-4 py-2 border rounded-lg text-black"
                    />
                    <input
                      type="text"
                      placeholder="County"
                      value={guestInfo.county}
                      onChange={(e) => setGuestInfo({ ...guestInfo, county: e.target.value })}
                      className="w-full px-4 py-2 border rounded-lg text-black"
                    />
                  </div>
                  <p className="text-sm text-gray-500">
                    <Link href="/login" className="text-blue-600 hover:underline">
                      Already have an account? Login here →
                    </Link>
                  </p>
                </div>
              </div>
            )}

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

            <button
              onClick={handlePlaceOrder}
              disabled={processing || cartItems.length === 0}
              className={`w-full py-3 rounded-full mt-6 font-medium transition-all ${
                processing 
                  ? 'bg-gray-400 cursor-not-allowed' 
                  : 'bg-black text-white hover:bg-gray-800 cursor-pointer'
              }`}
            >
              {processing ? 'Processing...' : `Place Order • KSh ${total.toLocaleString()}`}
            </button>

            {isGuest && (
              <p className="text-xs text-center text-gray-500 mt-3">
                Continue as guest or <Link href="/login" className="text-blue-600">Login</Link>
              </p>
            )}
          </div>
        </div>
      </div>
    </>
  );
}