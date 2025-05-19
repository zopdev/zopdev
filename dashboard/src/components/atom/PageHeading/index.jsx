const PageHeading = ({ title, subTitle, titleAction, actions, subTitleAction }) => {
  return (
    <div className="w-full flex mb-3 items-center justify-between">
      <div>
        <div className={`${titleAction ? 'flex gap-1' : ''}`}>
          <p className="text-left font-medium text-gray-600 text-xl">{title}</p>
          {titleAction && titleAction}
        </div>
        {subTitle && <div className="pt-0.5 text-gray-400 text-left text-sm">{subTitle}</div>}
        {subTitleAction && <>{subTitleAction}</>}
      </div>
      {actions && <div className="flex flex-row gap-2 items-center justify-end">{actions}</div>}
    </div>
  );
};

export default PageHeading;
