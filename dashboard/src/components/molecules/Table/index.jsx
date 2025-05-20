import TableHeader from './tableHeader';
import TableBody from './tableBody';

const Table = ({
  headers,
  data,
  handleRowClick = () => {},
  enableRowClick = false,
  stickyHeader = false,
}) => {
  return (
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
      />
    </table>
  );
};

export default Table;
