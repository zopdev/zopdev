import React from 'react';
import TableHeader from './tableHeader';
import TableBody from './tableBody';

const Table = ({
  headers,
  data,
  onEdit,
  onDelete,
  action,
  handleRowClick = () => {},
  renderComponent,
  enableRowClick,
}) => {
  return (
    <div className="-mx-4 mt-4 min-w-full h-[75vh] overflow-auto flow-root sm:mx-0">
      <table className="table-auto w-full border-collapse">
        <colgroup>
          {headers.map((header) => (
            <col key={header.key} className={header.colClassName || ''} />
          ))}
          <col className=" w-[100vw]" />
        </colgroup>
        <TableHeader headers={headers} action={action} />
        <TableBody
          data={data}
          headers={headers}
          onEdit={onEdit}
          onDelete={onDelete}
          action={action}
          handleRowClick={handleRowClick}
          renderComponent={renderComponent}
          enableRowClick={enableRowClick}
        />
      </table>
    </div>
  );
};

export default Table;
