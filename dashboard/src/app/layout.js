import localFont from 'next/font/local';
import './globals.css';
import ErrorCatcher from '../../components/ErrorCatcher/ErrorCatcher';
import NotFound from './not-found';
import dynamic from 'next/dynamic';
import ContextProvider from '../../libs/context';

const TopBarWrapper = dynamic(() => import('../../partials/Header'), {
  // ssr: false,
});

const geistSans = localFont({
  src: './fonts/GeistVF.woff',
  variable: '--font-geist-sans',
  weight: '100 900',
});
const geistMono = localFont({
  src: './fonts/GeistMonoVF.woff',
  variable: '--font-geist-mono',
  weight: '100 900',
});

export const metadata = {
  title: 'Zop.dev: Multicloud infrastructure automation',
  description:
    'Easily provision, manage, and monitor resources across AWS, Google Cloud, and Azure with Zop.dev.',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <ErrorCatcher fallback={<NotFound />}>
          <ContextProvider>
            <TopBarWrapper />
            {/* <div className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8">{children}</div> */}
            <div>{children}</div>
          </ContextProvider>
        </ErrorCatcher>
      </body>
    </html>
  );
}
