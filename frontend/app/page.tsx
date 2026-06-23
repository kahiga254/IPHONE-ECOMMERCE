'use client';
export const dynamic = 'force-dynamic';


import Link from 'next/link';
import Navbar from '@/components/Navbar';
import { useEffect, useState } from 'react';
import api from '@/services/api';

export default function Home() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [productsRes, categoriesRes] = await Promise.all([
        api.get('/products'),
        api.get('/categories')
      ]);
      setProducts(productsRes.data.data.data || []);
      setCategories(categoriesRes.data.data?.categories || []);
    } catch (error) {
      console.error('Failed to fetch data:', error);
    } finally {
      setLoading(false);
    }
  };

  const getProductImage = (product: any) => {
    if (product.variants && product.variants.length > 0) {
      const images = product.variants[0].images;
      if (images && images.length > 0) {
        return images[0];
      }
    }
    return null;
  };

  // Define the categories we want to display in order
  const displayCategories = [
    { name: 'iPhones', icon: '📱', slug: 'iphones' },
    { name: 'Apple Accessories', icon: '🎧', slug: 'apple-accessories' },
    { name: 'Samsung Phones', icon: '📲', slug: 'samsung-phones' },
    { name: 'Samsung Accessories', icon: '⌚', slug: 'samsung-accessories' },
  ];

  const featuredProducts = products.slice(0, 8);

  return (
    <>
      <Navbar />
      
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-black to-gray-900 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 md:py-24">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <h1 className="text-4xl md:text-6xl font-bold leading-tight mb-4">
                <span className="text-blue-500">NIFNOC</span> iMobile
              </h1>
              <p className="text-xl md:text-2xl mb-2 text-gray-300">
                ONLY NEW iPHONES
              </p>
              <p className="text-lg mb-6 text-gray-400">
                Buy the latest iPhones at the best prices in Kenya
              </p>
              <div className="flex flex-wrap gap-4">
                <Link href="/products" className="bg-blue-600 text-white px-8 py-3 rounded-full font-medium hover:bg-blue-700 transition-all inline-block">
                  Shop Now →
                </Link>
                <Link href="/products" className="border border-white text-white px-8 py-3 rounded-full font-medium hover:bg-white hover:text-black transition-all inline-block">
                  View All Products
                </Link>
              </div>
            </div>
            <div className="hidden md:block">
              <div className="bg-gradient-to-br from-blue-600 to-purple-600 rounded-3xl p-8 text-center">
                <div className="text-7xl mb-4">📱</div>
                <p className="text-2xl font-bold">Latest iPhones</p>
                <p className="text-gray-200">Up to 12 Months Warranty</p>
                <div className="mt-4 inline-block bg-white/20 px-6 py-2 rounded-full">
                  🔥 New Arrivals
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Shop by Category Section */}
      <section className="py-16 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-black text-center mb-4">Shop by Category</h2>
          <p className="text-center text-gray-500 mb-10">Find exactly what you're looking for</p>
          
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            {displayCategories.map((cat) => (
              <Link
                key={cat.slug}
                href={`/products?category=${cat.slug}`}
                className="bg-white rounded-2xl shadow-md p-8 text-center hover:shadow-xl transition-all duration-300 hover:-translate-y-1 border border-gray-100 group"
              >
                <div className="text-5xl mb-4 group-hover:scale-110 transition-transform duration-300">
                  {cat.icon}
                </div>
                <h3 className="font-semibold text-black text-lg">{cat.name}</h3>
                <p className="text-sm text-blue-600 mt-2 group-hover:underline">
                  Shop now →
                </p>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* Featured Products */}
      <section className="py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-black">Featured Products</h2>
            <Link href="/products" className="text-blue-600 hover:underline text-sm">
              View All →
            </Link>
          </div>
          
          {loading ? (
            <div className="text-center py-12 text-black">Loading products...</div>
          ) : featuredProducts.length === 0 ? (
            <div className="text-center py-12 text-gray-500">No products available</div>
          ) : (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {featuredProducts.map((product: any) => {
                const image = getProductImage(product);
                return (
                  <div key={product.id} className="bg-white rounded-xl shadow-sm overflow-hidden hover:shadow-lg transition border border-gray-100">
                    <div className="p-4">
                      <div className="bg-gray-100 rounded-lg h-40 flex items-center justify-center mb-3">
                        {image ? (
                          <img
                            src={image}
                            alt={product.name}
                            className="w-full h-40 object-contain"
                            onError={(e) => {
                              (e.target as HTMLImageElement).src = '';
                              (e.target as HTMLImageElement).alt = 'No image';
                            }}
                          />
                        ) : (
                          <div className="text-6xl">📱</div>
                        )}
                      </div>
                      <h3 className="font-semibold text-black text-sm truncate">{product.name}</h3>
                      <p className="text-blue-600 font-bold text-base mt-1">
                        KSh {product.base_price?.toLocaleString()}
                      </p>
                      <Link
                        href={`/products/${product.slug}`}
                        className="block w-full bg-black text-white text-center py-2 rounded-full mt-3 text-sm hover:bg-gray-800 transition"
                      >
                        Buy Now
                      </Link>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </section>

      {/* Features Section */}
      <section className="py-12 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div className="text-center">
              <div className="text-3xl mb-2">🚚</div>
              <p className="font-semibold text-black">Free Delivery</p>
              <p className="text-xs text-gray-500">On orders over KSh 50,000</p>
            </div>
            <div className="text-center">
              <div className="text-3xl mb-2">🛡️</div>
              <p className="font-semibold text-black">Authentic Products</p>
              <p className="text-xs text-gray-500">100% Genuine iPhones</p>
            </div>
            <div className="text-center">
              <div className="text-3xl mb-2">🔄</div>
              <p className="font-semibold text-black">7-Day Returns</p>
              <p className="text-xs text-gray-500">Money back guarantee</p>
            </div>
            <div className="text-center">
              <div className="text-3xl mb-2">📞</div>
              <p className="font-semibold text-black">24/7 Support</p>
              <p className="text-xs text-gray-500">We're here to help</p>
            </div>
          </div>
        </div>
      </section>

      {/* Newsletter */}
      <section className="py-12 bg-black text-white">
        <div className="max-w-2xl mx-auto px-4 text-center">
          <h2 className="text-2xl font-bold mb-2">Stay Updated</h2>
          <p className="text-gray-400 mb-4">Get notified about new arrivals and exclusive offers</p>
          <div className="flex gap-3 max-w-md mx-auto">
            <input
              type="email"
              placeholder="Enter your email"
              className="flex-1 px-4 py-3 rounded-full text-black outline-none"
            />
            <button className="bg-blue-600 px-6 py-3 rounded-full hover:bg-blue-700 transition">
              Subscribe
            </button>
          </div>
        </div>
      </section>
    </>
  );
}