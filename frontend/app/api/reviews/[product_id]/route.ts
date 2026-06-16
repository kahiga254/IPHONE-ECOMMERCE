import { NextRequest, NextResponse } from 'next/server';

export async function GET(
  request: NextRequest,
  { params }: { params: { product_id: string } }
) {
  // This is handled by the backend API via the service
  // The frontend calls api.get(`/reviews/${product_id}`)
  return NextResponse.json({ success: true });
}
