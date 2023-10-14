import React, { useContext, useState } from 'react';
import { Context } from '../state';
import { models } from '../../wailsjs/go/models';
import DeletePopup from './DeletePopup';
import FileList from '../ui/FileList';
import { UpdateWorkspace } from '../../wailsjs/go/api/Api';
import editIcon from '../../assets/icons/edit.svg';
import deleteIcon from '../../assets/icons/delete.svg';
import copyIcon from '../../assets/icons/copy.svg';

const ScriptCollectionView: React.FC = () => {
  const ctx = useContext(Context);
  const [isOpenDeletePopup, setIsOpenDeletePopup] = useState('');
  const [isOpenEditInput, setIsOpenEditInput] = useState(false);

  const activeScript = ctx.state.scriptIdFile;

  const setActiveScript = (id: string) => {
    ctx.dispatch({ type: 'setActiveScriptId', id });
  };

  // Add
  const addFile = (value: string) => {
    let updatedWs = new models.Workspace({
      ...ctx.state.activeWorkspace,
      scriptFiles: [
        ...(ctx.state.activeWorkspace?.scriptFiles || []),
        {
          name: value,
          createdAt: new Date(),
          content: '',
          headers: '',
          id: value,
        },
      ],
    });

    UpdateWorkspace(updatedWs).then((res) => {
      ctx.dispatch({
        type: 'updateWorkspace',
        workspace: updatedWs,
      });
    });
  };

  // Delete
  const deleteFile = (id: string) => {
    let updatedWs = new models.Workspace({
      ...ctx.state.activeWorkspace,
      scriptFiles: ctx.state.activeWorkspace?.scriptFiles?.filter(
        (file) => file.id !== activeScript
      ),
    });
    UpdateWorkspace(updatedWs).then((res) => {
      ctx.dispatch({
        type: 'updateWorkspace',
        workspace: updatedWs,
      });
    });
  };

  // Edit
  const edeitFile = (rename: string) => {
    let updatedWs = new models.Workspace({
      ...ctx.state.activeWorkspace,
      scriptFiles: ctx.state.activeWorkspace?.scriptFiles?.map((file) => {
        if (file.id === activeScript) {
          file.name = rename;
        }
        return file;
      }),
    });

    UpdateWorkspace(updatedWs).then((res) => {
      ctx.dispatch({
        type: 'updateWorkspace',
        workspace: updatedWs,
      });
    });
  };

  // Copy

  const CopyFile = () => {};

  // sub menu items
  const items = [
    {
      icon: editIcon,
      text: 'Edit',
      onClick: () => {
        setIsOpenEditInput(true);
      },
    },

    {
      icon: copyIcon,
      text: 'Copy',
      onClick: () => {
        CopyFile();
      },
    },
    {
      icon: deleteIcon,
      text: 'Delete',
      onClick: () => {
        setIsOpenDeletePopup(activeScript);
      },
    },
  ];

  return (
    <>
      {ctx.state.activeWorkspace ? (
        <>
          <FileList
            addFile={(value: string) => addFile(value)}
            activeWorkspace={ctx.state.activeWorkspace}
            setActiveScript={setActiveScript}
            items={items}
            isOpenEditInput={isOpenEditInput}
            onCloseInput={() => setIsOpenEditInput(false)}
            edeitFile={(value) => edeitFile(value)}
          />
          <DeletePopup
            id={isOpenDeletePopup}
            isOpen={isOpenDeletePopup !== ''}
            onClose={() => setIsOpenDeletePopup('')}
            deleteScript={() => deleteFile(activeScript)}
            title="Delete script?"
          />
        </>
      ) : null}
    </>
  );
};

export default ScriptCollectionView;
