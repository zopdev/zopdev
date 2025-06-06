import TableHeader from './tableHeader';
import TableBody from './tableBody';
import { MagnifyingGlassIcon } from '@heroicons/react/20/solid';

const Table = ({
  headers,
  data,
  handleRowClick = () => {},
  enableRowClick = false,
  stickyHeader = true,
  emptyStateTitle,
  emptyStateDescription,
  renderRow,
  maxHeight,
}) => {
  return (
    <>
      <div className="overflow-x-auto overflow-y-auto scroll-hidden" style={{ maxHeight }}>
        <table className="table-auto w-full border-collapse">
          <colgroup>
            {headers.map((header) => (
              <col
                key={header.key}
                style={{ width: header.width || 'auto' }}
                className={header.colClassName || ''}
              />
            ))}
          </colgroup>
          <TableHeader headers={headers} sticky={stickyHeader} />
          <TableBody
            headers={headers}
            data={data}
            handleRowClick={handleRowClick}
            enableRowClick={enableRowClick}
            renderRow={renderRow}
          />
        </table>
      </div>
      {data?.length === 0 && (
        <div className="w-full  flex flex-col items-center border-t border-gray-200 pt-20">
          <MagnifyingGlassIcon className="size-20  fill-gray-400  mb-4" />
          <p className="text-gray-400 font-semibold text-2xl">
            {emptyStateTitle || 'No data found'}
          </p>
          <p className="text-gray-400">{emptyStateDescription && emptyStateDescription}</p>
        </div>
      )}
    </>
  );
};

export default Table;
