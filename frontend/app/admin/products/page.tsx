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
  description: string;
  base_price: number;
  is_active: boolean;
  variants: Variant[];
}

interface Category {
  id: string;
  name: string;
  slug: string;
}

export default function AdminProducts() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingStock, setEditingStock] = useState<{ productId: string; variantId: string; stock: number } | null>(null);
  const [showAddForm, setShowAddForm] = useState(false);
  const [editingProductId, setEditingProductId] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 20,
    total: 0,
    totalPages: 0,
  });
  const [newProduct, setNewProduct] = useState({
    name: '',
    slug: '',
    description: '',
    base_price: '',
    category_id: '',
    variants: [{ sku: '', color: '', storage: '', price: '', stock: '', images: [] as string[] }]
  });

  const excludedCategories = ['apple', 'samsung', 'electronics'];

  useEffect(() => {
    fetchCategories();
  }, []);

  useEffect(() => {
    fetchProducts();
  }, [pagination.page, selectedCategory]);

  const fetchCategories = async () => {
    try {
      const response = await api.get('/categories');
      const allCategories = response.data.data?.categories || [];
      const filtered = allCategories.filter(
        (cat: Category) => !excludedCategories.includes(cat.slug)
      );
      setCategories(filtered);
    } catch (error) {
      console.error('Failed to fetch categories:', error);
    }
  };

  const fetchProducts = async () => {
    setLoading(true);
    try {
      let url = `/products?page=${pagination.page}&limit=${pagination.limit}`;
      if (selectedCategory !== 'all') {
        url += `&category=${selectedCategory}`;
      }
      const response = await api.get(url);
      const data = response.data.data;
      setProducts(data.data || []);
      setPagination({
        ...pagination,
        total: data.total || 0,
        totalPages: data.total_pages || 0,
      });
    } catch (error) {
      console.error('Failed to fetch products:', error);
      toast.error('Failed to load products');
    } finally {
      setLoading(false);
    }
  };

  const goToPage = (page: number) => {
    setPagination({ ...pagination, page });
  };

  const updateStock = async (variantId: string, newStock: number) => {
    try {
      await api.patch(`/admin/variants/${variantId}/stock`, { stock: newStock });
      toast.success('Stock updated');
      fetchProducts();
      setEditingStock(null);
    } catch (error) {
      toast.error('Failed to update stock');
    }
  };

  const getStockStatus = (stock: number) => {
    if (stock === 0) return { label: 'Out', color: 'bg-red-100 text-red-700' };
    if (stock < 10) return { label: 'Low', color: 'bg-yellow-100 text-yellow-700' };
    return { label: 'In', color: 'bg-green-100 text-green-700' };
  };

  const toggleProductStatus = async (id: string, currentStatus: boolean) => {
    try {
      await api.patch(`/admin/products/${id}/status`, { is_active: !currentStatus });
      toast.success(`Product ${!currentStatus ? 'activated' : 'deactivated'}`);
      fetchProducts();
    } catch (error) {
      toast.error('Failed to update');
    }
  };

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>, variantIndex: number, productId: string | null = null) => {
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
          if (productId) {
            // Editing existing product
            const updatedProducts = products.map(p => {
              if (p.id === productId) {
                const newVariants = [...p.variants];
                newVariants[variantIndex].images.push(data.url);
                return { ...p, variants: newVariants };
              }
              return p;
            });
            setProducts(updatedProducts);
          } else {
            // Adding new product
            const variants = [...newProduct.variants];
            variants[variantIndex].images.push(data.url);
            setNewProduct({ ...newProduct, variants });
          }
          toast.success('Image uploaded');
        }
      } catch (error) {
        toast.error('Upload failed');
      }
    }

    setUploading(false);
  };

  const removeImage = (variantIndex: number, imageIndex: number, productId: string | null = null) => {
    if (productId) {
      const updatedProducts = products.map(p => {
        if (p.id === productId) {
          const newVariants = [...p.variants];
          newVariants[variantIndex].images.splice(imageIndex, 1);
          return { ...p, variants: newVariants };
        }
        return p;
      });
      setProducts(updatedProducts);
    } else {
      const variants = [...newProduct.variants];
      variants[variantIndex].images.splice(imageIndex, 1);
      setNewProduct({ ...newProduct, variants });
    }
  };

  const addVariant = (productId: string | null = null) => {
    if (productId) {
      const updatedProducts = products.map(p => {
        if (p.id === productId) {
          return {
            ...p,
            variants: [...p.variants, { id: '', sku: '', color: '', storage: '', price: 0, stock: 0, images: [] }]
          };
        }
        return p;
      });
      setProducts(updatedProducts);
    } else {
      setNewProduct({
        ...newProduct,
        variants: [...newProduct.variants, { sku: '', color: '', storage: '', price: '', stock: '', images: [] }]
      });
    }
  };

  const removeVariant = (index: number, productId: string | null = null) => {
    if (productId) {
      const updatedProducts = products.map(p => {
        if (p.id === productId) {
          const newVariants = [...p.variants];
          newVariants.splice(index, 1);
          return { ...p, variants: newVariants };
        }
        return p;
      });
      setProducts(updatedProducts);
    } else {
      const variants = [...newProduct.variants];
      variants.splice(index, 1);
      setNewProduct({ ...newProduct, variants });
    }
  };

  const updateVariant = (index: number, field: string, value: string, productId: string | null = null) => {
    if (productId) {
      const updatedProducts = products.map(p => {
        if (p.id === productId) {
          const newVariants = [...p.variants];
          newVariants[index] = { ...newVariants[index], [field]: field === 'price' || field === 'stock' ? parseFloat(value) || 0 : value };
          return { ...p, variants: newVariants };
        }
        return p;
      });
      setProducts(updatedProducts);
    } else {
      const variants = [...newProduct.variants];
      variants[index] = { ...variants[index], [field]: value };
      setNewProduct({ ...newProduct, variants });
    }
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
      category_id: newProduct.category_id,
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
      toast.success('Product created');
      setShowAddForm(false);
      setNewProduct({
        name: '',
        slug: '',
        description: '',
        base_price: '',
        category_id: '',
        variants: [{ sku: '', color: '', storage: '', price: '', stock: '', images: [] }]
      });
      fetchProducts();
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to create');
    }
  };

  const updateProduct = async (productId: string) => {
    const product = products.find(p => p.id === productId);
    if (!product) return;

    const productData = {
      name: product.name,
      description: product.description,
      base_price: product.base_price,
      is_active: product.is_active,
    };

    try {
      await api.put(`/admin/products/${productId}`, productData);
      
      for (const variant of product.variants) {
        if (variant.id) {
          await api.put(`/admin/variants/${variant.id}`, {
            sku: variant.sku,
            color: variant.color,
            storage: variant.storage,
            price: variant.price,
            stock: variant.stock,
            images: variant.images
          });
        }
      }
      
      toast.success('Product updated');
      setEditingProductId(null);
      fetchProducts();
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to update');
    }
  };

  const openEditForm = (productId: string) => {
    setEditingProductId(editingProductId === productId ? null : productId);
  };

  if (loading) {
    return <div className="text-center py-16 text-black">Loading...</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-xl font-bold text-black">Products ({pagination.total})</h1>
        <button
          onClick={() => setShowAddForm(true)}
          className="bg-black text-white px-3 py-1.5 rounded-lg text-sm hover:bg-gray-800"
        >
          + Add Product
        </button>
      </div>

      {/* Category Filter */}
      <div className="flex gap-2 mb-4 flex-wrap">
        <button
          onClick={() => setSelectedCategory('all')}
          className={`px-3 py-1.5 rounded-lg text-sm font-medium transition ${
            selectedCategory === 'all' 
              ? 'bg-black text-white' 
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          }`}
        >
          All Products
        </button>
        {categories.map((cat) => (
          <button
            key={cat.id}
            onClick={() => setSelectedCategory(cat.slug)}
            className={`px-3 py-1.5 rounded-lg text-sm font-medium transition ${
              selectedCategory === cat.slug 
                ? 'bg-black text-white' 
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            {cat.name}
          </button>
        ))}
      </div>

      {/* Add Product Form */}
      {showAddForm && (
        <div className="bg-white rounded-xl shadow-md p-4 mb-4">
          <h2 className="text-lg font-bold text-black mb-3">Add Product</h2>
          <div className="space-y-3">
            <div className="grid md:grid-cols-3 gap-3">
              <input type="text" placeholder="Product Name *" value={newProduct.name} onChange={(e) => setNewProduct({ ...newProduct, name: e.target.value })} className="px-3 py-1.5 border rounded-lg text-sm text-black" />
              <input type="number" placeholder="Base Price *" value={newProduct.base_price} onChange={(e) => setNewProduct({ ...newProduct, base_price: e.target.value })} className="px-3 py-1.5 border rounded-lg text-sm text-black" />
              <select value={newProduct.category_id} onChange={(e) => setNewProduct({ ...newProduct, category_id: e.target.value })} className="px-3 py-1.5 border rounded-lg text-sm text-black">
                <option value="">Select Category</option>
                {categories.map((cat) => (
                  <option key={cat.id} value={cat.id}>{cat.name}</option>
                ))}
              </select>
            </div>
            <textarea placeholder="Description" value={newProduct.description} onChange={(e) => setNewProduct({ ...newProduct, description: e.target.value })} rows={2} className="w-full px-3 py-1.5 border rounded-lg text-sm text-black" />
            
            <h3 className="font-semibold text-black text-sm">Variants</h3>
            {newProduct.variants.map((variant, idx) => (
              <div key={idx} className="p-3 bg-gray-50 rounded-lg border border-gray-200">
                <div className="grid md:grid-cols-5 gap-2">
                  <input type="text" placeholder="SKU" value={variant.sku} onChange={(e) => updateVariant(idx, 'sku', e.target.value)} className="px-2 py-1 border rounded text-sm text-black" />
                  <input type="text" placeholder="Color" value={variant.color} onChange={(e) => updateVariant(idx, 'color', e.target.value)} className="px-2 py-1 border rounded text-sm text-black" />
                  <input type="text" placeholder="Storage" value={variant.storage} onChange={(e) => updateVariant(idx, 'storage', e.target.value)} className="px-2 py-1 border rounded text-sm text-black" />
                  <input type="number" placeholder="Price" value={variant.price} onChange={(e) => updateVariant(idx, 'price', e.target.value)} className="px-2 py-1 border rounded text-sm text-black" />
                  <input type="number" placeholder="Stock" value={variant.stock} onChange={(e) => updateVariant(idx, 'stock', e.target.value)} className="px-2 py-1 border rounded text-sm text-black" />
                </div>
                <div className="mt-2">
                  <input type="file" multiple accept="image/*" onChange={(e) => handleImageUpload(e, idx)} disabled={uploading} className="text-xs" />
                  {variant.images.length > 0 && (
                    <div className="flex gap-1 mt-1 flex-wrap">
                      {variant.images.map((img, imgIdx) => (
                        <div key={imgIdx} className="relative">
                          <img src={img} className="w-8 h-8 object-cover rounded border" />
                          <button onClick={() => removeImage(idx, imgIdx)} className="absolute -top-1 -right-1 bg-red-500 text-white rounded-full w-4 h-4 text-xs">×</button>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
                {idx > 0 && <button onClick={() => removeVariant(idx)} className="text-red-600 text-xs mt-1">Remove</button>}
              </div>
            ))}
            <button onClick={() => addVariant()} className="text-blue-600 text-xs">+ Add Variant</button>
            
            <div className="flex gap-2 pt-2">
              <button onClick={createProduct} className="bg-black text-white px-4 py-1.5 rounded-lg text-sm">Create</button>
              <button onClick={() => setShowAddForm(false)} className="border px-4 py-1.5 rounded-lg text-sm">Cancel</button>
            </div>
          </div>
        </div>
      )}

      {/* Products List */}
      {products.map((product) => {
        const totalStock = product.variants?.reduce((sum, v) => sum + v.stock, 0) || 0;
        const stockStatus = getStockStatus(totalStock);
        const isEditing = editingProductId === product.id;
        
        return (
          <div key={product.id} className="bg-white rounded-lg shadow-sm mb-2 overflow-hidden border border-gray-100">
            {/* Product Header - Slim */}
            <div className="px-4 py-2 bg-gray-50 flex justify-between items-center border-b border-gray-100">
              <div className="flex items-center gap-4">
                <span className="font-semibold text-black text-sm">{product.name}</span>
                <span className="text-xs text-gray-500">{product.slug}</span>
                <span className={`px-1.5 py-0.5 rounded text-xs font-medium ${product.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'}`}>
                  {product.is_active ? 'Active' : 'Inactive'}
                </span>
                <button onClick={() => toggleProductStatus(product.id, product.is_active)} className="text-blue-600 text-xs hover:underline">
                  {product.is_active ? 'Deactivate' : 'Activate'}
                </button>
              </div>
              <div className="flex items-center gap-3">
                <button onClick={() => openEditForm(product.id)} className="text-blue-600 text-xs hover:underline">
                  {isEditing ? 'Cancel Edit' : 'Edit'}
                </button>
                <span className="text-sm font-bold text-black">KSh {product.base_price?.toLocaleString()}</span>
                <span className={`text-xs font-medium ${stockStatus.color}`}>{stockStatus.label}</span>
              </div>
            </div>

            {/* Edit Form - Inline */}
            {isEditing && (
              <div className="p-4 bg-white border-b border-gray-100">
                <h3 className="font-semibold text-black mb-3">Editing: {product.name}</h3>
                <div className="space-y-3">
                  <div className="grid md:grid-cols-2 gap-3">
                    <div>
                      <label className="block text-xs text-gray-600 mb-1">Product Name</label>
                      <input 
                        type="text" 
                        value={product.name} 
                        onChange={(e) => {
                          const updated = products.map(p => p.id === product.id ? {...p, name: e.target.value} : p);
                          setProducts(updated);
                        }}
                        className="w-full px-3 py-1.5 border rounded-lg text-sm text-black"
                      />
                    </div>
                    <div>
                      <label className="block text-xs text-gray-600 mb-1">Base Price</label>
                      <input 
                        type="number" 
                        value={product.base_price} 
                        onChange={(e) => {
                          const updated = products.map(p => p.id === product.id ? {...p, base_price: parseFloat(e.target.value) || 0} : p);
                          setProducts(updated);
                        }}
                        className="w-full px-3 py-1.5 border rounded-lg text-sm text-black"
                      />
                    </div>
                  </div>
                  <div>
                    <label className="block text-xs text-gray-600 mb-1">Description</label>
                    <textarea 
                      value={product.description || ''} 
                      onChange={(e) => {
                        const updated = products.map(p => p.id === product.id ? {...p, description: e.target.value} : p);
                        setProducts(updated);
                      }}
                      rows={2}
                      className="w-full px-3 py-1.5 border rounded-lg text-sm text-black"
                    />
                  </div>
                  <div className="flex items-center gap-3">
                    <label className="text-sm text-black">Active:</label>
                    <input 
                      type="checkbox" 
                      checked={product.is_active} 
                      onChange={(e) => {
                        const updated = products.map(p => p.id === product.id ? {...p, is_active: e.target.checked} : p);
                        setProducts(updated);
                      }}
                      className="w-4 h-4"
                    />
                  </div>
                  
                  <h4 className="font-semibold text-black text-sm mt-3">Variants</h4>
                  {product.variants.map((variant, idx) => (
                    <div key={idx} className="p-3 bg-gray-50 rounded-lg border border-gray-200">
                      <div className="grid md:grid-cols-5 gap-2">
                        <div>
                          <label className="block text-xs text-gray-600 mb-1">SKU</label>
                          <input 
                            type="text" 
                            value={variant.sku} 
                            onChange={(e) => updateVariant(idx, 'sku', e.target.value, product.id)}
                            className="w-full px-2 py-1 border rounded text-sm text-black"
                          />
                        </div>
                        <div>
                          <label className="block text-xs text-gray-600 mb-1">Color</label>
                          <input 
                            type="text" 
                            value={variant.color} 
                            onChange={(e) => updateVariant(idx, 'color', e.target.value, product.id)}
                            className="w-full px-2 py-1 border rounded text-sm text-black"
                          />
                        </div>
                        <div>
                          <label className="block text-xs text-gray-600 mb-1">Storage</label>
                          <input 
                            type="text" 
                            value={variant.storage} 
                            onChange={(e) => updateVariant(idx, 'storage', e.target.value, product.id)}
                            className="w-full px-2 py-1 border rounded text-sm text-black"
                          />
                        </div>
                        <div>
                          <label className="block text-xs text-gray-600 mb-1">Price</label>
                          <input 
                            type="number" 
                            value={variant.price} 
                            onChange={(e) => updateVariant(idx, 'price', e.target.value, product.id)}
                            className="w-full px-2 py-1 border rounded text-sm text-black"
                          />
                        </div>
                        <div>
                          <label className="block text-xs text-gray-600 mb-1">Stock</label>
                          <input 
                            type="number" 
                            value={variant.stock} 
                            onChange={(e) => updateVariant(idx, 'stock', e.target.value, product.id)}
                            className="w-full px-2 py-1 border rounded text-sm text-black"
                          />
                        </div>
                      </div>
                      <div className="mt-2">
                        <label className="block text-xs text-gray-600 mb-1">Images</label>
                        <input 
                          type="file" 
                          multiple 
                          accept="image/*" 
                          onChange={(e) => handleImageUpload(e, idx, product.id)} 
                          disabled={uploading} 
                          className="text-xs"
                        />
                        {variant.images.length > 0 && (
                          <div className="flex gap-1 mt-1 flex-wrap">
                            {variant.images.map((img, imgIdx) => (
                              <div key={imgIdx} className="relative">
                                <img src={img} className="w-8 h-8 object-cover rounded border" />
                                <button 
                                  onClick={() => removeImage(idx, imgIdx, product.id)} 
                                  className="absolute -top-1 -right-1 bg-red-500 text-white rounded-full w-4 h-4 text-xs"
                                >
                                  ×
                                </button>
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                      {idx > 0 && (
                        <button 
                          onClick={() => removeVariant(idx, product.id)} 
                          className="text-red-600 text-xs mt-1"
                        >
                          Remove Variant
                        </button>
                      )}
                    </div>
                  ))}
                  <button onClick={() => addVariant(product.id)} className="text-blue-600 text-xs">
                    + Add Variant
                  </button>
                  
                  <div className="flex gap-2 pt-2">
                    <button 
                      onClick={() => updateProduct(product.id)} 
                      className="bg-black text-white px-4 py-1.5 rounded-lg text-sm"
                    >
                      Save Changes
                    </button>
                    <button 
                      onClick={() => setEditingProductId(null)} 
                      className="border px-4 py-1.5 rounded-lg text-sm"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              </div>
            )}

            {/* Variants Table */}
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead className="bg-gray-50">
                  <tr className="text-left text-xs">
                    <th className="px-3 py-1.5 text-gray-600">SKU</th>
                    <th className="px-3 py-1.5 text-gray-600">Color</th>
                    <th className="px-3 py-1.5 text-gray-600">Storage</th>
                    <th className="px-3 py-1.5 text-gray-600">Price</th>
                    <th className="px-3 py-1.5 text-gray-600">Stock</th>
                    <th className="px-3 py-1.5 text-gray-600">Status</th>
                    <th className="px-3 py-1.5 text-gray-600">Action</th>
                  </tr>
                </thead>
                <tbody>
                  {product.variants?.map((variant) => {
                    const vs = getStockStatus(variant.stock);
                    return (
                      <tr key={variant.id} className="border-b border-gray-50">
                        <td className="px-3 py-1.5 text-black text-xs">{variant.sku}</td>
                        <td className="px-3 py-1.5 text-black text-xs">{variant.color}</td>
                        <td className="px-3 py-1.5 text-black text-xs">{variant.storage}</td>
                        <td className="px-3 py-1.5 text-black text-xs">KSh {variant.price?.toLocaleString()}</td>
                        <td className="px-3 py-1.5">
                          {editingStock?.variantId === variant.id ? (
                            <input type="number" value={editingStock.stock} onChange={(e) => setEditingStock({ ...editingStock, stock: parseInt(e.target.value) || 0 })} className="w-16 px-1 py-0.5 border rounded text-xs" />
                          ) : (
                            <span className={`text-xs font-medium ${vs.color}`}>{variant.stock}</span>
                          )}
                        </td>
                        <td className="px-3 py-1.5"><span className={`px-1.5 py-0.5 rounded text-xs ${vs.color}`}>{vs.label}</span></td>
                        <td className="px-3 py-1.5">
                          {editingStock?.variantId === variant.id ? (
                            <div className="flex gap-1">
                              <button onClick={() => updateStock(variant.id, editingStock.stock)} className="text-green-600 text-xs">Save</button>
                              <button onClick={() => setEditingStock(null)} className="text-gray-600 text-xs">Cancel</button>
                            </div>
                          ) : (
                            <button onClick={() => setEditingStock({ productId: product.id, variantId: variant.id, stock: variant.stock })} className="text-blue-600 text-xs">Update</button>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
            <div className="px-3 py-1 bg-gray-50 text-xs text-gray-500 flex justify-between">
              <span>Total Stock: {totalStock} units</span>
              <span>Status: {stockStatus.label}</span>
            </div>
          </div>
        );
      })}

      {/* Pagination */}
      {pagination.totalPages > 1 && (
        <div className="flex justify-center gap-2 mt-4">
          <button
            onClick={() => goToPage(pagination.page - 1)}
            disabled={pagination.page === 1}
            className="px-3 py-1 border rounded text-sm disabled:opacity-50 hover:bg-gray-100"
          >
            Previous
          </button>
          <span className="px-3 py-1 text-sm text-black">
            Page {pagination.page} of {pagination.totalPages}
          </span>
          <button
            onClick={() => goToPage(pagination.page + 1)}
            disabled={pagination.page === pagination.totalPages}
            className="px-3 py-1 border rounded text-sm disabled:opacity-50 hover:bg-gray-100"
          >
            Next
          </button>
        </div>
      )}

      {products.length === 0 && !showAddForm && (
        <div className="text-center py-8 bg-white rounded-lg shadow-sm">
          <p className="text-black text-sm">No products found in this category</p>
        </div>
      )}

      {/* Stats */}
      <div className="grid grid-cols-4 gap-3 mt-4">
        <div className="bg-green-50 p-2 rounded-lg text-center">
          <p className="text-sm font-bold text-green-700">{pagination.total}</p>
          <p className="text-xs text-green-600">Products</p>
        </div>
        <div className="bg-blue-50 p-2 rounded-lg text-center">
          <p className="text-sm font-bold text-blue-700">{products.reduce((s, p) => s + (p.variants?.length || 0), 0)}</p>
          <p className="text-xs text-blue-600">Variants</p>
        </div>
        <div className="bg-yellow-50 p-2 rounded-lg text-center">
          <p className="text-sm font-bold text-yellow-700">{products.reduce((s, p) => s + (p.variants?.filter(v => v.stock > 0 && v.stock < 10).length || 0), 0)}</p>
          <p className="text-xs text-yellow-600">Low Stock</p>
        </div>
        <div className="bg-red-50 p-2 rounded-lg text-center">
          <p className="text-sm font-bold text-red-700">{products.reduce((s, p) => s + (p.variants?.filter(v => v.stock === 0).length || 0), 0)}</p>
          <p className="text-xs text-red-600">Out of Stock</p>
        </div>
      </div>
    </div>
  );
}
