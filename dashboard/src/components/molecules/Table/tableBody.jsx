const TableBody = ({ data, headers, handleRowClick, enableRowClick }) => {
  return (
    <tbody>
      {data?.map((row, rowIndex) => (
        <tr
          key={row.id || rowIndex}
          className={`border-b border-gray-200 group even:bg-gray-50  ${
            enableRowClick ? 'hover:bg-gray-100 cursor-pointer' : ''
          }`}
          onClick={() => enableRowClick && handleRowClick(row)}
        >
          {headers.map((header) => (
            <td
              key={header.key}
              className={`px-3 py-5 text-sm text-gray-500 ${
                header.align === 'right' ? 'text-right' : 'text-left'
              } ${header.colClassName || ''}`}
            >
              {typeof row[header.key] === 'function' ? row[header.key]() : row[header.key]}
            </td>
          ))}
        </tr>
      ))}
    </tbody>
  );
};

export default TableBody;
