import React, { useContext, useState } from 'react';
import { Context } from '../state';
import { models } from '../../wailsjs/go/models';
import DeletePopup from './DeletePopup';
import FileList from '../ui/FileList';
import { RemoveScriptFile, RenameScriptFile } from '../../wailsjs/go/api/Api';
import editIcon from '../../assets/icons/edit.svg';
import deleteIcon from '../../assets/icons/delete.svg';
import copyIcon from '../../assets/icons/copy.svg';
import { CreateScriptFile } from './../../wailsjs/go/api/Api';

const ScriptCollectionView: React.FC = () => {
  const ctx = useContext(Context);
  const [isOpenDeletePopup, setIsOpenDeletePopup] = useState('');
  const [isOpenEditInput, setIsOpenEditInput] = useState(false);
  const [isModeSubMenu, setIsModeSubMenu] = useState(false);

  const activeScript = ctx.state.scriptIdFile;
  const workspaceId = ctx.state.activeWorkspace?.id;

  const setActiveScript = (id: string) => {
    ctx.dispatch({ type: 'setActiveScriptId', id });
  };

  // Add
  const addFile = (value: string) => {
    CreateScriptFile(workspaceId || '', value, '').then((res) => {
      ctx.dispatch({ type: 'addScriptFile', scriptFile: res });
    });
  };

  // Delete
  const deleteFile = () => {
    RemoveScriptFile(workspaceId || '', activeScript).then((res) =>
      ctx.dispatch({ type: 'deleteScriptFile', listFiles: res })
    );
  };

  // Edit
  const editFile = (rename: string) => {
    RenameScriptFile(workspaceId || '', activeScript, rename).then((res) =>
      ctx.dispatch({
        type: 'renameScriptFile',
        idFile: activeScript,
        value: rename,
      })
    );
  };

  // Copy

  const CopyFile = () => {
    let nameScript = '';
    let contentScript = '';
    ctx.state.activeWorkspace?.scriptFiles.forEach((file) => {
      if (file.id === activeScript) {
        if (
          file.name.slice(-1).match(/[0-9]/) &&
          file.name.includes(` copy `)
        ) {
          let copyNumber = Number(file.name.slice(-1));
          let strlength = file.name.length - 1;
          let str = file.name.slice(0, strlength);

          nameScript = str + ++copyNumber;
          contentScript = file.content;
          return;
        }
        if (file.name.includes(` copy`)) {
          nameScript = `${file.name} 2`;
          contentScript = file.content;
          return;
        }
        nameScript = `${file.name} copy`;
        contentScript = file.content;
        return;
      }
    });

    CreateScriptFile(workspaceId || '', nameScript, contentScript).then(
      (res) => {
        ctx.dispatch({ type: 'addScriptFile', scriptFile: res });
      }
    );
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
        CopyFile();
        setIsModeSubMenu(false);
      },
    },
    {
      icon: deleteIcon,
      text: 'Delete',
      onClick: () => {
        setIsOpenDeletePopup(activeScript);
        setIsModeSubMenu(false);
      },
    },
  ];

  return (
    <>
      {workspaceId ? (
        <>
          <FileList
            addFile={(value: string) => addFile(value)}
            activeWorkspace={ctx.state.activeWorkspace}
            setActiveScript={setActiveScript}
            items={items}
            activeScript={activeScript}
            isOpenEditInput={isOpenEditInput}
            onCloseInput={() => setIsOpenEditInput(false)}
            editFile={(value) => editFile(value)}
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
      ) : null}
    </>
  );
};

export default ScriptCollectionView;
