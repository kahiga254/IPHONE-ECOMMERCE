'use client';
export const dynamic = 'force-dynamic';


import Link from 'next/link';
import Navbar from '@/components/Navbar';
import { useEffect, useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import api from '@/services/api';



export default function ProductsPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [categories, setCategories] = useState([]);
  const [filters, setFilters] = useState({
    search: searchParams.get('search') || '',
    category: searchParams.get('category') || '',
    minPrice: searchParams.get('minPrice') || '',
    maxPrice: searchParams.get('maxPrice') || '',
    sortBy: searchParams.get('sortBy') || 'created_at',
    order: searchParams.get('order') || 'desc',
  });
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 12,
    total: 0,
    totalPages: 0,
  });

  useEffect(() => {
    fetchCategories();
    fetchProducts();
  }, [filters, pagination.page]);

  const fetchCategories = async () => {
    try {
      const response = await api.get('/categories');
      setCategories(response.data.data?.categories || []);
    } catch (error) {
      console.error('Failed to fetch categories:', error);
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

  const fetchProducts = async () => {
    setLoading(true);
    try {
      const params = new URLSearchParams();
      params.append('page', pagination.page.toString());
      params.append('limit', pagination.limit.toString());
      if (filters.search) params.append('search', filters.search);
      if (filters.category) params.append('category', filters.category);
      if (filters.minPrice) params.append('min_price', filters.minPrice);
      if (filters.maxPrice) params.append('max_price', filters.maxPrice);
      if (filters.sortBy) params.append('sort_by', filters.sortBy);
      if (filters.order) params.append('order', filters.order);
      
      const response = await api.get(`/products?${params.toString()}`);
      console.log('Products data:', response.data.data.data);
      setProducts(response.data.data.data || []);
      setPagination({
        ...pagination,
        total: response.data.data.total || 0,
        totalPages: response.data.data.total_pages || 0,
      });
    } catch (error) {
      console.error('Failed to fetch products:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (key: string, value: string) => {
    setFilters({ ...filters, [key]: value });
    setPagination({ ...pagination, page: 1 });
    
    const params = new URLSearchParams();
    Object.keys({ ...filters, [key]: value }).forEach((k) => {
      if (k !== 'page' && (filters as any)[k]) {
        params.append(k, (filters as any)[k]);
      }
    });
    router.push(`/products?${params.toString()}`);
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const [sortBy, order] = e.target.value.split('-');
    handleFilterChange('sortBy', sortBy);
    handleFilterChange('order', order);
  };

  const clearFilters = () => {
    setFilters({
      search: '',
      category: '',
      minPrice: '',
      maxPrice: '',
      sortBy: 'created_at',
      order: 'desc',
    });
    setPagination({ ...pagination, page: 1 });
    router.push('/products');
  };

  const goToPage = (page: number) => {
    setPagination({ ...pagination, page });
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

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
        <h1 className="text-3xl font-bold text-black mb-8">All Products</h1>

        {/* Search Bar */}
        <div className="mb-6">
          <input
            type="text"
            placeholder="Search products..."
            value={filters.search}
            onChange={(e) => handleFilterChange('search', e.target.value)}
            className="w-full max-w-md px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-black focus:border-transparent text-black"
          />
        </div>

        <div className="flex flex-col lg:flex-row gap-8">
          {/* Sidebar Filters */}
          <div className="lg:w-64 space-y-6">
            <div className="bg-white rounded-xl shadow-md p-4">
              <h3 className="font-semibold text-black mb-3">Categories</h3>
              <select
                value={filters.category}
                onChange={(e) => handleFilterChange('category', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-black"
              >
                <option value="">All Categories</option>
                {categories.map((cat: any) => (
                  <option key={cat.id} value={cat.slug}>{cat.name}</option>
                ))}
              </select>
            </div>

            <div className="bg-white rounded-xl shadow-md p-4">
              <h3 className="font-semibold text-black mb-3">Price Range</h3>
              <div className="space-y-3">
                <input
                  type="number"
                  placeholder="Min Price"
                  value={filters.minPrice}
                  onChange={(e) => handleFilterChange('minPrice', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-black"
                />
                <input
                  type="number"
                  placeholder="Max Price"
                  value={filters.maxPrice}
                  onChange={(e) => handleFilterChange('maxPrice', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-black"
                />
              </div>
            </div>

            <div className="bg-white rounded-xl shadow-md p-4">
              <h3 className="font-semibold text-black mb-3">Sort By</h3>
              <select
                value={`${filters.sortBy}-${filters.order}`}
                onChange={handleSortChange}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-black"
              >
                <option value="created_at-desc">Latest</option>
                <option value="base_price-asc">Price: Low to High</option>
                <option value="base_price-desc">Price: High to Low</option>
                <option value="name-asc">Name: A to Z</option>
                <option value="name-desc">Name: Z to A</option>
              </select>
            </div>

            <button
              onClick={clearFilters}
              className="w-full bg-gray-200 text-black px-4 py-2 rounded-lg hover:bg-gray-300 transition"
            >
              Clear Filters
            </button>
          </div>

          {/* Products Grid */}
          <div className="flex-1">
            {products.length === 0 ? (
              <div className="text-center py-16 bg-white rounded-xl">
                <p className="text-gray-500">No products found</p>
                <button
                  onClick={clearFilters}
                  className="mt-4 text-blue-600 hover:underline"
                >
                  Clear filters
                </button>
              </div>
            ) : (
              <>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {products.map((product: any) => {
                    const image = getProductImage(product);
                    return (
                      <div key={product.id} className="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition">
                        <div className="p-4">
                          {/* Product Image */}
                          <div className="bg-gray-100 rounded-lg h-48 flex items-center justify-center mb-3">
                            {image ? (
                              <img
                                src={image}
                                alt={product.name}
                                className="w-full h-48 object-contain"
                                onError={(e) => {
                                  console.error('Image load error:', image);
                                  (e.target as HTMLImageElement).style.display = 'none';
                                }}
                              />
                            ) : (
                              <div className="text-6xl">📱</div>
                            )}
                          </div>
                          <h3 className="font-semibold text-lg mb-2 text-black">{product.name}</h3>
                          <p className="text-blue-600 font-bold text-xl">
                            KSh {product.base_price?.toLocaleString()}
                          </p>
                          <Link 
                            href={`/products/${product.slug}`}
                            className="bg-black text-white block text-center py-2 rounded-full mt-4 text-sm hover:bg-gray-800"
                          >
                            View Details
                          </Link>
                        </div>
                      </div>
                    );
                  })}
                </div>

                {/* Pagination */}
                {pagination.totalPages > 1 && (
                  <div className="flex justify-center gap-2 mt-8">
                    <button
                      onClick={() => goToPage(pagination.page - 1)}
                      disabled={pagination.page === 1}
                      className="px-4 py-2 border rounded-lg disabled:opacity-50 hover:bg-gray-100"
                    >
                      Previous
                    </button>
                    <span className="px-4 py-2 text-black">
                      Page {pagination.page} of {pagination.totalPages}
                    </span>
                    <button
                      onClick={() => goToPage(pagination.page + 1)}
                      disabled={pagination.page === pagination.totalPages}
                      className="px-4 py-2 border rounded-lg disabled:opacity-50 hover:bg-gray-100"
                    >
                      Next
                    </button>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
