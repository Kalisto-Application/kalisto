import React, { useState, useContext, useMemo } from 'react';
import { Context } from "../state";

import { NewWorkspace, GetWorkspace, DeleteWorkspace, RenameWorkspace } from "../../wailsjs/go/api/Api";
import { models } from "../../wailsjs/go/models";

import Dropdown, {DropdownItemProps} from './../ui/Dropdown';

import folderIcon from '../../assets/icons/folder.svg'
import dropdownIcon from '../../assets/icons/dropdown.svg'
import editIcon from '../../assets/icons/edit.svg'
import deleteIcon from '../../assets/icons/delete.svg'
import plusIcon from '../../assets/icons/plus.svg'

let called = false;

interface WorkspaceListProps {
    items?: models.Workspace[];
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

  console.log('items:', items)
  if (items?.length === 0 && !called) {
    called = true;
    newWorkspace();
  }

  const renameWorkspace = (id: string, name: string) => {
    RenameWorkspace(id, name).then(_ => {
      setRenameI(-1);
      ctx.dispatch({type: 'renameWorkspace', id: id, name: name})
    }).catch(err => {
      console.log(`failed to rename workspace id=${id}, new name=${name}: ${err}`)
    })
  }

  const setActiveWorkspace = (id: string) => {
    GetWorkspace(id).then(res => {
      ctx.dispatch({type: 'activeWorkspace', workspace: res});
    }).catch(err => console.log(`error on get workspace by id==${id}: `, err))
  }

  const removeWorkspace = (id: string) => {
    DeleteWorkspace(id).then(_ => {
      ctx.dispatch({type: 'removeWorkspace', id: id})
    }).catch(err => {
      console.log(`failed to remove workspace id=${id}: ${err}`);
    })
  }

  let menuItems: DropdownItemProps[] = activeWorkspace? [{
    text: activeWorkspace?.name || "",
    icon: folderIcon,
    tip: <img src={dropdownIcon} />,
  }] : [{
    text: "No workspace found",
    icon: folderIcon,
    tip: <img src={dropdownIcon} />,
  }] 
  menuItems = menuItems.concat([
    {
      text: "Add new workspace",
      tip: <img src={plusIcon} />,
      onClick: () => newWorkspace(),
      divide: true,
    },
    ...items?.map((it, i) => {
      return {
      text: it.name,
      onClick: () => setActiveWorkspace(it.id),
      edit: {
        inEdit: i === renameI,
        onEditDone: (newName: string) => renameWorkspace(it.id, newName),
      },
      menu: [
        {icon: editIcon, text: "Edit", onClick: (e: React.MouseEvent) => {e.preventDefault(); setRenameI(i)}},
        {icon: deleteIcon, text: "Delete", onClick: (e: React.MouseEvent) => {e.preventDefault(); removeWorkspace(it.id)}},
      ],
      }
    }) || []
  ])

  const main = menuItems.shift();
  return (
    <Dropdown main={main!} items={menuItems} />
  );
};
