import { NextResponse } from 'next/server';

// This function can be marked `async` if using `await` inside
export function middleware(request) {
  const res = NextResponse.next();
  const url = request.url;
  const currentUrl = request.nextUrl.clone();

  if (url === `${currentUrl.origin}/`) {
    return NextResponse.redirect(`${currentUrl.origin}/cloud-accounts`);
  }
  return res;
}

// See "Matching Paths" below to learn more
export const config = {
  matcher: '/',
};
