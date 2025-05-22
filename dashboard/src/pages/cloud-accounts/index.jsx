import BlankCloudAccountSvg from '@/assets/svg/BlankCloudAccount';
import Button from '@/components/atom/Button';
import EmptyComponent from '@/components/atom/EmptyComponent';
import ErrorComponent from '@/components/atom/ErrorComponent';
import CompleteLoader from '@/components/atom/Loaders/CompleteLoader';
import PageHeading from '@/components/atom/PageHeading';
import CloudAccountCard from '@/components/molecules/Cards/CloudAccountCard';
import { useGetCloudAccounts } from '@/queries/cloud-account';
import { PlusCircleIcon } from '@heroicons/react/20/solid';

const CloudAccountPage = () => {
  const cloudAccounts = useGetCloudAccounts();

  return (
    <div className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8">
      {cloudAccounts?.data?.data?.length > 0 && (
        <PageHeading
          title={'Cloud Accounts'}
          actions={
            <Button
              href={'/cloud-setup'}
              size="md"
              startEndornment={<PlusCircleIcon className="-ml-0.5 h-5 w-5" aria-hidden="true" />}
            >
              Add Cloud Account
            </Button>
          }
        />
      )}

      {cloudAccounts?.isLoading && <CompleteLoader />}
      <div className="flex gap-4 w-full justify-start mt-4 flex-wrap">
        {cloudAccounts?.data?.data?.map((item, idx) => {
          return <CloudAccountCard key={idx} item={item} />;
        })}
      </div>

      {cloudAccounts?.data?.data?.length === 0 && (
        <EmptyComponent
          imageComponent={<BlankCloudAccountSvg />}
          redirectLink={'/cloud-setup'}
          buttonTitle={'Add Cloud Account'}
          title={'Please start by setting up your first cloud account'}
        />
      )}

      {cloudAccounts?.isError && (
        <ErrorComponent errorText={cloudAccounts?.error?.message || 'Something went wrong'} />
      )}
    </div>
  );
};

export default CloudAccountPage;
