import { NextRequest, NextResponse } from 'next/server';

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ product_id: string }> }
) {
  const { product_id } = await params;
  return NextResponse.json({ success: true });
}