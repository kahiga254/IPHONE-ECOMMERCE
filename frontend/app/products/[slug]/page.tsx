'use client';

import Link from 'next/link';
import Navbar from '@/components/Navbar';
import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import api from '@/services/api';
import toast from 'react-hot-toast';

interface Review {
  id: string;
  user_id: string;
  rating: number;
  comment: string;
  created_at: string;
  user?: {
    name: string;
  };
}

export default function ProductDetailPage() {
  const { slug } = useParams();
  const [product, setProduct] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [addingToWishlist, setAddingToWishlist] = useState(false);
  const [addingToCart, setAddingToCart] = useState(false);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [showReviewForm, setShowReviewForm] = useState(false);
  const [reviewRating, setReviewRating] = useState(5);
  const [reviewComment, setReviewComment] = useState('');
  const [submittingReview, setSubmittingReview] = useState(false);

  useEffect(() => {
    if (slug) {
      fetchProduct();
    }
  }, [slug]);

  const fetchProduct = async () => {
    try {
      const response = await api.get(`/products/${slug}`);
      setProduct(response.data.data);
      if (response.data.data?.id) {
        fetchReviews(response.data.data.id);
      }
    } catch (error) {
      console.error('Failed to fetch product:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchReviews = async (productId: string) => {
    try {
      const response = await api.get(`/reviews/${productId}`);
      setReviews(response.data.data || []);
    } catch (error) {
      console.log('Reviews not available');
      setReviews([]);
    }
  };

  const addToWishlist = async () => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      toast.error('Please login first');
      window.location.href = '/login';
      return;
    }

    setAddingToWishlist(true);
    try {
      await api.post('/wishlist', { variant_id: product.variants[0].id });
      toast.success('Added to wishlist!');
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to add');
    } finally {
      setAddingToWishlist(false);
    }
  };

  const addToCart = () => {
    if (addingToCart) return; // Prevent multiple clicks
    
    setAddingToCart(true);
    
    const existingCart = localStorage.getItem('cart');
    let cart = existingCart ? JSON.parse(existingCart) : [];
    
    const existingIndex = cart.findIndex((item: any) => item.id === product.variants[0].id);
    
    if (existingIndex >= 0) {
      cart[existingIndex].quantity += 1;
      toast.success('Quantity increased in cart!');
    } else {
      cart.push({
        id: product.variants[0].id,
        name: product.name,
        price: product.base_price,
        quantity: 1,
        image: product.variants[0]?.images?.[0] || '',
        variant: `${product.variants[0]?.color} / ${product.variants[0]?.storage}`
      });
      toast.success('Added to cart!');
    }
    
    localStorage.setItem('cart', JSON.stringify(cart));
    window.dispatchEvent(new Event('cartUpdated'));
    
    // Reset button state after a short delay
    setTimeout(() => {
      setAddingToCart(false);
    }, 500);
  };

  const submitReview = async () => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      toast.error('Please login to leave a review');
      window.location.href = '/login';
      return;
    }

    if (reviewComment.length < 10) {
      toast.error('Please write at least 10 characters');
      return;
    }

    setSubmittingReview(true);
    try {
      await api.post('/reviews', {
        product_id: product.id,
        rating: reviewRating,
        comment: reviewComment
      });
      toast.success('Review submitted!');
      setShowReviewForm(false);
      setReviewComment('');
      setReviewRating(5);
      if (product.id) {
        fetchReviews(product.id);
      }
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to submit review');
    } finally {
      setSubmittingReview(false);
    }
  };

  const renderStars = (rating: number) => {
    return '★'.repeat(rating) + '☆'.repeat(5 - rating);
  };

  const averageRating = reviews.length > 0 
    ? reviews.reduce((sum, r) => sum + r.rating, 0) / reviews.length 
    : 0;

  if (loading) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center text-black">Loading...</div>
      </>
    );
  }

  if (!product) {
    return (
      <>
        <Navbar />
        <div className="max-w-7xl mx-auto px-4 py-16 text-center text-black">Product not found</div>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 py-16">
        <Link href="/products" className="text-blue-600 hover:underline mb-4 inline-block">
          ← Back to Products
        </Link>
        
        <div className="grid md:grid-cols-2 gap-8 mt-4">
          <div className="bg-gray-100 rounded-xl p-8">
            {product.variants?.[0]?.images?.[0] && (
              <img 
                src={product.variants[0].images[0]} 
                alt={product.name}
                className="w-full h-96 object-contain"
              />
            )}
          </div>
          <div>
            <h1 className="text-3xl font-bold text-black mb-4">{product.name}</h1>
            <p className="text-gray-600 mb-4">{product.description}</p>
            
            <div className="flex items-center gap-2 mb-4">
              <span className="text-yellow-500 text-xl">{renderStars(Math.round(averageRating))}</span>
              <span className="text-gray-600">({reviews.length} reviews)</span>
            </div>
            
            <p className="text-3xl text-blue-600 font-bold mb-6">
              KSh {product.base_price?.toLocaleString()}
            </p>
            
            <div className="flex gap-4">
              <button 
                onClick={addToCart}
                disabled={addingToCart}
                className={`flex-1 px-8 py-3 rounded-full font-medium transition-all cursor-pointer ${addingToCart ? 'bg-gray-400 cursor-not-allowed' : 'bg-black text-white hover:bg-gray-800'}`}
              >
                {addingToCart ? 'Adding...' : 'Add to Cart'}
              </button>
              <button 
                onClick={addToWishlist}
                disabled={addingToWishlist}
                className={`px-8 py-3 rounded-full font-medium transition-all cursor-pointer border-2 ${addingToWishlist ? 'border-gray-400 text-gray-400 cursor-not-allowed' : 'border-black text-black hover:bg-black hover:text-white'}`}
              >
                {addingToWishlist ? 'Adding...' : 'Add to Wishlist'}
              </button>
            </div>
          </div>
        </div>

        {/* Reviews Section */}
        <div className="mt-16">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-black">Customer Reviews</h2>
            {!showReviewForm && (
              <button
                onClick={() => setShowReviewForm(true)}
                className="bg-black text-white px-4 py-2 rounded-lg hover:bg-gray-800 cursor-pointer"
              >
                Write a Review
              </button>
            )}
          </div>

          {showReviewForm && (
            <div className="bg-white rounded-xl shadow-md p-6 mb-6">
              <h3 className="text-lg font-bold text-black mb-4">Write Your Review</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-black mb-2">Rating</label>
                  <div className="flex gap-2">
                    {[1, 2, 3, 4, 5].map((star) => (
                      <button
                        key={star}
                        onClick={() => setReviewRating(star)}
                        className={`text-2xl cursor-pointer ${star <= reviewRating ? 'text-yellow-500' : 'text-gray-300'}`}
                      >
                        ★
                      </button>
                    ))}
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-black mb-2">Your Review</label>
                  <textarea
                    value={reviewComment}
                    onChange={(e) => setReviewComment(e.target.value)}
                    rows={4}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-black focus:border-transparent text-black"
                    placeholder="Share your experience with this product..."
                  />
                  <p className="text-xs text-gray-500 mt-1">Minimum 10 characters</p>
                </div>
                <div className="flex gap-3">
                  <button
                    onClick={submitReview}
                    disabled={submittingReview}
                    className={`px-6 py-2 rounded-lg font-medium ${submittingReview ? 'bg-gray-400 cursor-not-allowed' : 'bg-black text-white hover:bg-gray-800 cursor-pointer'}`}
                  >
                    {submittingReview ? 'Submitting...' : 'Submit Review'}
                  </button>
                  <button
                    onClick={() => {
                      setShowReviewForm(false);
                      setReviewComment('');
                      setReviewRating(5);
                    }}
                    className="border border-gray-300 text-black px-6 py-2 rounded-lg hover:bg-gray-100 cursor-pointer"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          )}

          {reviews.length === 0 ? (
            <div className="text-center py-8 bg-gray-50 rounded-xl">
              <p className="text-black">No reviews yet. Be the first to review this product!</p>
            </div>
          ) : (
            <div className="space-y-4">
              {reviews.map((review) => (
                <div key={review.id} className="bg-white rounded-xl shadow-md p-6">
                  <div className="flex justify-between items-start mb-2">
                    <div>
                      <p className="font-semibold text-black">{review.user?.name || 'Anonymous'}</p>
                      <div className="text-yellow-500 text-sm">{renderStars(review.rating)}</div>
                    </div>
                    <p className="text-sm text-gray-500">
                      {new Date(review.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <p className="text-gray-700 mt-2">{review.comment}</p>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </>
  );
}
