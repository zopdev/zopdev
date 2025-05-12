import SimpleLoader from '@/components/atom/Loaders/SimpleLoader.jsx';

/**
 * this function is to display loader when query is fetching
 * @returns loader
 */
export default function CompleteLoader() {
  return (
    <div
      className="flex justify-center items-center h-screen w-full"
      data-testid="testCompleteLoader"
    >
      <SimpleLoader size={50} thickness={3} />
    </div>
  );
}
