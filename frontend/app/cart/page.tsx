export const dynamic = 'force-dynamic';
'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import Navbar from '@/components/Navbar';
import toast from 'react-hot-toast';

interface CartItem {
  id: string;
  name: string;
  price: number;
  quantity: number;
  image: string;
  variant: string;
}

export default function CartPage() {
  const [cartItems, setCartItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadCart();
  }, []);

  const loadCart = () => {
    const savedCart = localStorage.getItem('cart');
    if (savedCart) {
      const items = JSON.parse(savedCart);
      const uniqueItems = items.reduce((acc: CartItem[], current: CartItem) => {
        const exists = acc.find(item => item.id === current.id);
        if (!exists) {
          acc.push(current);
        } else {
          exists.quantity += current.quantity;
        }
        return acc;
      }, []);
      setCartItems(uniqueItems);
    }
    setLoading(false);
  };

  const saveCart = (items: CartItem[]) => {
    localStorage.setItem('cart', JSON.stringify(items));
    window.dispatchEvent(new Event('cartUpdated'));
    setCartItems(items);
  };

  const updateQuantity = (id: string, newQuantity: number) => {
    if (newQuantity < 1) return;
    const updated = cartItems.map(item =>
      item.id === id ? { ...item, quantity: newQuantity } : item
    );
    saveCart(updated);
    toast.success('Cart updated');
  };

  const removeItem = (id: string) => {
    const updated = cartItems.filter(item => item.id !== id);
    saveCart(updated);
    toast.success('Item removed');
  };

  const subtotal = cartItems.reduce((sum, item) => sum + item.price * item.quantity, 0);
  const shipping = subtotal > 0 ? 200 : 0;
  const total = subtotal + shipping;

  if (loading) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center text-black">Loading...</div>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 py-16">
        <h1 className="text-3xl font-bold text-black mb-8">Shopping Cart</h1>

        {cartItems.length === 0 ? (
          <div className="text-center py-16">
            <p className="text-black mb-4">Your cart is empty</p>
            <Link href="/products" className="bg-black text-white px-6 py-2 rounded-full hover:bg-gray-800">
              Continue Shopping
            </Link>
          </div>
        ) : (
          <div className="grid lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 space-y-4">
              {cartItems.map((item) => (
                <div key={item.id} className="bg-white rounded-xl shadow-md p-4 flex gap-4">
                  <div className="w-24 h-24 bg-gray-100 rounded-lg flex items-center justify-center">
                    {item.image ? (
                      <img src={item.image} alt={item.name} className="w-20 h-20 object-contain" />
                    ) : (
                      <div className="w-20 h-20 bg-gray-200 rounded flex items-center justify-center text-gray-500 text-xs">
                        No image
                      </div>
                    )}
                  </div>
                  <div className="flex-1">
                    <h3 className="font-semibold text-lg text-black">{item.name}</h3>
                    <p className="text-gray-600 text-sm">{item.variant}</p>
                    <p className="text-blue-600 font-bold text-lg">KSh {item.price.toLocaleString()}</p>
                  </div>
                  <div className="flex flex-col items-end gap-2">
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() => updateQuantity(item.id, item.quantity - 1)}
                        className="w-8 h-8 border border-gray-300 rounded-full hover:bg-gray-100 text-black"
                      >
                        -
                      </button>
                      <span className="w-8 text-center text-black font-medium">{item.quantity}</span>
                      <button
                        onClick={() => updateQuantity(item.id, item.quantity + 1)}
                        className="w-8 h-8 border border-gray-300 rounded-full hover:bg-gray-100 text-black"
                      >
                        +
                      </button>
                    </div>
                    <button
                      onClick={() => removeItem(item.id)}
                      className="text-red-600 text-sm hover:underline"
                    >
                      Remove
                    </button>
                  </div>
                </div>
              ))}
            </div>

            <div className="bg-white rounded-xl shadow-md p-6 h-fit">
              <h2 className="text-xl font-bold text-black mb-4">Order Summary</h2>
              <div className="space-y-2 border-b border-gray-200 pb-4">
                <div className="flex justify-between text-black">
                  <span>Subtotal</span>
                  <span className="font-semibold">KSh {subtotal.toLocaleString()}</span>
                </div>
                <div className="flex justify-between text-black">
                  <span>Shipping</span>
                  <span className="font-semibold">KSh {shipping.toLocaleString()}</span>
                </div>
              </div>
              <div className="flex justify-between font-bold text-xl mt-4 text-black">
                <span>Total</span>
                <span className="text-blue-600">KSh {total.toLocaleString()}</span>
              </div>
              <Link
                href="/checkout"
                className="block bg-black text-white text-center py-3 rounded-full mt-6 hover:bg-gray-800 transition"
              >
                Proceed to Checkout
              </Link>
            </div>
          </div>
        )}
      </div>
    </>
  );
}
