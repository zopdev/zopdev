import React from 'react';
import { FullScreenDrawer } from '../FullScreenDrawer';
import { DocumentChartBarIcon } from '@heroicons/react/24/outline';

const TableBody = ({
  data,
  headers,
  onEdit,
  onDelete,
  action = true,
  handleRowClick = () => {},
  renderComponent,
  enableRowClick = false,
}) => {
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
          {action && (
            <td className="px-3 py-5 text-right">
              <div className="flex gap-2 justify-center">
                {onEdit && (
                  <button onClick={() => onEdit(row)} className="text-blue-500 hover:text-blue-700">
                    Edit
                  </button>
                )}
                {onDelete && (
                  <button onClick={() => onDelete(row)} className="text-red-500 hover:text-red-700">
                    Delete
                  </button>
                )}
                {renderComponent && (
                  <FullScreenDrawer
                    // tooltipTitle={row?.name}
                    isIcon
                    RenderIcon={
                      <DocumentChartBarIcon
                        className="-ml-0.5 h-6 w-6 text-gray-600"
                        aria-hidden="true"
                      />
                    }
                    title={row?.name}
                    RenderComponent={renderComponent}
                    formData={row}
                  />
                )}
              </div>
            </td>
          )}
        </tr>
      ))}
    </tbody>
  );
};

export default TableBody;
