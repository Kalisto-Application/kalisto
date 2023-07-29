import React, { useState, useContext } from 'react';
import { Context } from "../state";

import { NewWorkspace, GetWorkspace, DeleteWorkspace } from "../../wailsjs/go/api/Api";
import { models } from "../../wailsjs/go/models";

interface WorkspaceListProps {
    items: models.Workspace[];
    activeWorkspace?: models.Workspace;
}

export const WorkspaceList: React.FC<WorkspaceListProps> = ({items, activeWorkspace}) => {
  const ctx = useContext(Context);

  const newWorkspace = () => {
    NewWorkspace()
    .then(res => {
      ctx.dispatch({type: 'newWorkspace', workspace: res});
    })
    .catch(err => console.log('error on new workspace: ', err))
  }

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

  return (
    <div className="relative">
      <button
        type="button"
        className="flex items-center px-4 py-2 text-sm font-medium text-gray-800 bg-gray-200 border border-gray-300 rounded-md hover:bg-gray-300 focus:outline-none focus:ring focus:ring-gray-400"
        onClick={toggleDropdown}
      >
        {activeWorkspace && (<span className="mr-1">{activeWorkspace.name}</span>)}
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

          { items && (<ul className="py-1">
            {items.map((item, index) => (
              <WorkspaceItem key={index} id={item.id} onClick={()=> {setIsDropdownOpen(false)}} name={item.name} />
            ))}
          </ul>)}

          <button onClick={onNewWorkspace}>Add new workspace</button>
        </div>
      )}
    </div>
  );
};

type WorkspaceItemProps = {
  id: string;
  onClick: () => void;
  name: string;
}

const WorkspaceItem: React.FC<WorkspaceItemProps> = ({ id, onClick, name}) => {
  const ctx = useContext(Context);

  const [isHovered, setIsHovered] = useState(false);

  const onMouseEnter = () => {
    setIsHovered(true);
  }

  const onMouseLeave = () => {
    setIsHovered(false);
  }

  const removeWorkspace = (id: string) => {
    onClick();
    DeleteWorkspace(id).then(_ => {
      ctx.dispatch({type: 'removeWorkspace', id: id})
    }).catch(err => {
      console.log(`failed to remove workspace id=${id}: ${err}`);
    })
  }

  const renameWorkspace = (id: string) => {

  }

  const setActiveWorkspace = (id: string) => {
    GetWorkspace(id).then(res => {
      ctx.dispatch({type: 'activeWorkspace', workspace: res});
    }).catch(err => console.log(`error on get workspace by id==${id}: `, err))
  }


  const selectWorkspace = (id: string) => {
    setActiveWorkspace(id);
    onClick();
  };

  return (
    <div className="h-auto" onMouseEnter={onMouseEnter} 
    onMouseLeave={onMouseLeave} >
      <li className="px-4 py-2 text-sm text-gray-800 hover:bg-gray-100 cursor-pointer"
             onClick={() => {selectWorkspace(id)}}>{name}</li>

    {isHovered && <WorkspaceMenu 
      items={[
        {text: "Remove", onClick: () => {removeWorkspace(id)}},
        {text: "Rename", onClick: () => {renameWorkspace(id)}},
        ]} />}
             </div>);
}

type WorkspaceMenuProps = {
    items: WorkspaceMenuItemProps[];
}

type WorkspaceMenuItemProps = {
  text: string;
  onClick: () => void;
};

const WorkspaceMenu: React.FC<WorkspaceMenuProps> = ({items}) => {
  return (<div className="top-full z-10 w-40 mt-2 bg-white border border-gray-300 rounded-md shadow-lg">
    <ul className="py-1">
      {items.map((it, i) => {
        return <WorkspaceMenuItem key={i} text={it.text} onClick={it.onClick} />;
      })}
    </ul>
  </div>);
}

const WorkspaceMenuItem: React.FC<WorkspaceMenuItemProps> = ({text, onClick}) => {
  return (<li className="px-4 py-2 text-sm text-gray-800 hover:bg-gray-100 cursor-pointer" onClick={onClick}>{text}</li>);
}
