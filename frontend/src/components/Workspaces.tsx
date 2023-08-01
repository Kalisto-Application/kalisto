import React, { useState, useContext } from 'react';
import { Context } from "../state";

import { NewWorkspace, GetWorkspace, DeleteWorkspace, RenameWorkspace } from "../../wailsjs/go/api/Api";
import { models } from "../../wailsjs/go/models";
import iconFolderUrl from "../icons/folder.svg";
import iconDropdownUrl from "../icons/dropdown.svg";

interface WorkspaceListProps {
    items: models.Workspace[];
    activeWorkspace?: models.Workspace;
}

export const WorkspaceList: React.FC<WorkspaceListProps> = ({items, activeWorkspace}) => {
  const ctx = useContext(Context);
  const [renameI, setRenameI] = useState(-1);

  const newWorkspace = () => {
    NewWorkspace()
    .then(res => {
      ctx.dispatch({type: 'newWorkspace', workspace: res});
    })
    .catch(err => console.log('error on new workspace: ', err))
  }

  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  const toggleDropdown = () => {
    setIsDropdownOpen((prevIsDropdownOpen) => !prevIsDropdownOpen);
  };

  const onNewWorkspace = () => {
    if (isDropdownOpen) {
      toggleDropdown();
    }
    newWorkspace();
  }

  const renameWorkspace = (id: string, name: string) => {
    RenameWorkspace(id, name).then(_ => {
      setRenameI(-1);
      ctx.dispatch({type: 'renameWorkspace', id: id, name: name})
    }).catch(err => {
      console.log(`failed to rename workspace id=${id}, new name=${name}`)
    })
  }

  let containerClass = "w-[220px] bg-primaryFill relative border-2 border-borderFill"
  if (isDropdownOpen) {
    containerClass += " z-10"
  }

  return (
    <div className={containerClass}>
      <button
        type="button"
        className="flex items-center p-[12px] text-sm"
        onClick={toggleDropdown}
      >
        <div className='flex flex-1 m-[4px]'>
          <div className='mt-[4px]'>
        <img src={iconFolderUrl} />
        </div>
        {activeWorkspace && (<span className="ml-[12px]">{activeWorkspace.name}</span>)}
        </div>
        <div>
      <img src={iconDropdownUrl} />
        </div>
      </button>

      {isDropdownOpen && (
        <div className="absolute top-full">

          { items && (<ul className="py-1">
            {items.map((item, index) => {
              if (index === renameI) {
                return <WorkspaceRenameInput key={index} text={item.name} rename={(newName) => {renameWorkspace(item.id, newName)}}></WorkspaceRenameInput>;
              }
              return <WorkspaceItem key={index} 
                    id={item.id} 
                    onClick={()=> {setIsDropdownOpen(false)}} 
                    setRename={() => {setRenameI(index)}}
                    name={item.name} />
            })}
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
  setRename: () => void;
  name: string;
}

const WorkspaceItem: React.FC<WorkspaceItemProps> = ({ id, onClick, setRename, name}) => {
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
        {text: "Rename", onClick: setRename},
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

type WorkspaceRenameInputProps = {
  text: string;
  rename: (name: string) => void
}

const WorkspaceRenameInput: React.FC<WorkspaceRenameInputProps> = ({text, rename}) => {
  const [newText, setNewText] = useState(text)

  const onEnter = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      rename(newText);
    }
  }

  return (
    <div>
      <input type="text" value={newText} onChange={(e) => {setNewText(e.target.value)}} onKeyDown={onEnter} autoFocus></input>
    </div>
  );
}