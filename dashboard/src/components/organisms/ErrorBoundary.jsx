import { Component } from 'react';
import Button from "@/components/atom/Button/index.jsx";

export class ErrorBoundary extends Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  //   componentDidCatch(error, errorInfo) {
  //     console.error('ErrorBoundary caught an error', error, errorInfo);
  //   }

  render() {
    if (this.state.hasError) {
      return (
        <main className="grid min-h-screen place-items-center bg-primary px-6 py-24 sm:py-32 lg:px-8">
          <div className="text-center">
            <p className="text-4xl font-semibold text-primary-600">Aw, Snap!</p>
            <h1 className="mt-4 text-xl font-medium tracking-tight text-secondary-600 sm:text-2xl">
              Something went wrong while displaying this webpage.
            </h1>
            <div className="mt-10 flex items-center justify-center gap-x-6">
              <Button
                onClick={() => {
                  window.location.reload();
                }}
                className="rounded-md bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600"
              >
                Reload
              </Button>
            </div>
          </div>
        </main>
      );
    }

    return this.props.children;
  }
}
