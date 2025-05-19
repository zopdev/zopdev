const TableHeader = ({ headers, action = true }) => {
  return (
    <thead className="border-b border-gray-300 text-gray-900">
      <tr>
        {headers.map((header) => (
          <th
            key={header.key}
            scope="col"
            className={`px-3 py-3.5 text-sm font-semibold text-gray-900 whitespace-nowrap ${
              header.align === 'right' ? 'text-right' : 'text-left'
            } ${header.className || ''}`}
          >
            {header.label}
          </th>
        ))}
        {action && (
          <th className="px-3 py-3.5 text-sm font-semibold text-gray-900 text-center">Actions</th>
        )}
      </tr>
    </thead>
  );
};

export default TableHeader;
