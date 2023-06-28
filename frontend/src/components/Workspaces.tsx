import React, { useState } from 'react';

interface WorkspaceListProps {
    items: WorkspaceItem[];
    setActive: (id: string) => void;
    newWorkspace: () => void;
}

export type WorkspaceItem = {
    id: string;
    name: string;
    active: boolean;
}

export const WorkspaceList: React.FC<WorkspaceListProps> = ({items, setActive, newWorkspace}) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState<boolean>(false);

  const toggleDropdown = () => {
    setIsDropdownOpen((prevIsDropdownOpen) => !prevIsDropdownOpen);
  };

  const onNewWorkspace = () => {
    if (isDropdownOpen) {
      toggleDropdown();
    }
    newWorkspace();
  }

  const selectedWorkspace = items.find(it => it.active)

  const selectWorkspace = (id: string) => {
    setActive(id);
    setIsDropdownOpen(false);
  };

  return (
    <div className="relative">
      <button
        type="button"
        className="flex items-center px-4 py-2 text-sm font-medium text-gray-800 bg-gray-200 border border-gray-300 rounded-md hover:bg-gray-300 focus:outline-none focus:ring focus:ring-gray-400"
        onClick={toggleDropdown}
      >
        {selectedWorkspace && (<span className="mr-1">{selectedWorkspace.name}</span>)}
        <svg
          className={`w-4 h-4 transition-transform duration-300 ${
            isDropdownOpen ? 'transform rotate-180' : ''
          }`}
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            fillRule="evenodd"
            d="M10.707 14.293a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L10 11.586l3.293-3.293a1 1 0 111.414 1.414l-4 4z"
            clipRule="evenodd"
          />
        </svg>
      </button>

      {isDropdownOpen && (
        <div className="absolute top-full z-10 w-40 mt-2 bg-white border border-gray-300 rounded-md shadow-lg">
          <ul className="py-1">
            {items.map((item, index) => (
              <li key={index} className="px-4 py-2 text-sm text-gray-800 hover:bg-gray-100 cursor-pointer" onClick={() => selectWorkspace(item.id)}>{item.name}</li>
            ))}
          </ul>
        </div>
      )}
      <div>
        <button onClick={onNewWorkspace}>Add new workspace</button>
      </div>
    </div>
  );
};
