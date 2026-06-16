'use client';

import Link from 'next/link';
import Navbar from '../components/Navbar';
import { useEffect, useState } from 'react';
import api from '../services/api';

export default function Home() {
  const [products, setProducts] = useState([]);

  useEffect(() => {
    fetchProducts();
  }, []);

  const fetchProducts = async () => {
    try {
      const response = await api.get('/products');
      setProducts(response.data.data.data || []);
    } catch (error) {
      console.error('Failed to fetch products:', error);
    }
  };

  return (
    <>
      <Navbar />
      
      <section className="bg-black text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 text-center">
          <h1 className="text-5xl md:text-7xl font-bold mb-4">
            NIFNOC iMobile
          </h1>
          <p className="text-xl md:text-2xl mb-2 text-gray-300">
            ONLY NEW iPHONES
          </p>
          <p className="text-lg mb-8 text-gray-400">
            BUY NEW iPHONES FOR SALE
          </p>
          <Link href="/products" className="bg-blue-600 text-white px-8 py-3 rounded-full font-medium hover:bg-blue-700 transition-all inline-block">
            Shop Now →
          </Link>
        </div>
      </section>

      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <h2 className="text-3xl font-bold text-center mb-12">NEW iPHONES FOR SALE</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {products.length === 0 ? (
            <p className="text-center col-span-4 text-gray-500">Loading products...</p>
          ) : (
            products.slice(0, 4).map((product: any) => (
              <div key={product.id} className="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition-shadow">
                <div className="p-4">
                  {product.variants?.[0]?.images?.[0] && (
                    <img 
                      src={product.variants[0].images[0]} 
                      alt={product.name}
                      className="w-full h-48 object-contain mb-4"
                    />
                  )}
                  <h3 className="font-semibold text-lg mb-2">{product.name}</h3>
                  <p className="text-blue-600 font-bold text-xl">
                    KSh {product.base_price?.toLocaleString()}
                  </p>
                  <Link 
                    href={`/products/${product.slug}`}
                    className="bg-black text-white block text-center py-3 rounded-full mt-4 text-sm hover:bg-gray-800 transition-all"
                  >
                    Buy Now
                  </Link>
                </div>
              </div>
            ))
          )}
        </div>
      </section>
    </>
  );
}
