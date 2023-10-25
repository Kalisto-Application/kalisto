import React, { useContext, useState } from 'react';
import copyIcon from '../../assets/icons/copy.svg';
import deleteIcon from '../../assets/icons/delete.svg';
import editIcon from '../../assets/icons/edit.svg';
import { RemoveScriptFile, RenameWorkspace } from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import FileList from '../ui/FileList';
import { CreateScriptFile } from '../../wailsjs/go/api/Api';
import DeletePopup from './DeletePopup';

const ScriptCollectionView: React.FC = () => {
  const ctx = useContext(Context);
  if (!ctx.state.activeWorkspace) {
    return <></>;
  }

  const [isOpenDeletePopup, setIsOpenDeletePopup] = useState('');
  const [isOpenEditInput, setIsOpenEditInput] = useState(false);
  const [isModeSubMenu, setIsModeSubMenu] = useState(false);

  const workspace = ctx.state.activeWorkspace;

  const activeScript = workspace.scriptFiles.find(
    (it) => it.id === ctx.state.activeScriptFileId
  );

  const getActiveScript = (id: string) => {

    ctx.dispatch({ type: 'setActiveScriptId', id });
  };

  // Add
  const addFile = (value: string) => {
    CreateScriptFile(workspace.id, value, '').then((res) => {
      ctx.dispatch({ type: 'addScriptFile', file: res });
    });
  };

  // Delete
  const deleteFile = () => {

    RemoveScriptFile(workspace.id, activeScript?.id || '').then((res) => {

      let ws = new models.Workspace({
        ...ctx.state.activeWorkspace,
        scriptFiles: [...res],
      });

      ctx.dispatch({ type: 'updateWorkspace', workspace: ws });
    });
  };

  // Edit
  const renameFile = (name: string) => {
    console.log('idd',workspace.id);
    debugger
    const renamed = new models.File({
      ...activeScript,
      name: name,
    });
    RenameWorkspace(workspace.id, name).then((res) =>
      ctx.dispatch({
        type: 'updateScriptFile',
        file: renamed,
      })
    );
  };

  // Copy

  const copyFile = () => {
    CreateScriptFile(
      workspace.id,
      `${activeScript?.name} copy`,
      activeScript?.content || ''
    ).then((res) => {
      ctx.dispatch({ type: 'addScriptFile', file: res });
    });
  };
  // sub menu items
  const items = [
    {
      icon: editIcon,
      text: 'Edit',
      onClick: () => {
        setIsOpenEditInput(true);
        setIsModeSubMenu(false);
      },
    },

    {
      icon: copyIcon,
      text: 'Copy',
      onClick: () => {
        copyFile();
        setIsModeSubMenu(false);
      },
    },
    {
      icon: deleteIcon,
      text: 'Delete',
      onClick: () => {
        setIsOpenDeletePopup(activeScript?.id || '');
        setIsModeSubMenu(false);
      },
    },
  ];

  return (
    <>
      <FileList
        addFile={addFile}
        activeWorkspace={ctx.state.activeWorkspace}
        getActiveScript={getActiveScript}
        items={items}
        activeScript={activeScript?.id || ''}
        isOpenEditInput={isOpenEditInput}
        onCloseInput={() => setIsOpenEditInput(false)}
        editFile={renameFile}
        isModeSubMenu={isModeSubMenu}
        closeSubMenu={() => setIsModeSubMenu(false)}
        openSubMenu={() => setIsModeSubMenu(true)}
      />
      <DeletePopup
        id={isOpenDeletePopup}
        isOpen={isOpenDeletePopup !== ''}
        onClose={() => setIsOpenDeletePopup('')}
        deleteScript={() => deleteFile()}
        title="Delete script?"
      />
    </>
  );
};

export default ScriptCollectionView;
