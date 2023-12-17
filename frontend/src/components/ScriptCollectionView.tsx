import React, { useContext, useState } from 'react';
import copyIcon from '../../assets/icons/copy.svg';
import deleteIcon from '../../assets/icons/delete.svg';
import editIcon from '../../assets/icons/edit.svg';
import gIcon from '../../assets/icons/g.svg';
import {
  CreateScriptFile,
  RemoveScriptFile,
  UpdateScriptFile,
} from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import CreateItem from '../ui/CreateItem';
import FileList from '../ui/FileList';
import DeletePopup from './DeletePopup';

const ScriptCollectionView: React.FC = () => {
  const ctx = useContext(Context);
  const workspace = ctx.state.activeWorkspace;
  if (!workspace) {
    return <></>;
  }

  const [isOpenDeletePopup, setIsOpenDeletePopup] = useState<string | null>('');
  const [isOpenEditInput, setIsOpenEditInput] = useState('');

  const activeScript = workspace.scriptFiles.find(
    (it) => it.id === ctx.state.activeScriptFileId
  );

  const setActiveScript = (id: string): void => {
    ctx.dispatch({ type: 'setActiveScriptId', id });
  };

  // Add
  const addFile = (value: string): void => {
    CreateScriptFile(workspace.id, value, '', '')
      .then((res) => {
        ctx.dispatch({ type: 'addScriptFile', file: res });
      })
      .catch(() => {});
  };

  // Delete
  const deleteFile = (): void => {
    if (isOpenDeletePopup != null) return;

    RemoveScriptFile(workspace.id, isOpenDeletePopup != null)
      .then((res) => {
        let ws = new models.Workspace({
          ...ctx.state.activeWorkspace,
          scriptFiles: [...res],
        });

        ctx.dispatch({ type: 'updateWorkspace', workspace: ws });
      })
      .catch(() => {});
  };

  // Edit
  const renameFile = (name: string): void => {
    const renamed = new models.File({
      ...workspace.scriptFiles.find((it) => it.id === isOpenEditInput),
      name: name,
    });

    UpdateScriptFile(workspace.id, renamed)
      .then((res) =>
        ctx.dispatch({
          type: 'updateScriptFile',
          file: renamed,
        })
      )
      .catch(() => {});
  };

  // Copy
  const copyFile = (id: string): void => {
    const file = workspace.scriptFiles.find((it) => it.id === id);

    if (!file) return;

    CreateScriptFile(
      workspace.id,
      `${file.name} copy`,
      file.content,
      file.headers
    )
      .then((res) => {
        ctx.dispatch({ type: 'addScriptFile', file: res });
      })
      .catch(() => {});
  };
  // sub menu items
  const items = ctx.state.activeWorkspace?.scriptFiles.map((it) => {
    return {
      file: it,
      inEdit: it.id === isOpenEditInput,
      isActive: it.id === activeScript?.id,
      onClick: () => setActiveScript(it.id),
      menu: [
        {
          icon: editIcon,
          text: 'Edit',
          onClick: () => {
            setIsOpenEditInput(it.id);
          },
        },
        {
          icon: copyIcon,
          text: 'Copy',
          onClick: () => {
            copyFile(it.id);
          },
        },
        {
          icon: deleteIcon,
          text: 'Delete',
          onClick: () => {
            setIsOpenDeletePopup(it.id);
          },
        },
      ],
    };
  });

  return (
    <div>
      <CreateItem
        text="Add Script"
        addItem={(value) => addFile(value)}
        placeholder="Name script"
      />
      <FileList
        items={items}
        onCloseInput={() => setIsOpenEditInput('')}
        editFile={renameFile}
        itemIcon={gIcon}
      />
      <DeletePopup
        id={isOpenDeletePopup}
        isOpen={isOpenDeletePopup !== ''}
        onClose={() => setIsOpenDeletePopup('')}
        deleteScript={() => deleteFile()}
        title="Delete script?"
      />
    </div>
  );
};

export default ScriptCollectionView;
