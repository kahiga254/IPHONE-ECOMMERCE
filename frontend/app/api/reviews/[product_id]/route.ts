import { NextRequest, NextResponse } from 'next/server';

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ product_id: string }> }
) {

  const { product_id } = await params;
  // This is handled by the backend API via the service
  // The frontend calls api.get(`/reviews/${product_id}`)
  return NextResponse.json({ success: true });
}
