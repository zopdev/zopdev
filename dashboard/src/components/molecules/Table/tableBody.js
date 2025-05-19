const TableBody = ({ data, headers, handleRowClick = () => {}, enableRowClick = false }) => {
  return (
    <tbody>
      {data?.map((row, rowIndex) => (
        <tr
          key={row.id || rowIndex}
          className={`border-b border-gray-200 group ${enableRowClick && 'hover:bg-gray-50 cursor-pointer'}`}
          onClick={() => enableRowClick && handleRowClick(row)}
        >
          {headers.map((header) => (
            <td
              key={header.key}
              className={`px-3 py-5 text-sm  ${
                header.align === 'right' ? 'text-right' : 'text-left'
              } text-gray-500 ${header.colClassName}`}
              style={{ minWidth: 175, flexGrow: 1 }}
            >
              {/* Allow React components or plain text */}
              {typeof row[header.key] === 'function' ? row[header.key]() : row[header.key]}
            </td>
          ))}
        </tr>
      ))}
    </tbody>
  );
};

export default TableBody;
