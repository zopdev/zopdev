const TableHeader = ({ headers, sticky = false }) => {
  return (
    <thead
      className={`border-b border-gray-300 text-gray-900 ${
        sticky ? 'sticky top-0 z-10 bg-white' : ''
      }`}
    >
      <tr>
        {headers.map((header) => (
          <th
            key={header.key}
            scope="col"
            className={`px-3 py-3.5 text-sm font-semibold whitespace-nowrap ${
              header.align === 'right' ? 'text-right' : 'text-left'
            } ${header.className || ''} ${sticky ? 'bg-white' : ''}`} // <== Add this
          >
            {header.label}
          </th>
        ))}
      </tr>
    </thead>
  );
};

export default TableHeader;
