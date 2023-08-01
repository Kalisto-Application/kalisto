import React, { useState, useContext } from 'react';
import { Context } from "../state";

import { NewWorkspace, GetWorkspace, DeleteWorkspace, RenameWorkspace } from "../../wailsjs/go/api/Api";
import { models } from "../../wailsjs/go/models";

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
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" xmlns="http://www.w3.org/2000/svg">
          <g clip-path="url(#clip0_400_226)">
          <path d="M12.2569 4.06432H6.74063L4.94081 2.26449C4.85706 2.18046 4.75752 2.11382 4.64791 2.06843C4.5383 2.02304 4.42079 1.99978 4.30216 2H0.903137C0.66361 2 0.433894 2.09515 0.264523 2.26452C0.0951517 2.4339 0 2.66361 0 2.90314V12.2325C0.00017086 12.4614 0.0911636 12.6809 0.252997 12.8427C0.41483 13.0045 0.634274 13.0955 0.863141 13.0957H12.3143C12.5385 13.0955 12.7535 13.0064 12.9121 12.8478C13.0707 12.6892 13.1598 12.4742 13.16 12.25V4.96745C13.16 4.72793 13.0648 4.49821 12.8955 4.32884C12.7261 4.15947 12.4964 4.06432 12.2569 4.06432ZM0.903137 2.77412H4.30216C4.33632 2.77426 4.36903 2.78795 4.39312 2.81218L5.64525 4.06432H0.774118V2.90314C0.774118 2.86892 0.787711 2.8361 0.811907 2.81191C0.836102 2.78771 0.868919 2.77412 0.903137 2.77412ZM12.3859 12.25C12.3859 12.269 12.3783 12.2872 12.3649 12.3006C12.3515 12.314 12.3333 12.3216 12.3143 12.3216H0.863141C0.839635 12.3212 0.817184 12.3118 0.800561 12.2951C0.783937 12.2785 0.774451 12.2561 0.774118 12.2325V4.83843H12.2569C12.2911 4.83843 12.3239 4.85203 12.3481 4.87622C12.3723 4.90042 12.3859 4.93323 12.3859 4.96745V12.25Z" fill="#BEBEC3" stroke="#BEBEC3" stroke-width="0.5"/>
          </g>
          <defs>
          <clipPath id="clip0_400_226">
          <rect width="14" height="14" fill="white"/>
          </clipPath>
          </defs>
        </svg>
        </div>
        {activeWorkspace && (<span className="ml-[12px]">{activeWorkspace.name}</span>)}
        </div>
        <div>
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M5 9L7 11L9 9" stroke="#BEBEC3" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M9 5L7 3L5 5" stroke="#BEBEC3" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
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