import React, { useContext, useState } from 'react';
import copyIcon from '../../assets/icons/copy.svg';
import deleteIcon from '../../assets/icons/delete.svg';
import editIcon from '../../assets/icons/edit.svg';
import {
  CreateScriptFile,
  RemoveScriptFile,
  UpdateScriptFile,
} from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import FileList from '../ui/FileList';
import NewScript from './../ui/NewScript';
import DeletePopup from './DeletePopup';

const ScriptCollectionView: React.FC = () => {
  const ctx = useContext(Context);
  if (!ctx.state.activeWorkspace) {
    return <></>;
  }

  const [isOpenDeletePopup, setIsOpenDeletePopup] = useState('');
  const [isOpenEditInput, setIsOpenEditInput] = useState(false);
  const [isModeSubMenu, setIsModeSubMenu] = useState('');

  const workspace = ctx.state.activeWorkspace;

  const activeScript = workspace.scriptFiles.find(
    (it) => it.id === ctx.state.activeScriptFileId
  );

  const setActiveScript = (id: string) => {
    ctx.dispatch({ type: 'setActiveScriptId', id });
  };

  // Add
  const addFile = (value: string) => {
    CreateScriptFile(workspace.id, value, '', '').then((res) => {
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
    const renamed = new models.File({
      ...activeScript,
      name: name,
    });

    UpdateScriptFile(workspace.id, renamed).then((res) =>
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
      activeScript?.content || '',
      activeScript?.headers || ''
    ).then((res) => {
      ctx.dispatch({ type: 'addScriptFile', file: res });
    });
  };
  // sub menu items
  const items = ctx.state.activeWorkspace?.scriptFiles.map((it) => {
    return {
      file: it,
      menu: [
        {
          icon: editIcon,
          text: 'Edit',
          onClick: () => {
            setIsOpenEditInput(true);
            setIsModeSubMenu('');
          },
        },

        {
          icon: copyIcon,
          text: 'Copy',
          onClick: () => {
            copyFile();
            setIsModeSubMenu('');
          },
        },
        {
          icon: deleteIcon,
          text: 'Delete',
          onClick: () => {
            setIsOpenDeletePopup(activeScript?.id || '');
            setIsModeSubMenu('');
          },
        },
      ],
    };
  });

  return (
    <>
      <NewScript addFile={addFile} />
      <FileList
        activeWorkspace={ctx.state.activeWorkspace}
        setActiveScript={setActiveScript}
        items={items}
        activeScriptId={activeScript?.id || ''}
        isOpenEditInput={isOpenEditInput}
        onCloseInput={() => setIsOpenEditInput(false)}
        editFile={renameFile}
        isModeSubMenu={isModeSubMenu}
        closeSubMenu={() => setIsModeSubMenu('')}
        openSubMenu={() => setIsModeSubMenu(activeScript?.id || '')}
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
