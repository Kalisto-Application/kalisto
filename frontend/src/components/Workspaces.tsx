import React, { useState, useContext, useMemo, useEffect } from 'react';
import { Context } from '../state';

import {
  RenameWorkspace,
  WorkspaceList as GetWorkspaceList,
} from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';

import Dropdown, { DropdownItemProps } from './../ui/Dropdown';

import folderIcon from '../../assets/icons/folder.svg';
import dropdownIcon from '../../assets/icons/dropdown.svg';
import editIcon from '../../assets/icons/edit.svg';
import deleteIcon from '../../assets/icons/delete.svg';
import plusIcon from '../../assets/icons/plus.svg';
import CreateWorkspacePopup from './CreateWorkspacePopup';
import DeleteWorkspaceConfirmationPopup from './DeleteWorkspaceConfirmationPopup';

export const WorkspaceList: React.FC = () => {
  const ctx = useContext(Context);
  const [renameI, setRenameI] = useState(-1);
  const [isOpenCreateWorkspace, setIsOpenCreateWorkspace] = useState(false);
  const [isOpenDeletePopup, setIsOpenDeletePopup] = useState('');

  const items = ctx.state.workspaceList;
  const activeWorkspace = ctx.state.activeWorkspace;

  useEffect(() => {
    if (items?.length === 0) {
      setIsOpenCreateWorkspace(true);
    }
  }, [items]);

  const renameWorkspace = (id: string, name: string) => {
    RenameWorkspace(id, name)
      .then((_) => {
        setRenameI(-1);
        ctx.dispatch({ type: 'renameWorkspace', id: id, name: name });
      })
      .catch((err) => {
        console.log(
          `failed to rename workspace id=${id}, new name=${name}: ${err}`
        );
      });
  };

  const setActiveWorkspace = (id: string) => {
    GetWorkspaceList(id)
      .then((res) => {
        ctx.dispatch({ type: 'workspaceList', workspaceList: res });
      })
      .catch((err) =>
        console.log(`failed to get active workspace, id==${id}: ${err}`)
      );
  };

  let menuItems: DropdownItemProps[] = activeWorkspace
    ? [
        {
          text: activeWorkspace?.name || '',
          icon: folderIcon,
          tip: <img src={dropdownIcon} />,
        },
      ]
    : [
        {
          text: 'No workspace found',
          icon: folderIcon,
          tip: <img src={dropdownIcon} />,
        },
      ];
  menuItems = menuItems.concat([
    {
      text: 'Add new workspace',
      tip: <img src={plusIcon} />,
      onClick: () => setIsOpenCreateWorkspace(true),
      divide: true,
    },
    ...(items?.map((it, i) => {
      return {
        text: it.name,
        onClick: () => setActiveWorkspace(it.id),
        edit: {
          inEdit: i === renameI,
          onEditDone: (newName: string) => renameWorkspace(it.id, newName),
        },
        menu: [
          {
            icon: editIcon,
            text: 'Edit',
            onClick: (e: React.MouseEvent) => {
              e.preventDefault();
              setRenameI(i);
            },
          },

          {
            icon: deleteIcon,
            text: 'Delete',
            onClick: (e: React.SyntheticEvent) => {
              e.preventDefault();
              setIsOpenDeletePopup(it.id);
            },
          },
        ],
      };
    }) || []),
  ]);

  const main = menuItems.shift();
  return (
    <>
      <CreateWorkspacePopup
        open={isOpenCreateWorkspace}
        onClose={() => setIsOpenCreateWorkspace(false)}
      />
      <DeleteWorkspaceConfirmationPopup
        id={isOpenDeletePopup}
        isOpen={isOpenDeletePopup !== ''}
        onClose={() => setIsOpenDeletePopup('')}
      />
      <Dropdown main={main!} items={menuItems} />
    </>
  );
};
