import { Suspense } from 'react';
import Navbar from '@/components/Navbar';
import ProductsContent from './ProductsContent';

export const dynamic = 'force-dynamic';

export default function ProductsPage() {
  return (
    <>
      <Navbar />
      <Suspense fallback={
        <div className="max-w-7xl mx-auto px-4 py-16 text-center text-black">
          Loading products...
        </div>
      }>
        <ProductsContent />
      </Suspense>
    </>
  );
}
