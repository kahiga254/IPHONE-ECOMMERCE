'use client';

import { useEffect, useState } from 'react';
import api from '@/services/api';
import toast from 'react-hot-toast';

interface Variant {
  id: string;
  sku: string;
  color: string;
  storage: string;
  price: number;
  stock: number;
  images: string[];
}

interface Product {
  id: string;
  name: string;
  slug: string;
  base_price: number;
  is_active: boolean;
  variants: Variant[];
}

export default function AdminProducts() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingStock, setEditingStock] = useState<{ productId: string; variantId: string; stock: number } | null>(null);
  const [showAddForm, setShowAddForm] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [newProduct, setNewProduct] = useState({
    name: '',
    slug: '',
    description: '',
    base_price: '',
    variants: [{ sku: '', color: '', storage: '', price: '', stock: '', images: [] as string[] }]
  });

  useEffect(() => {
    fetchProducts();
  }, []);

  const fetchProducts = async () => {
    try {
      const response = await api.get('/products');
      setProducts(response.data.data?.data || []);
    } catch (error) {
      console.error('Failed to fetch products:', error);
      toast.error('Failed to load products');
    } finally {
      setLoading(false);
    }
  };

  const updateStock = async (variantId: string, newStock: number) => {
    try {
      await api.patch(`/admin/variants/${variantId}/stock`, { stock: newStock });
      toast.success('Stock updated successfully');
      fetchProducts();
      setEditingStock(null);
    } catch (error) {
      toast.error('Failed to update stock');
    }
  };

  const getStockStatus = (stock: number) => {
    if (stock === 0) return { label: 'Out of Stock', color: 'bg-red-100 text-red-800', textColor: 'text-red-700' };
    if (stock < 10) return { label: 'Low Stock', color: 'bg-yellow-100 text-yellow-800', textColor: 'text-yellow-700' };
    return { label: 'In Stock', color: 'bg-green-100 text-green-800', textColor: 'text-green-700' };
  };

  const toggleProductStatus = async (id: string, currentStatus: boolean) => {
    try {
      await api.patch(`/admin/products/${id}`, { is_active: !currentStatus });
      toast.success(`Product ${!currentStatus ? 'activated' : 'deactivated'}`);
      fetchProducts();
    } catch (error) {
      toast.error('Failed to update product status');
    }
  };

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>, variantIndex: number) => {
    const files = Array.from(e.target.files || []);
    if (files.length === 0) return;

    setUploading(true);

    for (const file of files) {
      const formData = new FormData();
      formData.append('file', file);

      try {
        const response = await fetch('/api/upload', {
          method: 'POST',
          body: formData,
        });
        const data = await response.json();
        
        if (data.success) {
          const variants = [...newProduct.variants];
          variants[variantIndex].images.push(data.url);
          setNewProduct({ ...newProduct, variants });
          toast.success('Image uploaded');
        } else {
          toast.error('Upload failed');
        }
      } catch (error) {
        console.error('Upload failed:', error);
        toast.error('Failed to upload image');
      }
    }

    setUploading(false);
  };

  const removeImage = (variantIndex: number, imageIndex: number) => {
    const variants = [...newProduct.variants];
    variants[variantIndex].images.splice(imageIndex, 1);
    setNewProduct({ ...newProduct, variants });
  };

  const addVariant = () => {
    setNewProduct({
      ...newProduct,
      variants: [...newProduct.variants, { sku: '', color: '', storage: '', price: '', stock: '', images: [] }]
    });
  };

  const removeVariant = (index: number) => {
    const variants = [...newProduct.variants];
    variants.splice(index, 1);
    setNewProduct({ ...newProduct, variants });
  };

  const updateVariant = (index: number, field: string, value: string) => {
    const variants = [...newProduct.variants];
    variants[index] = { ...variants[index], [field]: value };
    setNewProduct({ ...newProduct, variants });
  };

  const createProduct = async () => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      toast.error('Please login');
      return;
    }

    if (!newProduct.name || !newProduct.base_price) {
      toast.error('Please fill product name and price');
      return;
    }

    const slug = newProduct.slug || newProduct.name.toLowerCase().replace(/ /g, '-');
    
    const productData = {
      name: newProduct.name,
      slug: slug,
      description: newProduct.description,
      base_price: parseFloat(newProduct.base_price),
      variants: newProduct.variants.filter(v => v.sku && v.color && v.storage && v.price).map(v => ({
        sku: v.sku,
        color: v.color,
        storage: v.storage,
        price: parseFloat(v.price),
        stock: parseInt(v.stock) || 0,
        images: v.images.length > 0 ? v.images : ['/product-placeholder.jpg']
      }))
    };

    try {
      await api.post('/admin/products', productData);
      toast.success('Product created successfully');
      setShowAddForm(false);
      setNewProduct({
        name: '',
        slug: '',
        description: '',
        base_price: '',
        variants: [{ sku: '', color: '', storage: '', price: '', stock: '', images: [] }]
      });
      fetchProducts();
    } catch (error: any) {
      console.error('Failed to create product:', error);
      toast.error(error.response?.data?.error || 'Failed to create product');
    }
  };

  if (loading) {
    return <div className="text-center py-16 text-black">Loading products...</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-black">Products & Stock Management</h1>
        <button
          onClick={() => setShowAddForm(true)}
          className="bg-black text-white px-4 py-2 rounded-lg hover:bg-gray-800"
        >
          + Add New Product
        </button>
      </div>

      {showAddForm && (
        <div className="bg-white rounded-xl shadow-md p-6 mb-6">
          <h2 className="text-xl font-bold text-black mb-4">Add New Product</h2>
          <div className="space-y-4">
            <div className="grid md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-black mb-1">Product Name *</label>
                <input
                  type="text"
                  value={newProduct.name}
                  onChange={(e) => setNewProduct({ ...newProduct, name: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg text-black"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-black mb-1">Base Price *</label>
                <input
                  type="number"
                  value={newProduct.base_price}
                  onChange={(e) => setNewProduct({ ...newProduct, base_price: e.target.value })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg text-black"
                />
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-black mb-1">Description</label>
              <textarea
                value={newProduct.description}
                onChange={(e) => setNewProduct({ ...newProduct, description: e.target.value })}
                rows={3}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg text-black"
              />
            </div>
            
            <h3 className="font-semibold text-black mt-4">Variants & Images</h3>
            {newProduct.variants.map((variant, idx) => (
              <div key={idx} className="p-4 bg-gray-50 rounded-lg mb-3 border border-gray-200">
                <div className="grid md:grid-cols-5 gap-3 mb-3">
                  <div>
                    <label className="block text-xs font-medium text-black mb-1">SKU</label>
                    <input type="text" placeholder="e.g., IP15-BLK" value={variant.sku} onChange={(e) => updateVariant(idx, 'sku', e.target.value)} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-black text-sm" />
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-black mb-1">Color</label>
                    <input type="text" placeholder="e.g., Black" value={variant.color} onChange={(e) => updateVariant(idx, 'color', e.target.value)} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-black text-sm" />
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-black mb-1">Storage</label>
                    <input type="text" placeholder="e.g., 128GB" value={variant.storage} onChange={(e) => updateVariant(idx, 'storage', e.target.value)} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-black text-sm" />
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-black mb-1">Price</label>
                    <input type="number" placeholder="0" value={variant.price} onChange={(e) => updateVariant(idx, 'price', e.target.value)} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-black text-sm" />
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-black mb-1">Stock</label>
                    <input type="number" placeholder="0" value={variant.stock} onChange={(e) => updateVariant(idx, 'stock', e.target.value)} className="w-full px-3 py-2 border border-gray-300 rounded-lg text-black text-sm" />
                  </div>
                </div>
                
                <div className="mb-2">
                  <label className="block text-sm font-medium text-black mb-2">Product Images</label>
                  <input 
                    type="file" 
                    multiple 
                    accept="image/*" 
                    onChange={(e) => handleImageUpload(e, idx)} 
                    disabled={uploading}
                    className="block w-full text-sm text-black border border-gray-300 rounded-lg p-2 file:mr-3 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-semibold file:bg-black file:text-white hover:file:bg-gray-800"
                  />
                  {uploading && <p className="text-xs text-gray-500 mt-1">Uploading...</p>}
                  {variant.images.length > 0 && (
                    <div className="flex gap-2 mt-3 flex-wrap">
                      {variant.images.map((img, imgIdx) => (
                        <div key={imgIdx} className="relative">
                          <img src={img} alt="Preview" className="w-16 h-16 object-cover rounded-lg border border-gray-300" />
                          <button
                            onClick={() => removeImage(idx, imgIdx)}
                            className="absolute -top-2 -right-2 bg-red-500 text-white rounded-full w-5 h-5 text-xs hover:bg-red-600"
                          >
                            ×
                          </button>
                        </div>
                      ))}
                    </div>
                  )}
                  <p className="text-xs text-gray-500 mt-1">Upload JPG or PNG images</p>
                </div>
                
                {idx > 0 && (
                  <button onClick={() => removeVariant(idx)} className="text-red-600 text-sm mt-2 hover:underline">
                    Remove Variant
                  </button>
                )}
              </div>
            ))}
            <button onClick={addVariant} className="text-blue-600 text-sm hover:underline">
              + Add Another Variant
            </button>
            
            <div className="flex gap-3 pt-4">
              <button onClick={createProduct} className="bg-black text-white px-6 py-2 rounded-lg hover:bg-gray-800">
                Create Product
              </button>
              <button onClick={() => setShowAddForm(false)} className="border border-gray-300 text-black px-6 py-2 rounded-lg hover:bg-gray-100">
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {products.map((product) => {
        const totalStock = product.variants?.reduce((sum, v) => sum + v.stock, 0) || 0;
        const stockStatus = getStockStatus(totalStock);
        return (
          <div key={product.id} className="bg-white rounded-xl shadow-md mb-6 overflow-hidden">
            <div className="p-6 border-b border-gray-200 bg-gray-50">
              <div className="flex justify-between items-center">
                <div>
                  <h2 className="text-xl font-bold text-black">{product.name}</h2>
                  <p className="text-gray-600 text-sm mt-1">Slug: {product.slug}</p>
                </div>
                <div className="text-right">
                  <p className="text-sm text-gray-600">Base Price</p>
                  <p className="text-lg font-bold text-black">KSh {product.base_price?.toLocaleString()}</p>
                </div>
              </div>
              <div className="mt-3 flex justify-between items-center">
                <span className={`inline-block px-2 py-1 rounded text-xs font-medium ${product.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'}`}>
                  {product.is_active ? 'Active' : 'Inactive'}
                </span>
                <button onClick={() => toggleProductStatus(product.id, product.is_active)} className="text-blue-600 hover:underline text-sm">
                  {product.is_active ? 'Deactivate' : 'Activate'}
                </button>
              </div>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-100">
                  <tr className="text-left">
                    <th className="px-6 py-3 text-black font-medium">SKU</th>
                    <th className="px-6 py-3 text-black font-medium">Color</th>
                    <th className="px-6 py-3 text-black font-medium">Storage</th>
                    <th className="px-6 py-3 text-black font-medium">Price</th>
                    <th className="px-6 py-3 text-black font-medium">Stock</th>
                    <th className="px-6 py-3 text-black font-medium">Status</th>
                    <th className="px-6 py-3 text-black font-medium">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {product.variants?.map((variant) => {
                    const vs = getStockStatus(variant.stock);
                    return (
                      <tr key={variant.id} className="border-b border-gray-100">
                        <td className="px-6 py-4 text-black text-sm">{variant.sku}</td>
                        <td className="px-6 py-4 text-black">{variant.color}</td>
                        <td className="px-6 py-4 text-black">{variant.storage}</td>
                        <td className="px-6 py-4 text-black">KSh {variant.price?.toLocaleString()}</td>
                        <td className="px-6 py-4">
                          {editingStock?.variantId === variant.id ? (
                            <input type="number" value={editingStock.stock} onChange={(e) => setEditingStock({ ...editingStock, stock: parseInt(e.target.value) || 0 })} className="w-20 px-2 py-1 border rounded text-black" />
                          ) : (
                            <span className={`font-semibold ${vs.textColor}`}>{variant.stock}</span>
                          )}
                        </td>
                        <td className="px-6 py-4"><span className={`px-2 py-1 rounded text-xs font-medium ${vs.color}`}>{vs.label}</span></td>
                        <td className="px-6 py-4">
                          {editingStock?.variantId === variant.id ? (
                            <div className="flex gap-2">
                              <button onClick={() => updateStock(variant.id, editingStock.stock)} className="text-green-600 text-sm hover:underline">Save</button>
                              <button onClick={() => setEditingStock(null)} className="text-gray-600 text-sm hover:underline">Cancel</button>
                            </div>
                          ) : (
                            <button onClick={() => setEditingStock({ productId: product.id, variantId: variant.id, stock: variant.stock })} className="text-blue-600 text-sm hover:underline">
                              Update Stock
                            </button>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
            <div className="px-6 py-3 bg-gray-50 border-t border-gray-200">
              <div className="flex justify-between items-center">
                <span className="text-sm text-black">Total Stock: {totalStock} units</span>
                <span className={`text-sm font-medium ${stockStatus.textColor}`}>Status: {stockStatus.label}</span>
              </div>
            </div>
          </div>
        );
      })}

      {products.length === 0 && !showAddForm && (
        <div className="text-center py-16 bg-white rounded-xl shadow-md">
          <p className="text-black mb-4">No products found</p>
          <button onClick={() => setShowAddForm(true)} className="bg-black text-white px-6 py-2 rounded-lg">
            Add Your First Product
          </button>
        </div>
      )}

      <div className="grid md:grid-cols-4 gap-4 mt-6">
        <div className="bg-green-50 rounded-lg p-4 border border-green-200">
          <p className="text-sm text-green-700">Total Products</p>
          <p className="text-2xl font-bold text-green-800">{products.length}</p>
        </div>
        <div className="bg-blue-50 rounded-lg p-4 border border-blue-200">
          <p className="text-sm text-blue-700">Total Variants</p>
          <p className="text-2xl font-bold text-blue-800">{products.reduce((s, p) => s + (p.variants?.length || 0), 0)}</p>
        </div>
        <div className="bg-yellow-50 rounded-lg p-4 border border-yellow-200">
          <p className="text-sm text-yellow-700">Low Stock</p>
          <p className="text-2xl font-bold text-yellow-800">{products.reduce((s, p) => s + (p.variants?.filter(v => v.stock > 0 && v.stock < 10).length || 0), 0)}</p>
        </div>
        <div className="bg-red-50 rounded-lg p-4 border border-red-200">
          <p className="text-sm text-red-700">Out of Stock</p>
          <p className="text-2xl font-bold text-red-800">{products.reduce((s, p) => s + (p.variants?.filter(v => v.stock === 0).length || 0), 0)}</p>
        </div>
      </div>
    </div>
  );
}
