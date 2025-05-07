'use client';

import Button from '@/components/atom/Button/index.jsx';

const DashBoardCard = ({
  title,
  description,
  features,
  buttonText,
  buttonVariant = 'primary',
  onClick,
  buttonIcon,
}) => {
  return (
    <div
      className="rounded-lg  bg-card text-card-foreground border border-borderDefault flex flex-col"
      data-v0-t="card"
    >
      <div className="flex flex-col space-y-1.5 p-6 pb-3">
        <h3 className="font-semibold tracking-tight text-xl">{title}</h3>
        <p className="text-muted-foreground text-base">{description}</p>
      </div>
      <div className="p-8 pt-0 flex-1">
        <ul className="space-y-2 text-sm">
          {features.map((feature, index) => (
            <li key={index} className="flex items-center">
              <span className="mr-2 flex h-5 w-5 items-center justify-center rounded-full bg-primary/10 text-xs text-primary-600">
                ✓
              </span>
              <span>{feature}</span>
            </li>
          ))}
        </ul>
      </div>
      <div className="flex items-center p-6 pt-0">
        <Button
          onClick={onClick}
          variant={buttonVariant}
          fullWidth={true}
          icon={buttonIcon}
          startEndornment={true}
        >
          {buttonText}
        </Button>
      </div>
    </div>
  );
};

export default DashBoardCard;
