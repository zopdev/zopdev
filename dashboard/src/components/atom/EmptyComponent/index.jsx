import { PlusCircleIcon } from '@heroicons/react/20/solid';
import Button from '@/components/atom/Button';

const EmptyComponent = ({
  imageComponent,
  redirectLink,
  buttonTitle,
  title,
  buttonIcon,
  height = 'h-[70vh]',
  customButton,
  onClickHandler,
}) => {
  return (
    <div className={`w-full flex justify-center items-center ${height}`}>
      <div className="flex flex-col gap-4 items-center">
        <div className="flex flex-col justify-center items-center">
          {imageComponent}
          <p className="text-gray-400 text-base font-medium text-wrap mt-2 text-center">{title}</p>
        </div>
        {redirectLink && (
          <Button
            href={redirectLink}
            startEndornment={
              buttonIcon || <PlusCircleIcon className="h-5 w-5" aria-hidden="true" />
            }
          >
            {buttonTitle}
          </Button>
        )}
        {onClickHandler && (
          <Button
            onClick={onClickHandler}
            startEndornment={
              buttonIcon || <PlusCircleIcon className="h-5 w-5" aria-hidden="true" />
            }
          >
            {buttonTitle}
          </Button>
        )}
        {customButton}
      </div>
    </div>
  );
};

export default EmptyComponent;
