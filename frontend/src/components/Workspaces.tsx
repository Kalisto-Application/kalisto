import React, { useState, useContext } from 'react';
import { Context } from "../state";

import { NewWorkspace, GetWorkspace, DeleteWorkspace, RenameWorkspace } from "../../wailsjs/go/api/Api";
import { models } from "../../wailsjs/go/models";

import Dropdown from './../ui/Dropdown';
import Menu from '../ui/Menu';
import folderIcon from '../icons/folder.svg'
import dropdownIcon from '../icons/dropdown.svg'
import editIcon from '../icons/edit.svg'
import deleteIcon from '../icons/delete.svg'
import plusIcon from '../icons/plus.svg'

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

  const menuItems = [
    {
      text: "Add new workspace",
      icon: plusIcon,
      onClick: () => newWorkspace(),
    },
    ... items.map((it, i) => {
    return {
      text: it.name,
      onClick: () => setActiveWorkspace(it.id),
      edit: {
        inEdit: i === renameI,
        onEditDone: (newName: string) => renameWorkspace(it.id, newName),
      },
      icon: folderIcon,
      tip: <img src={dropdownIcon} />,
      menu: [
        {icon: editIcon, text: "Edit", onClick: (e: React.MouseEvent) => {e.preventDefault(); setRenameI(i)}},
        {icon: deleteIcon, text: "Delete", onClick: () => removeWorkspace(it.id)},
      ],
    }
  })];

  return (
    <Dropdown main={{text: activeWorkspace?.name || ""}} items={menuItems} />
  );
};
