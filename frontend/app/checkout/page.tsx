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
  const [paymentStatus, setPaymentStatus] = useState<'idle' | 'pending' | 'success' | 'failed'>('idle');
  const [paymentMessage, setPaymentMessage] = useState('');
  const [orderId, setOrderId] = useState('');
  const [guestInfo, setGuestInfo] = useState({
    full_name: '',
    email: '',
    phone: '',
    address: '',
    city: '',
    county: ''
  });
  const [statusCheckCount, setStatusCheckCount] = useState(0);

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
    setPaymentStatus('pending');
    setPaymentMessage('Processing your order...');

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

      const orderResponse = await api.post('/orders/guest', orderData);
      const order = orderResponse.data.data;
      setOrderId(order.id);
      
      if (paymentMethod === 'mpesa') {
        setPaymentMessage('Initiating M-Pesa payment...');
        
        const paymentResponse = await api.post('/payments/mpesa/guest/stkpush', {
          order_id: order.id,
          phone: mpesaPhone,
          amount: total
        });
        
        if (paymentResponse.data.success) {
          setPaymentStatus('pending');
          setPaymentMessage('Payment initiated! Check your phone for M-Pesa prompt.');
          toast.success('Payment initiated! Check your phone for M-Pesa prompt');
          localStorage.removeItem('cart');
          
          // Start checking payment status
          setStatusCheckCount(0);
          checkPaymentStatus(order.id);
          
        } else {
          setPaymentStatus('failed');
          setPaymentMessage(paymentResponse.data.error || 'Payment initiation failed');
          toast.error(paymentResponse.data.error || 'Payment initiation failed');
        }
      } else {
        setPaymentStatus('success');
        setPaymentMessage('Order placed successfully! You will pay on delivery.');
        toast.success('Order placed successfully!');
        localStorage.removeItem('cart');
      }
      
    } catch (error: any) {
      console.error('Order failed:', error);
      setPaymentStatus('failed');
      const errorMsg = error.response?.data?.error || 'Failed to place order';
      setPaymentMessage(errorMsg);
      toast.error(errorMsg);
    } finally {
      setProcessing(false);
    }
  };

  const checkPaymentStatus = async (orderId: string) => {
    // Only check up to 10 times (50 seconds)
    if (statusCheckCount > 24) {
      setPaymentStatus('failed');
      setPaymentMessage('Payment status check timed out. Please check your orders page.');
      return;
    }

    setStatusCheckCount(prev => prev + 1);
    
    try {
      // Use a different endpoint that doesn't require auth for guests
      const response = await api.get(`/payments/guest/${orderId}/status`);
      const payment = response.data.data;
      console.log('💳 Payment status response:', payment);
      
      if (payment.status === 'success') {
        setPaymentStatus('success');
        setPaymentMessage('Payment successful! Your order has been confirmed.');
        toast.success('Payment successful!');
        return;
      } else if (payment.status === 'failed') {
        setPaymentStatus('failed');
        setPaymentMessage(`Payment failed: ${payment.failure_reason || 'Unknown error'}`);
        toast.error('Payment failed');
        return;
      } else {
        // Still pending, check again after 5 seconds
        setTimeout(() => {
          checkPaymentStatus(orderId);
        }, 5000);
      }
    } catch (error) {
      console.error('Failed to check payment status:', error);
      // If it's a 401, try again (might be token issue)
      setTimeout(() => {
        checkPaymentStatus(orderId);
      }, 5000);
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

  if (cartItems.length === 0 && paymentStatus === 'idle') {
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

        {/* Payment Status Display */}
        {paymentStatus !== 'idle' && (
          <div className={`mb-6 p-6 rounded-xl shadow-md ${
            paymentStatus === 'success' ? 'bg-green-50 border border-green-200' :
            paymentStatus === 'failed' ? 'bg-red-50 border border-red-200' :
            'bg-yellow-50 border border-yellow-200'
          }`}>
            <div className="flex items-center gap-3">
              <span className="text-3xl">
                {paymentStatus === 'success' ? '✅' :
                 paymentStatus === 'failed' ? '❌' :
                 '⏳'}
              </span>
              <div>
                <h3 className={`font-bold ${
                  paymentStatus === 'success' ? 'text-green-700' :
                  paymentStatus === 'failed' ? 'text-red-700' :
                  'text-yellow-700'
                }`}>
                  {paymentStatus === 'success' ? 'Payment Successful!' :
                   paymentStatus === 'failed' ? 'Payment Failed' :
                   'Processing Payment...'}
                </h3>
                <p className="text-black">{paymentMessage}</p>
                {orderId && (
                  <p className="text-sm text-gray-600 mt-1">Order ID: {orderId}</p>
                )}
              </div>
            </div>
            {paymentStatus !== 'pending' && (
              <div className="mt-4 flex gap-4 flex-wrap">
                <Link href="/orders" className="bg-black text-white px-6 py-2 rounded-lg hover:bg-gray-800">
                  View Orders
                </Link>
                <Link href="/" className="border border-gray-300 text-black px-6 py-2 rounded-lg hover:bg-gray-100">
                  Continue Shopping
                </Link>
                {paymentStatus === 'failed' && (
                  <button 
                    onClick={() => {
                      setPaymentStatus('idle');
                      setPaymentMessage('');
                      window.location.reload();
                    }}
                    className="border border-blue-600 text-blue-600 px-6 py-2 rounded-lg hover:bg-blue-50"
                  >
                    Try Again
                  </button>
                )}
              </div>
            )}
          </div>
        )}

        {paymentStatus === 'idle' && (
          <div className="grid lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 space-y-6">
              {/* Guest Checkout Form */}
              {isGuest && (
                <div className="bg-white rounded-xl shadow-md p-6">
                  <h2 className="text-xl font-bold text-black mb-4">Guest Details</h2>
                  <div className="space-y-4">
                    <div className="grid md:grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-black mb-1">Full Name *</label>
                        <input
                          type="text"
                          value={guestInfo.full_name}
                          onChange={(e) => setGuestInfo({ ...guestInfo, full_name: e.target.value })}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg text-black"
                          placeholder="John Doe"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-black mb-1">Email *</label>
                        <input
                          type="email"
                          value={guestInfo.email}
                          onChange={(e) => setGuestInfo({ ...guestInfo, email: e.target.value })}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg text-black"
                          placeholder="john@example.com"
                        />
                      </div>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-black mb-1">Phone Number *</label>
                      <input
                        type="text"
                        value={guestInfo.phone}
                        onChange={(e) => setGuestInfo({ ...guestInfo, phone: e.target.value })}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg text-black"
                        placeholder="0712345678"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-black mb-1">Delivery Address *</label>
                      <input
                        type="text"
                        value={guestInfo.address}
                        onChange={(e) => setGuestInfo({ ...guestInfo, address: e.target.value })}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg text-black"
                        placeholder="123 Kenyatta Ave"
                      />
                    </div>
                    <div className="grid md:grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-black mb-1">City</label>
                        <input
                          type="text"
                          value={guestInfo.city}
                          onChange={(e) => setGuestInfo({ ...guestInfo, city: e.target.value })}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg text-black"
                          placeholder="Nairobi"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-black mb-1">County</label>
                        <input
                          type="text"
                          value={guestInfo.county}
                          onChange={(e) => setGuestInfo({ ...guestInfo, county: e.target.value })}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg text-black"
                          placeholder="Nairobi"
                        />
                      </div>
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
            <div className="bg-white rounded-xl shadow-md p-6 h-fit sticky top-24">
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
        )}
      </div>
    </>
  );
}
